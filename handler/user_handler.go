package handler

import (
	"douyin-app/service"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

type req struct {
	Username string `validate:"required,max=32"`
	Password string `validate:"required,min=6,max=32"`
}

func RegisterHandler(c *gin.Context) {
	var Req req
	Req.Username = c.Query("username")
	Req.Password = c.Query("password")
	//去首尾空格
	Req.Username = strings.TrimSpace(Req.Username)

	validate := validator.New() // 创建验证器
	err := validate.Struct(Req) // 执行验证
	if err != nil {
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
			ReturnVal: map[string]interface{}{
				"user_id": 0,
				"token":   "",
			},
		})
		return
	}
	service.Register(c, Req.Username, Req.Password)
}

func LoginHandler(c *gin.Context) {
	var Req req
	Req.Username = c.Query("username")
	Req.Password = c.Query("password")
	//去username首尾空格
	Req.Username = strings.TrimSpace(Req.Username)

	validate := validator.New() // 创建验证器
	err := validate.Struct(Req) // 执行验证
	if err != nil {
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
			ReturnVal: map[string]interface{}{
				"user_id": 0,
				"token":   "",
			},
		})
		return
	}
	service.Login(c, Req.Username, Req.Password)
}
