package models

import (
	"log"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB = Init()

var RDB = InitRedisDB()

func Init() *gorm.DB {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm init error:", err)
	}
	return db
}

func InitRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
