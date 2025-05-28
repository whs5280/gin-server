package config

import (
	"github.com/spf13/viper"
)

func init() {
	if err := Init("../config.yaml"); err != nil {
		panic(err)
	}
}

func GetDBConf() map[string]interface{} {
	dbConf := viper.GetStringMap("database")
	return dbConf
}

func GetRedisConf() map[string]interface{} {
	redisConf := viper.GetStringMap("redis")
	return redisConf
}

func GetString(key string) string {
	return viper.GetString(key)
}
