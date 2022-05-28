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

func FollowActionHandler(c *gin.Context) {
	var req RelationReq
	token := c.Query("token")
	req.ToUserID, _ = util.Str2Int64(c.Query("to_user_id"))
	req.ActionType = c.Query("action_type") // 1-关注 2-取消关注
	userId, err := util.ParseToken(token)
	if err != nil {
		log.Println("RelationHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}
	req.UserID = userId // 客户端的请求里没有接口文档里说的userId，直接通过token解析
	err = service.RelationAction(c, req.UserID, req.ToUserID, req.ActionType)
	if err != nil {
		log.Println("RelationAction() Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
	})

}

func GetFollowListHandler(c *gin.Context) {
	// TODO
}

func GetFollowerListHandler(c *gin.Context) {
	// TODO
}
