package middleware

import (
	"douyin-app/util"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

func JWT(c *gin.Context) {

	token := c.Query("token")
	//解析token
	userId, err := util.ParseToken(token)
	if err != nil { //ParseToken只会返回两种错误
		if errors.Is(err, util.ErrNoAuth) {
			log.Println("JWTMiddleWare Token <Nil>")
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.NoAuth,
			})
			c.Abort()
			return
		}
		if errors.Is(err, util.ErrWrongAuth) {
			log.Println("JWTMiddleWare Token Wrong ,Err=", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
			c.Abort()
			return
		}
	}
	c.Set("userId", userId)
	c.Next()
}

func GetUserId(ctx *gin.Context) int64 {
	u, _ := ctx.Get("userId")
	return u.(int64)
}
