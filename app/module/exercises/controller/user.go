package controller

import (
	"gin-server/app/module/exercises/helper"
	"gin-server/app/module/exercises/model"
	"gin-server/app/module/exercises/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Register(g *gin.Context) {
	userService := service.UserService{G: g}

	var req model.ExamUserRegisterReq
	if err := g.ShouldBindQuery(&req); err != nil {
		helper.ResponseJson(g, true, "参数错误", err, 422)
		return
	}

	user, err := userService.Register(req)
	if err != nil {
		helper.ResponseJson(g, true, "注册失败", err)
		return
	}

	userResp := new(model.ExamUserResp)
	userResp.UserId = user.ID
	userResp.Nickname = user.Nickname
	userResp.Avatar = user.Avatar

	helper.ResponseJson(g, false, "注册成功", userResp)
}

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
