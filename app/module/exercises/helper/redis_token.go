package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gin-server/app/module/exercises/db"
	"time"
)

// GenerateToken 时间戳+随机数
func GenerateToken(userId string) string {
	timestamp := time.Now().UnixNano()
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	data := fmt.Sprintf("%d%s", timestamp, base64.URLEncoding.EncodeToString(randBytes))
	token := base64.URLEncoding.EncodeToString([]byte(data))

	db.RedisClient.SetNX(GetAuthorizeKey(userId), token, 6*time.Hour)
	return token
}

func CleanToken(userId string) {
	db.RedisClient.Del(GetAuthorizeKey(userId))
}
