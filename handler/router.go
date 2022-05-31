package handler

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	r.GET("ping", Ping)

	// user
	r.POST("/douyin/user/register/", RegisterHandler)
	r.POST("/douyin/user/login/", LoginHandler)
	r.GET("/douyin/user/", GetUserInfoHandler)

	// video
	r.POST("/douyin/publish/action/", VideoPublishHandler)
	r.GET("/douyin/publish/list/", VideoPublishedListHandler)
	r.GET("/douyin/feed/", VideoFeedHandler)

	//r.POST("")
	// favorite
	r.POST("/douyin/favorite/action/", VideoFavoriteHandler)
	r.GET("/douyin/favorite/list/", VedioFavoriteListHandler)

	// relation
	r.POST("douyin/relation/action/", FollowActionHandler)
	r.GET("douyin/relation/follow/list/", GetFollowListHandler)
	r.GET("douyin/relation/follower/list/", GetFollowerListHandler)
}
