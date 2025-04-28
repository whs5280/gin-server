package utils

import (
	"fmt"
	"gin-server/app/models"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"math/rand"
	"time"
)

// GetSTSToken 获取STS token
func GetSTSToken(userId int32) (response *models.STSTokenResponse, err error) {
	config := &models.STSConfig{
		AccessKeyID:     "", // 填写您的AccessKey ID
		AccessKeySecret: "", // 填写您的AccessKey Secret
		RoleArn:         "", // 填写您的角色的ARN
		Bucket:          "", // 填写Bucket所在地域对应的Bucket名称
		Region:          "", // 填写CDN的域名
	}

	sdkConfig := sdk.NewConfig().WithTimeout(10 * time.Second)
	cred := credentials.NewAccessKeyCredential(config.AccessKeyID, config.AccessKeySecret)
	client, err := sts.NewClientWithOptions(config.Region, sdkConfig, cred)

	token := fmt.Sprintf("yidoutang%d%d", time.Now().Unix(), rand.Intn(10000))
	path := fmt.Sprintf("test/content/%s/%s/%d", time.Now().Format("200601"), time.Now().Format("02"), userId)
	policy := fmt.Sprintf(`
{
  "Statement": [
    {
      "Action": [
        "oss:Put*"
      ],
      "Effect": "Allow",
      "Resource": [
        "acs:oss:*:*:%s",
        "acs:oss:*:*:%s/%s/*"
      ]
    }
  ],
  "Version": "1"
}`, config.Bucket, config.Bucket, path)

	// 发起请求
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = config.RoleArn
	request.RoleSessionName = token
	request.Policy = policy
	request.DurationSeconds = "1800"

	res, err := client.AssumeRole(request)
	if err != nil {
		return nil, err
	}

	return &models.STSTokenResponse{
		AccessKeyId:     res.Credentials.AccessKeyId,
		AccessKeySecret: res.Credentials.AccessKeySecret,
		SecurityToken:   res.Credentials.SecurityToken,
		Expiration:      res.Credentials.Expiration,
		Bucket:          config.Bucket,
		Region:          config.Region,
		Path:            path,
	}, nil
}
