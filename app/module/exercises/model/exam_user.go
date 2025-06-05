package model

import (
	"errors"
	"gin-server/app/module/exercises/helper"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ExamUser struct {
	BaseModel
	Nickname  string `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar    string `gorm:"type:varchar(200)" json:"avatar"`
	Account   string `gorm:"type:varchar(50);not null" json:"account"`
	Password  string `gorm:"type:varchar(100);not null" json:"password"`
	CreatedAt string `gorm:"type:datetime" json:"created_at"`
}

type ExamUserLoginReq struct {
	Account  string `form:"account" json:"account" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type ExamUserRegisterReq struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
	Account  string `form:"account" json:"account" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type ExamUserResp struct {
	UserId   int    `gorm:"type:int(11);not null" json:"user_id"`
	Nickname string `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar   string `gorm:"type:varchar(200)" json:"avatar"`
	Token    string `json:"token"`
}

func FindUser(req ExamUserLoginReq) (user ExamUser, err error) {
	err = DB.Where("account = ? and password = ?", req.Account, helper.Md5Encrypt(req.Password)).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, &helper.BusinessError{
			Code:    http.StatusUnauthorized,
			Message: "用户不存在或密码错误",
		}
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
		return user, &helper.BusinessError{
			Code:    http.StatusConflict,
			Message: "账号已存在",
		}
	}
	user.Nickname = req.Nickname
	user.Account = req.Account
	user.Password = helper.Md5Encrypt(req.Password)
	user.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	user.CreatedAt = helper.GetNowTime()
	err = DB.Create(&user).Error
	return user, err
}

func (ExamUser) TableName() string {
	return "exam_user"
}
