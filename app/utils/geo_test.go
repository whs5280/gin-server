package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_IpToAddress(t *testing.T) {
	//address, err := IpToAddress("95.142.107.181")
	address, err := IpToAddress("119.131.138.96")
	if err != nil {
		t.Error(err)
		return
	}
	jsonData, _ := json.MarshalIndent(address, "", "  ")
	fmt.Println(string(jsonData))
}
