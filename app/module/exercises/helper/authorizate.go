package helper

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func CommonGetUserId(g *gin.Context) int {
	userIdValue, _ := g.Get("userId")
	userIdStr, _ := userIdValue.(string)
	userId, _ := strconv.Atoi(userIdStr)

	return userId
}
