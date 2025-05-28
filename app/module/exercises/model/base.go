package model

import (
	"gin-server/app/module/exercises/db"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type BaseModel struct {
	ID int `gorm:"AUTO_INCREMENT;NOT NULL;primary_key" json:"id"`
}

func init() {
	DB = db.GetDB()
}
