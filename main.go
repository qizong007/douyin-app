package main

import (
	"douyin-app/conf"
	"douyin-app/handler"
	"douyin-app/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	repository.InitDB()
	conf.InitConf("./conf/conf.yaml")

	r := gin.Default()
	handler.Register(r)
	r.Run(conf.Config.Server.Port) // listen and serve on 0.0.0.0:8080
}
