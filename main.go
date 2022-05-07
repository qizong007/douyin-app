package main

import (
	"douyin-app/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	handler.Register(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
