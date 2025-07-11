package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gin-server/app/models"
	"gin-server/app/module/exercises/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	LocalTmpDir = "/tmp/video_processing/"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "r-bp1u7zhpk8qbntsxmspd.redis.rds.aliyuncs.com:6379",
		Password: "aFO8wFOl1m1ffZGi",
		DB:       1,
	})
	// jwtSecret := []byte("0a0cdfbf-1385-4019-92e5-696855b5af3a")
	jwtSecret []byte
)

func init() {
	jwtSecret = []byte(config.GetString("jwt.secret"))
}

// HandleOSSCallback 事件通知：当PutObject到指定目录（如uploads/）时触发HTTP回调
func HandleOSSCallback(ctx context.Context, objectKey string, typeKey string) {
	retry := 0
	maxRetry := 3
	for retry < maxRetry {
		if typeKey == "local" {
			path, err := processLocalVideo(ctx, objectKey)
			log.Println("path" + path)
			if err == nil {
				break
			}
			retry++
			log.Printf("Retry %d for %s: %v", retry, objectKey, err)
			time.Sleep(time.Second * time.Duration(retry*2))
		}

		if typeKey == "oss" {
			err := processOssVideo(ctx, objectKey)
			if err == nil {
				break
			}
			retry++
			log.Printf("Retry %d for %s: %v", retry, objectKey, err)
			time.Sleep(time.Second * time.Duration(retry*2))
		}
	}

	// 上报转码结果到监控系统
	/*metrics.Send("video_transcode", map[string]interface{}{
		"success": err == nil,
		"retries": retry,
		"key":     objectKey,
	})*/

}

// ProcessVideo 处理本地视频
func processLocalVideo(ctx context.Context, videoPath string) (string, error) {
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("input file not exist: %s", videoPath)
	}

	// 输出目录
	outputDir := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + "_hls"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output dir: %v", err)
	}

	// 设置输出路径
	m3u8Path := filepath.Join(outputDir, "playlist.m3u8")
	segmentPattern := filepath.Join(outputDir, "segment_%03d.ts")

	// 构建FFmpeg命令
	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-i", videoPath, // 输入文件
		"-c:v", "libx264", // H.264视频编码
		"-c:a", "aac", // AAC音频编码
		"-hls_time", "10", // 每个TS片段10秒
		"-hls_list_size", "0", // 无限播放列表长度（适合点播）
		"-hls_segment_filename", segmentPattern, // TS片段命名模式
		"-f", "hls", // 输出HLS格式
		m3u8Path, // M3U8文件路径
	)

	logFile, err := os.Create(filepath.Join(outputDir, "transcode.log"))
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %v", err)
	}
	defer logFile.Close()
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// 执行转码
	startTime := time.Now()
	if err := cmd.Run(); err != nil {
		// 清理失败时产生的部分文件
		_ = os.RemoveAll(outputDir)
		return "", fmt.Errorf("ffmpeg failed: %v (see log %s)", err, logFile.Name())
	}

	if _, err := os.Stat(m3u8Path); err != nil {
		return "", fmt.Errorf("m3u8 file not generated: %v", err)
	}

	duration := time.Since(startTime).Round(time.Second)
	fmt.Printf("Transcode success in %s: %s\n", duration, m3u8Path)
	return m3u8Path, nil
}

// ProcessVideo 处理OSS视频
func processOssVideo(ctx context.Context, objectKey string) error {
	oosConfig := &models.OSSConfig{
		AccessKeyId:     "",                             // 阿里云AccessKeyId
		AccessKeySecret: "",                             // 阿里云AccessKeySecret
		Bucket:          "",                             // 阿里云Bucket
		Region:          "oss-cn-hangzhou.aliyuncs.com", //  阿里云Region
	}

	// 创建OSS客户端
	client, err := oss.New(oosConfig.Region, oosConfig.AccessKeyId, oosConfig.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("OSS init failed: %v", err)
	}

	bucket, err := client.Bucket(oosConfig.Bucket)
	if err != nil {
		return fmt.Errorf("get bucket failed: %v", err)
	}

	// 下载原始视频到本地
	localMP4 := filepath.Join(LocalTmpDir, filepath.Base(objectKey))
	if err := os.MkdirAll(LocalTmpDir, 0755); err != nil {
		return err
	}
	if err := bucket.GetObjectToFile(objectKey, localMP4); err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer os.Remove(localMP4)

	// 转码为HLS
	outputPrefix := strings.TrimSuffix(objectKey, filepath.Ext(objectKey))
	outputDir := filepath.Join(LocalTmpDir, outputPrefix)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(outputDir)

	m3u8Path := filepath.Join(outputDir, "playlist.m3u8")
	cmd := exec.CommandContext(ctx, LocalTmpDir,
		"-i", localMP4,
		"-c:v", "libx264", // H.264编码
		"-hls_time", "10", // 每个TS片段10秒
		"-hls_list_size", "0", // 不限播放列表长度
		"-hls_segment_filename", filepath.Join(outputDir, "segment_%03d.ts"),
		"-f", "hls",
		m3u8Path,
	)
	cmd.Stderr = os.Stdout // 输出FFmpeg日志
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("transcode failed: %v", err)
	}

	// 上传HLS文件到OSS
	return filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		relPath, _ := filepath.Rel(LocalTmpDir, path)
		return bucket.PutObjectFromFile(
			filepath.Join("hls", relPath),
			path,
			oss.ContentType("application/vnd.apple.mpegurl"), // 正确设置MIME类型
		)
	})
}

