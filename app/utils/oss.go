package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"gin-server/app/models"
	"net/url"
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
