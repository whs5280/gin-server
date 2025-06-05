package db

import (
	"fmt"
	"gin-server/app/module/exercises/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
)

var gDB *gorm.DB

func InitDB() (*gorm.DB, error) {
	if gDB != nil {
		return gDB, nil
	}

	dbConf := config.GetDBConf()
	dbUrl := fmt.Sprintf("%s:%v@tcp(%s:3306)/%s?charset=%s&parseTime=True&loc=Local", dbConf["user"], dbConf["password"], dbConf["host"], dbConf["dbname"], dbConf["charset"])
	db, err := gorm.Open("mysql", dbUrl)
	gDB = db

	gDB.DB().SetMaxIdleConns(20)
	gDB.DB().SetMaxOpenConns(25)

	if viper.GetBool("app.debug") {
		gDB = gDB.Debug()
		log.Println("GORM Debug 模式已启用")
	}

	return gDB, err
}

func GetDB() *gorm.DB {
	var err error
	if gDB != nil {
		return gDB
	}

	gDB, err = InitDB()
	if err != nil {
		return nil
	}
	return gDB
}

func init() {
	_, err := InitDB()
	if err != nil {
		panic(err)
	}
}
