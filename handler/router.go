package handler

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	r.GET("ping", Ping)

	r.POST("/douyin/user/register/", RegisterHandler)

	r.POST("/douyin/user/login/", LoginHandler)
}
