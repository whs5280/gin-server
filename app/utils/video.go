package utils

import (
	"context"
	"fmt"
	"gin-server/app/models"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	config := &models.OSSConfig{
		AccessKeyId:     "",                             // 阿里云AccessKeyId
		AccessKeySecret: "",                             // 阿里云AccessKeySecret
		Bucket:          "",                             // 阿里云Bucket
		Region:          "oss-cn-hangzhou.aliyuncs.com", //  阿里云Region
	}

	// 创建OSS客户端
	client, err := oss.New(config.Region, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("OSS init failed: %v", err)
	}

	bucket, err := client.Bucket(config.Bucket)
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
