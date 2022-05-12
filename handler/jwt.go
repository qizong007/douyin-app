package handler

import (
	"douyin-app/util"
	"github.com/gin-gonic/gin"
)

//进行token检查
//errCode 可能为Success 0, NoAuth 2 ,ErrorAuth 3
func JWTHandler(c *gin.Context) (Id int64, userId int64, errCode int) {

	//获取query中的token
	token := c.Request.URL.Query().Get("token")
	if token == "" {
		return 0, 0, util.NoAuth
	}

	//我们使用之前定义好的解析JWT的函数来解析它
	claim, err := util.ParseToken(token)
	if err != nil {
		return 0, 0, util.ErrorAuth
	}

	return claim.Id, claim.UserId, util.Success
}
