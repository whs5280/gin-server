package controller

import (
	"gin-server/app/module/exercises/helper"
	"github.com/gin-gonic/gin"
)

func Index(g *gin.Context) {
	helper.ResponseJson(g, false, "success", "hello world")
}
