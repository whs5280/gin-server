package middleware

import (
	"gin-server/app/module/exercises/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}

func CheckToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		helper.ResponseJson(c, true, "token is empty", nil, 401)
		c.Abort()
		return
	}
	userId, err := helper.ValidToken(token)
	if err != nil {
		helper.ResponseJson(c, true, "token is invalid", nil, 401)
		c.Abort()
		return
	}
	if userId == "" {
		helper.ResponseJson(c, true, "token is invalid", nil, 401)
		c.Abort()
	}
	c.Set("userId", userId)
}
