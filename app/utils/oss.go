package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"gin-server/app/models"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetSignUrl 阿里云OSS签名
func GetSignUrl(objectKey string) string {
	config := &models.OSSConfig{
		AccessKeyId:     "",                                                         // 阿里云AccessKeyId
		AccessKeySecret: "",                                                         // 阿里云AccessKeySecret
		Bucket:          "",                                                         // 阿里云Bucket
		Region:          fmt.Sprintf("https://%s.oss-cn-hangzhou.aliyuncs.com", ""), //  阿里云Region
	}

	expires := time.Now().UTC().Add(20 * time.Minute).Unix()

	// 1. 构造待签名的字符串
	canonicalizedResource := fmt.Sprintf("/%s/%s", config.Bucket, objectKey)
	stringToSign := fmt.Sprintf("GET\n\n\n%d\n%s", expires, canonicalizedResource)

	// 2. 计算签名
	h := hmac.New(sha1.New, []byte(config.AccessKeySecret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 3. 构造最终的URL
	signedURL := fmt.Sprintf("%s/%s?OSSAccessKeyId=%s&Expires=%d&Signature=%s",
		config.Region,
		url.QueryEscape(objectKey),
		config.AccessKeyId,
		expires,
		url.QueryEscape(signature),
	)

	return signedURL
}

// GetSignUrlByCDN CDN的加密链接功能
func GetSignUrlByCDN(objectKey string) string {
	config := &models.CDNConfig{
		Domain: "", // CDN的域名
		Key:    "", // CDN控制台配置的鉴权密钥
	}
	expires := time.Now().Add(20 * time.Minute).Unix()

	// CDN签名计算（算法可能为md5或sha1，需根据CDN配置调整）
	authStr := fmt.Sprintf("%s-%d-0-0-%s", objectKey, expires, config.Key)
	hash := md5.Sum([]byte(authStr))
	signature := hex.EncodeToString(hash[:])

	return fmt.Sprintf("https://%s/%s?auth_key=%d-0-0-%s",
		config.Domain,
		objectKey,
		expires,
		signature)
}

// GetMediaInfo 获取媒体信息
// https://help.aliyun.com/zh/oss/user-guide/video-information-extraction?spm=a2c4g.11186623.help-menu-31815.d_4_11_3_0_5.6f443097VPmOn2&scm=20140722.H_2362757._.OR_help-T_cn~zh-V_1 参考链接
func GetMediaInfo(objectName string, objectType string) (*models.VideoMeta, error) {
	config := &models.OSSConfig{
		AccessKeyId:     "",                             // 阿里云AccessKeyId
		AccessKeySecret: "",                             // 阿里云AccessKeySecret
		Bucket:          "",                             // 阿里云Bucket
		Region:          "oss-cn-hangzhou.aliyuncs.com", //  阿里云Region
	}

	// 创建OSS客户端
	client, err := oss.New(config.Region, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取Bucket
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取存储空间失败: %v", err)
	}
	isExist, err := bucket.IsObjectExist(objectName)
	if err != nil {
		return nil, fmt.Errorf("检查文件存在性失败: %v", err)
	}
	if !isExist {
		return nil, fmt.Errorf("文件%s不存在于Bucket%s中", objectName, config.Bucket)
	}

	// TODO 兼容audio音频
	props, err := bucket.GetObjectDetailedMeta(objectName)
	if err != nil {
		return nil, fmt.Errorf("获取视频元数据失败: %v", err)
	}

	return buildMediaInfo(objectName, props)
}

// buildMediaInfo 构建媒体信息
func buildMediaInfo(objectName string, props http.Header) (*models.VideoMeta, error) {
	// 解析LastModified时间
	lastModified, _ := time.Parse(time.RFC1123, props.Get("Last-Modified"))

	// 解析TransitionTime时间
	var transitionTime time.Time
	if tt := props.Get("X-Oss-Transition-Time"); tt != "" {
		transitionTime, _ = time.Parse(time.RFC1123, tt)
	}

	// 获取内容长度
	contentLength, err := strconv.ParseInt(props.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("解析Content-Length失败: %v", err)
	}

	// 组装元数据
	meta := &models.VideoMeta{
		BasicInfo: models.BasicInfo{
			FileName:     objectName,
			ContentType:  props.Get("Content-Type"),
			Size:         contentLength,
			LastModified: lastModified,
			ETag:         props.Get("Etag"),
		},
		StorageInfo: models.StorageInfo{
			ObjectType:     props.Get("X-Oss-Object-Type"),
			StorageClass:   props.Get("X-Oss-Storage-Class"),
			TransitionTime: transitionTime,
		},
		ContentInfo: models.ContentInfo{
			MD5:          props.Get("Content-Md5"),
			CRC64:        props.Get("X-Oss-Hash-Crc64ecma"),
			AcceptRanges: props.Get("Accept-Ranges"),
		},
		RequestInfo: models.RequestInfo{
			RequestID: props.Get("X-Oss-Request-Id"),
			Date:      props.Get("Date"),
		},
	}

	return meta, nil
}
