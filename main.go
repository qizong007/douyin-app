package main

import (
	"douyin-app/conf"
	"douyin-app/handler"
	"douyin-app/repository"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
)

const (
	ConfPath = "./conf/conf.yaml"
)

func main() {
	conf.InitConf(ConfPath)
	repository.InitDB()
	util.InitIdGenerator()
	util.InitJWTVal()
	util.InitValidate()
	util.InitOSSClient(conf.Config)

	r := gin.Default()
	handler.Register(r)
	r.Run(conf.Config.Server.Port) // listen and serve on 0.0.0.0:8080
}
