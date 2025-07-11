package utils

import (
	"context"
	"fmt"
	"testing"
)

func Test_HandleOSSCallback(t *testing.T) {
	path := "../resource/m3u8/9580b9b06c.mp4"
	HandleOSSCallback(context.Background(), path, "local")
	fmt.Println("done")
}

func Test_EncryptAndTranscode(t *testing.T) {
	err := EncryptAndTranscode("../resource/m3u8/9580b9b06c.mp4", "../resource/m3u8/9580b9b06c_enc")
	if err != nil {
		panic(err)
	}
	fmt.Println("done")
}

func Test_IssueTempURL(t *testing.T) {
	key, token, err := IssueTempURL("9580b9b06c")
	if err != nil {
		panic(err)
	}
	fmt.Println(key, token)
}

func Test_ServeKey(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE3NDk0NjE0ODAsImtleSI6IjQ5YzE3ZmU0ZGFkYmY0NjI5ZTNjYTQ0ZjQzZWM2ODA5IiwidmlkZW9JZCI6Ijk1ODBiOWIwNmMifQ.7cRcjw1UblxRWeS3JsV_ba-j_iByGceI7H6tQnga9pc"
	key, err := ServeKey(token, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(key)
}
