package utils

import (
	"fmt"
	"testing"
)

func Test_GetSignUrl(t *testing.T) {
	url := GetSignUrl("test/common/202504/23/8764da/8cebba701d83be1a.xls")
	fmt.Println(url)
}

func Test_GetSignUrlByCDN(t *testing.T) {
	url := GetSignUrlByCDN("test/common/202504/23/8764da/8cebba701d83be1a.xls")
	fmt.Println(url)
}

func Test_GetMediaInfo(t *testing.T) {
	meta, err := GetMediaInfo(
		"prod/content/202005/22/5203e6/51e7358b8464e7d21b9c130f4cd46f74.mp4",
		"video",
	)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("视频元数据:\n")
	fmt.Printf("文件名: %s\n", meta.BasicInfo.FileName)
	fmt.Printf("类型: %s\n", meta.BasicInfo.ContentType)
	fmt.Printf("大小: %.2f MB\n", float64(meta.BasicInfo.Size)/(1024*1024))
	fmt.Printf("修改时间: %s\n", meta.BasicInfo.LastModified.Format("2006-01-02 15:04:05"))
	fmt.Printf("存储类型: %s\n", meta.StorageInfo.StorageClass)
	fmt.Printf("MD5: %s\n", meta.ContentInfo.MD5)
	fmt.Printf("请求ID: %s\n", meta.RequestInfo.RequestID)
}
