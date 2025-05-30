package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encrypt(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
