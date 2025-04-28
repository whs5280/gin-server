package utils

import (
	"fmt"
	"testing"
)

func Test_GetSTSToken(t *testing.T) {
	token, err := GetSTSToken(1)
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", token)
}
