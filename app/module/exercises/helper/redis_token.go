package helper

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"gin-server/app/module/exercises/db"
	"strings"
	"time"
)

// GenerateToken 时间戳+随机数
func GenerateToken(userId string) string {
	timestamp := time.Now().UnixNano()
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	data := fmt.Sprintf("%d%s", timestamp, base64.URLEncoding.EncodeToString(randBytes))
	token := base64.URLEncoding.EncodeToString([]byte(data))
	prefix, _ := AESGCMEncrypt(userId)

	db.RedisClient.Set(GetAuthorizeKey(userId), token, 48*time.Hour)
	return fmt.Sprintf("%s.%s", prefix, token)
}

func CleanToken(userId string) {
	db.RedisClient.Del(GetAuthorizeKey(userId))
}

func ValidToken(token string) (userId string, err error) {
	parts := strings.Split(token, ".")
	prefix := parts[0]
	userToken := parts[1]
	userId, err = AESGCMDecrypt(prefix)
	if err != nil {
		fmt.Println("Failed to decrypt token:", err)
		return "", err
	}

	if db.RedisClient.Get(GetAuthorizeKey(userId)).Val() != userToken {
		fmt.Println("redis token not equal")
		return "", errors.New("token is invalid")
	}
	return userId, nil
}
