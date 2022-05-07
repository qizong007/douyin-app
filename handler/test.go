package handler

import (
	"douyin-app/util"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"msg": "ping",
		},
	})
}
