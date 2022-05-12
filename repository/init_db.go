package repository

import (
	"context"
	"douyin-app/conf"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitDB() {
	InitMySQL()
	InitRedis()
}

var (
	//Redis相关全局变量
	RedisCtx context.Context
	RedisDB  *redis.Client

	//gorm全局变量
	DB *gorm.DB
)

func InitRedis() {
	RedisCtx = context.Background()
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     conf.Config.Redis.Addr,
		Password: conf.Config.Redis.Password,
		DB:       conf.Config.Redis.DB,
	})
	_, err := RedisDB.Ping(RedisCtx).Result()
	if err != nil {
		panic(err)
	}
}

func InitMySQL() {
	dsn := conf.Config.MYSQL.Username + ":" +
		conf.Config.MYSQL.Password + "@tcp(" +
		conf.Config.MYSQL.Addr + ")/" +
		conf.Config.MYSQL.Database +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) //这里用短变量声明会有歧义
	if err != nil {
		panic(err)
	}

	//绑定结构体
	err = DB.AutoMigrate(
		User{},
	)
	if err != nil {
		log.Panicln(err)
	}
}
