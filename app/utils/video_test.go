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
