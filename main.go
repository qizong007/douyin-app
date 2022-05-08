package main

import (
	"douyin-app/conf"
	"douyin-app/handler"
	"douyin-app/repository"
	"github.com/gin-gonic/gin"
)

const (
	ConfPath = "./conf/conf.yaml"
)

func main() {
	conf.InitConf(ConfPath)
	repository.InitDB()

	r := gin.Default()
	handler.Register(r)
	r.Run(conf.Config.Server.Port) // listen and serve on 0.0.0.0:8080
}
