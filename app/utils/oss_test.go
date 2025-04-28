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
