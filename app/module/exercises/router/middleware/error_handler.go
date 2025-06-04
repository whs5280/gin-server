package middleware

import (
	"errors"
	"gin-server/app/module/exercises/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var bizErr *helper.BusinessError
			if errors.As(err, &bizErr) {
				// 业务错误（返回用户友好提示）
				helper.ResponseJson(c, true, bizErr.Message, nil, bizErr.Code)
			} else {
				// 系统错误（日志记录详细错误，返回通用提示）
				log.Printf("Internal Error: %v", err)
				helper.ResponseJson(c, true, "系统繁忙，请稍后重试", nil, http.StatusInternalServerError)
			}
			return
		}
	}
}
