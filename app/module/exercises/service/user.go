package service

import (
	"gin-server/app/module/exercises/model"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	G *gin.Context
}

func (usv UserService) Register(req model.ExamUserRegisterReq) (u model.ExamUser, err error) {
	return model.Register(req)
}

func (usv UserService) Login(req model.ExamUserLoginReq) (u model.ExamUser, err error) {
	return model.FindUser(req)
}
