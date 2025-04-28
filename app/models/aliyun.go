package models

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
