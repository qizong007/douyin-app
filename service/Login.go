package service

import (
	"douyin-app/repository"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, username string, password string) {
	//若用户名不存在
	if exist := repository.ExistUserByName(username); !exist {
		resp := util.HttpResponse{
			StatusCode: util.ParamError,
			StatusMsg:  "username is not exist",
			ReturnVal: map[string]interface{}{
				"user_id": 0,
				"token":   "",
			},
		}
		util.MakeResponse(c, &resp)
		return
	}

	ID, UserId, ok := repository.VerifyPassword(username, password)
	if !ok {
		resp := util.HttpResponse{
			StatusCode: util.ErrorAuth,
			StatusMsg:  "password is wrong",
			ReturnVal: map[string]interface{}{
				"user_id": 0,
				"token":   "",
			},
		}
		util.MakeResponse(c, &resp)
		return
	}
	//生成token
	token, _ := util.GenerateToken(ID, UserId)
	resp := util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"user_id": UserId,
			"token":   token,
		},
	}
	util.MakeResponse(c, &resp)
	return
}
