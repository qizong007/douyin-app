package handler

import (
	"douyin-app/service"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Req struct {
	Username string `validate:"required,max=32"`
	Password string `validate:"required,min=6,max=32"`
}

func RegisterHandler(c *gin.Context) {
	//解析,处理参数
	var req Req
	var resp util.HttpResponse
	req.Username = c.Query("username")
	req.Password = c.Query("password")
	req.Username = strings.TrimSpace(req.Username) //去首尾空格

	err := util.Validate.Struct(req) // 执行验证
	if err != nil {
		resp.StatusCode = util.ParamError
		util.MakeResponse(c, &resp)
		return
	}
	service.Register(c, Req.Username, Req.Password)
}

func LoginHandler(c *gin.Context) {
	//解析,处理参数
	var Req Req
	var resp util.HttpResponse
	Req.Username = c.Query("username")
	Req.Password = c.Query("password")
	//去username首尾空格
	Req.Username = strings.TrimSpace(Req.Username)

	validate := validator.New() // 创建验证器
	err := validate.Struct(Req) // 执行验证
	if err != nil {
		resp.StatusCode = util.ParamError
		util.MakeResponse(c, &resp)
		return
	}

	userId, token, err := service.Login(c, Req.Username, Req.Password)
	resp.ReturnVal = map[string]interface{}{
		"user_id": userId,
		"token":   token,
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //用户未注册
			resp.StatusCode = util.UserNotExist
			util.MakeResponse(c, &resp)
			return
		}

		if errors.Is(err, util.ErrWrongPassword) { //用户密码错误
			resp.StatusCode = util.WrongPassword
			util.MakeResponse(c, &resp)
			return
		}

		resp.StatusCode = util.InternalServerError
		util.MakeResponse(c, &resp)
		return
	}
	resp.StatusCode = util.Success
	util.MakeResponse(c, &resp)
	return
}
