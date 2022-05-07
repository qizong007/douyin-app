package handler

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	r.GET("/ping", Ping)
}
