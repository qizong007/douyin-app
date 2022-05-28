package handler

import (
	"douyin-app/service"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
	"log"
)

type RelationReq struct {
	UserID     int64
	ToUserID   int64
	ActionType string
}

func RelationHandler(c *gin.Context) {
	var req RelationReq
	token := c.Query("token")
	req.ToUserID, _ = util.Str2Int64(c.Query("to_user_id"))
	req.ActionType = c.Query("action_type") // 1-关注 2-取消关注
	userId, err := util.ParseToken(token)
	req.UserID = userId // 客户端的请求里没有接口文档里说的userId，直接通过token解析

	if err != nil {
		log.Println("RelationHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	if ok := service.Execute(c, req.UserID, req.ToUserID, req.ActionType); ok {
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.Success,
		})
		return
	}

}

func FollowListHandler(c *gin.Context) {
	//token := c.Query("token")
	//userId, err := util.ParseToken(token)
	//if err != nil {
	//	log.Println("RelationHandler ParseToken Failed", err)
	//	util.MakeResponse(c, &util.HttpResponse{
	//		StatusCode: util.WrongAuth,
	//	})
	//}
	//
	//util.MakeResponse(c, &util.HttpResponse{
	//	StatusCode: util.Success,
	//	ReturnVal: map[string]interface{}{
	//		"user_list": service.FollowList(userId),
	//	},
	//})
}

func FollowerListHandler(c *gin.Context) {

}
