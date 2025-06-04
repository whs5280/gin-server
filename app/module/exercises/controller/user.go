package controller

import (
	"gin-server/app/module/exercises/helper"
	"gin-server/app/module/exercises/model"
	"gin-server/app/module/exercises/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Register(g *gin.Context) {
	userService := service.UserService{G: g}

	var req model.ExamUserRegisterReq
	if err := g.ShouldBindQuery(&req); err != nil {
		helper.ResponseJson(g, true, "参数错误", err, http.StatusFailedDependency)
		return
	}

	user, err := userService.Register(req)
	if err != nil {
		g.Error(err)
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

	var req model.ExamUserLoginReq
	if err := g.ShouldBindQuery(&req); err != nil {
		helper.ResponseJson(g, true, "参数错误", err, http.StatusFailedDependency)
		return
	}

	user, err := userService.Login(req)
	if err != nil {
		g.Error(err)
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
		helper.ResponseJson(g, true, "token不能为空", nil, http.StatusFailedDependency)
		return
	}

	helper.CleanToken(token)
	helper.ResponseJson(g, false, "退出成功", nil)
}
