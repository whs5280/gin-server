package service

import (
	"gin-server/app/module/exercises/model"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	G *gin.Context
}

func (usv UserService) Login(account string, password string) (u model.ExamUser, err error) {
	s, err := model.FindUser(account, password)
	if err != nil {
		return
	}
	return s, nil
}
