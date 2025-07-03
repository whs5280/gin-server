package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type WeatherInfo struct {
	City        string `json:"city"`
	Temperature string `json:"temperature"`
	Weather     string `json:"weather"`
	Wind        string `json:"wind"`
	Humidity    string `json:"humidity"`
}

func GeoPoint(address string) (res string, err error) {
	apiKey := os.Getenv("REST_API_KEY")
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s", apiKey, address)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	res = string(body)
	return
}
