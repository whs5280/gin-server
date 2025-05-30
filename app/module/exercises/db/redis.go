package db

import (
	"fmt"
	"gin-server/app/module/exercises/config"
	"github.com/go-redis/redis"
	"log"
)

var RedisClient *redis.Client

func init() {
	InitRedis()
}

func GetClient() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}
	return nil
}

func InitRedis() {
	redisConf := config.GetRedisConf()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf["host"], redisConf["port"]),
		Password: redisConf["password"].(string),
		DB:       redisConf["db"].(int),
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	RedisClient = client
}
