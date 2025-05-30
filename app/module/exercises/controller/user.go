package controller

import (
	"gin-server/app/module/exercises/helper"
	"gin-server/app/module/exercises/model"
	"gin-server/app/module/exercises/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Login(g *gin.Context) {
	userService := service.UserService{G: g}
	account := g.Query("account")
	password := g.Query("password")

	user, err := userService.Login(account, password)
	if err != nil {
		helper.ResponseJson(g, true, "登录失败", err)
		return
	}

	userResp := new(model.ExamUserResp)
	userResp.Token = helper.GenerateToken(strconv.Itoa(user.ID))
	userResp.UserId = user.ID
	userResp.Nickname = user.Nickname

	helper.ResponseJson(g, false, "登录成功", userResp)
}

func Logout(g *gin.Context) {
	token := g.GetHeader("token")
	if token == "" {
		helper.ResponseJson(g, true, "token不能为空", nil)
		return
	}

	helper.CleanToken(token)
	helper.ResponseJson(g, false, "退出成功", nil)
}
