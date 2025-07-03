package main

import (
	"fmt"
	"gin-server/app/module/weather/service"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("please input city name")
		return
	}

	city := args[1]
	cityArr := strings.Split(city, ",")

	for i, singleCity := range cityArr {
		fmt.Printf("input %d: ------%s------\n", i, singleCity)
		info, err := service.GeoPoint(singleCity)
		if err != nil {
			return
		}
		fmt.Println(info)
	}
}
