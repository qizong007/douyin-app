package handler

import (
	"douyin-app/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("ping", Ping)

	// user
	userRouters := r.Group("/douyin/user")
	userRouters.POST("/register/", RegisterHandler)
	userRouters.POST("/login/", LoginHandler)

	//userInfo
	r.GET("/douyin/user/", GetUserInfoHandler).Use(middleware.JWT)

	// video
	publishRouters := r.Group("/douyin/publish").Use(middleware.JWT)
	publishRouters.POST("/action/", VideoPublishHandler)
	publishRouters.GET("/list/", VideoPublishedListHandler)

	// feed
	r.GET("/douyin/feed/", VideoFeedHandler).Use(middleware.JWT)

	// favorite
	favoriteRouters := r.Group("/douyin/favorite").Use(middleware.JWT)
	favoriteRouters.POST("/action/", VideoFavoriteHandler)
	favoriteRouters.GET("/list/", VedioFavoriteListHandler)

	// relation
	relationRouters := r.Group("/douyin/relation").Use(middleware.JWT)
	relationRouters.POST("/action/", FollowActionHandler)
	relationRouters.GET("/follow/list/", GetFollowListHandler)
	relationRouters.GET("/follower/list/", GetFollowerListHandler)

	// comment
	commentRouters := r.Group("/douyin/comment").Use(middleware.JWT)
	commentRouters.POST("/action/", CommentHandler)
	commentRouters.GET("/list/", CommentListHandler)
}
