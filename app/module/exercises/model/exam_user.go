package model

import (
	"gin-server/app/module/exercises/helper"
)

type ExamUser struct {
	BaseModel
	Nickname string `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar   string `gorm:"type:varchar(200)" json:"avatar"`
	Account  string `gorm:"type:varchar(50);not null" json:"account"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
	CreateAt string `gorm:"type:datetime" json:"create_at"`
}

type ExamUserResp struct {
	UserId   int    `gorm:"type:int(11);not null" json:"user_id"`
	Nickname string `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar   string `gorm:"type:varchar(200)" json:"avatar"`
	Token    string `json:"token"`
}

func FindUser(account string, password string) (user ExamUser, err error) {
	err = DB.Where("account = ? and password = ?", account, helper.Md5Encrypt(password)).Find(&user).Error
	return user, err
}

func (ExamUser) TableName() string {
	return "exam_user"
}
