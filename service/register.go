package service

import (
	"douyin-app/repository"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context, username string, password string) {
	//判断用户名是否在使用
	if exist := repository.ExistUserByName(username); exist {
		resp := util.HttpResponse{
			StatusCode: util.ParamError,
			StatusMsg:  "username is already in use",
			ReturnVal: map[string]interface{}{
				"user_id": 0,
				"token":   "",
			},
		}
		util.MakeResponse(c, &resp)
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//创建用户实例,存入注册信息
	ID, userId, err := repository.CreateUserInfo(username, string(hashPassword))
	if err != nil {
		resp := util.HttpResponse{
			StatusCode: util.InternalServerError,
			ReturnVal: map[string]interface{}{
				"user_id": 0,
				"token":   "",
			},
		}
		util.MakeResponse(c, &resp)
		return
	}
	//生成token
	token, _ := util.GenerateToken(ID, userId)
	resp := util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"user_id": userId,
			"token":   token,
		},
	}
	util.MakeResponse(c, &resp)
	return
}
