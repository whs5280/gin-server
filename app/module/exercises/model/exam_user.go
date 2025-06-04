package model

import (
	"errors"
	"fmt"
	"gin-server/app/module/exercises/helper"
	"github.com/jinzhu/gorm"
)

type ExamUser struct {
	BaseModel
	Nickname  string `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar    string `gorm:"type:varchar(200)" json:"avatar"`
	Account   string `gorm:"type:varchar(50);not null" json:"account"`
	Password  string `gorm:"type:varchar(100);not null" json:"password"`
	CreatedAt string `gorm:"type:datetime" json:"created_at"`
}

type ExamUserRegisterReq struct {
	Nickname string `form:"nickname" json:"nickname"`
	Account  string `form:"account" json:"account"`
	Password string `form:"password" json:"password"`
}

type ExamUserResp struct {
	UserId   int    `gorm:"type:int(11);not null" json:"user_id"`
	Nickname string `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar   string `gorm:"type:varchar(200)" json:"avatar"`
	Token    string `json:"token"`
}

func FindUser(account string, password string) (user ExamUser, err error) {
	err = DB.Where("account = ? and password = ?", account, helper.Md5Encrypt(password)).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, fmt.Errorf("用户不存在或密码错误")
	}
	return user, err
}

func IsRegister(account string) (bool, error) {
	var count int
	err := DB.Model(&ExamUser{}).Where("account = ?", account).Count(&count).Error
	return count > 0, err
}

func Register(req ExamUserRegisterReq) (user ExamUser, err error) {
	if isRegister, _ := IsRegister(req.Account); isRegister {
		return user, fmt.Errorf("账号已存在")
	}
	user.Nickname = req.Nickname
	user.Account = req.Account
	user.Password = helper.Md5Encrypt(req.Password)
	user.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	err = DB.Create(&user).Error
	return user, err
}

func (ExamUser) TableName() string {
	return "exam_user"
}
