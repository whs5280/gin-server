package utils

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"path/filepath"
	"strings"
)

type Address struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Network  string `json:"network"`
}

func IpToAddress(ip string) (address *Address, err error) {
	address = &Address{}
	absPath, err := filepath.Abs("../data/ip2region.xdb")
	if err != nil {
		return nil, err
	}
	fmt.Println("dbPath:", absPath)

	searcher, err := xdb.NewWithFileOnly(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create searcher: %v", err)
	}
	defer searcher.Close()

	region, err := searcher.SearchByStr(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to create searcher: %v", err)
	}

	parts := strings.Split(region, "|")
	address.Country = parts[0]
	address.Province = parts[2]
	address.City = parts[3]
	address.Network = parts[4]

	return address, nil
}
