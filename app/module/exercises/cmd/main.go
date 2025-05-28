package main

import (
	"fmt"
	"gin-server/app/module/exercises/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	gin.SetMode(gin.TestMode)

	r := router.InitRouter()

	port := "8080"
	fmt.Printf("Start to listening the incoming requests on http address: :%s\n", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
