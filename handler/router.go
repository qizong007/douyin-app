package handler

import (
	"douyin-app/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("ping", Ping)

	// user
	r.POST("/douyin/user/register/", RegisterHandler)
	r.POST("/douyin/user/login/", LoginHandler)

	// feed
	r.GET("/douyin/feed/", VideoFeedHandler)

	// commentList
	r.GET("/douyin/comment/list/", CommentListHandler)

	r.Use(middleware.JWT)
	{
		//userInfo
		r.GET("/douyin/user/", GetUserInfoHandler)

		// video
		r.POST("/douyin/publish/action/", VideoPublishHandler)
		r.GET("/douyin/publish/list/", VideoPublishedListHandler)

		// favorite
		r.POST("/douyin/favorite/action/", VideoFavoriteHandler)
		r.GET("/douyin/favorite/list/", VideoFavoriteListHandler)

		// relation
		r.POST("douyin/relation/action/", FollowActionHandler)
		r.GET("douyin/relation/follow/list/", GetFollowListHandler)
		r.GET("douyin/relation/follower/list/", GetFollowerListHandler)

		// comment
		r.POST("/douyin/comment/action/", CommentHandler)
	}

}