// EncryptAndTranscode 每个视频独立密钥
func EncryptAndTranscode(inputPath, outputDir string) error {
	keyPath := filepath.Join(outputDir, "encryption.key")

	dir := filepath.Dir(keyPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return err
	}

	openssl := exec.Command("openssl", "rand", "-out", keyPath, "16")
	if output, err := openssl.CombinedOutput(); err != nil {
		fmt.Printf("Failed to generate key: %v\nOutput: %s\n", err, string(output))
		return err
	}

	// 创建keyInfo文件
	keyInfo := fmt.Sprintf(
		"https://your-server.com/keys/encryption.key\n%s\n",
		keyPath,
	)
	if err := os.WriteFile(filepath.Join(outputDir, "keyinfo.txt"), []byte(keyInfo), 0644); err != nil {
		fmt.Printf("Failed to create keyinfo.txt: %v\n", err)
		return err
	}

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-hls_time", "10",
		"-hls_key_info_file", filepath.Join(outputDir, "keyinfo.txt"),
		"-hls_playlist_type", "vod",
		filepath.Join(outputDir, "playlist.m3u8"),
	)
	return cmd.Run()
}

func IssueTempURL(videoID string) (key, token string, err error) {
	key, token, err = generateKeyAndToken(videoID)
	return key, token, err
}

// generateKeyAndToken 生成密钥和JWT令牌
func generateKeyAndToken(videoID string) (key string, token string, err error) {
	// 生成AES-128密钥（16字节）
	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", "", err
	}
	key = hex.EncodeToString(keyBytes)

	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(15 * time.Second).Unix(),
		"videoId":   videoID,
		"key":       key,
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenObj.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	err = redisClient.Set(fmt.Sprintf("key_%s", videoID), key, 20*time.Second).Err()
	return key, token, err
}

// ServeKey 验证并返回密钥
func ServeKey(tokenString string, viaRedis bool) (string, error) {
	// JWT验证
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println(ok)
		return "", fmt.Errorf("invalid token claims")
	}

	// Redis 验证
	if viaRedis {
		videoID := claims["videoId"].(string)
		storedKey, err := redisClient.Get(fmt.Sprintf("key_%s", videoID)).Result()
		if err != nil || storedKey != claims["key"] {
			log.Println(err)
			return "", fmt.Errorf("redis vaild fail")
		}
	}

	// 返回密钥（实际应加密传输）
	return claims["key"].(string), nil
}

// 流程

// 预备流程： `ffmpeg` 预生成 HLS + OSS/CDN 签名 【TS 分片缓存 30 天，m3u8 缓存 1 分钟】

// 1、客户端 -> 服务端（通过videoID 下发m3u8临时链接） 【30s 临时链接过期(jwt生成)】

// 2、客户端 -> 服务端（验证jwt, 下发视频加密的key、m3u8） 【cdn的URL鉴权、Referer防盗链】

// 3、客户端 （通过key 和 m3u8链接）播放视频

//    Client->>Server: 1. 请求播放 videoID
//    Server->>Client: 2. 返回临时 m3u8 URL（JWT 加密，60秒过期）
//    Client->>Server: 3. 用 JWT 请求 key 和 m3u8
//    Server->>Client: 4. 返回动态 m3u8（含签名 TS 链接）和 key（双重加密）
//    Client->>CDN:    5. 用 key 解密播放 TS 片段
//    CDN->>Client:    6. 返回加密视频流

//特性	      静态 m3u8	                        动态 m3u8
//生成方式	提前由 ffmpeg 生成，文件固定	    服务端按请求实时生成
//TS 链接	直接暴露原始路径（如 segment1.ts）	临时签名 URL（如 segment1.ts?token=xxx）
