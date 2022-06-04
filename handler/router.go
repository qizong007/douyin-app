package handler

import (
	"douyin-app/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("ping", Ping)

	r.GET("/douyin/user/", GetUserInfoHandler)

	r.GET("/douyin/feed/", VideoFeedHandler)

	r.Use(middleware.JWT)
	{
		// user
		r.POST("/douyin/user/register/", RegisterHandler)
		r.POST("/douyin/user/login/", LoginHandler)

		// video
		r.POST("/douyin/publish/action/", VideoPublishHandler)
		r.GET("/douyin/publish/list/", VideoPublishedListHandler)

		// favorite
		r.POST("/douyin/favorite/action/", VideoFavoriteHandler)
		r.GET("/douyin/favorite/list/", VedioFavoriteListHandler)

		// relation
		r.POST("douyin/relation/action/", FollowActionHandler)
		r.GET("douyin/relation/follow/list/", GetFollowListHandler)
		r.GET("douyin/relation/follower/list/", GetFollowerListHandler)

		// comment
		r.POST("/douyin/comment/action/", CommentHandler)
		r.GET("/douyin/comment/list/", CommentListHandler)
	}

}
