package models

import "time"

type STSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	RoleArn         string
	Bucket          string
	Region          string
}

type STSTokenResponse struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
	Bucket          string `json:"Bucket"`
	Region          string `json:"Region"`
	Path            string `json:"Path"`
}

type OSSConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Region          string
}

type CDNConfig struct {
	Domain string
	Key    string
}

// VideoMeta 包含从OSS返回的视频元数据
type VideoMeta struct {
	BasicInfo   BasicInfo   `json:"basic_info"`
	StorageInfo StorageInfo `json:"storage_info"`
	ContentInfo ContentInfo `json:"content_info"`
	RequestInfo RequestInfo `json:"request_info"`
}

type BasicInfo struct {
	FileName     string    `json:"file_name"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"` // 字节
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
}

type StorageInfo struct {
	ObjectType     string    `json:"object_type"`
	StorageClass   string    `json:"storage_class"`
	TransitionTime time.Time `json:"transition_time"` // 转储时间
}

type ContentInfo struct {
	MD5          string `json:"md5"`
	CRC64        string `json:"crc64"`
	AcceptRanges string `json:"accept_ranges"`
}

type RequestInfo struct {
	RequestID string `json:"request_id"`
	Date      string `json:"date"`
}
