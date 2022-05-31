package handler

import (
	"douyin-app/domain"
	"douyin-app/repository"
	"douyin-app/service"
	"douyin-app/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const (
	FeedLimit = 30
)

func VideoPublishHandler(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		log.Println("VideoPublishHandler Token <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		log.Println("VideoPublishHandler Title <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	userId, err := util.ParseToken(token)
	if err != nil {
		log.Println("VideoPublishHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
		return
	}

	// 读取视频文件数据
	videoFileHeader, err := c.FormFile("data")
	if err != nil {
		log.Println("VideoPublishHandler FormFile Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	videoFile, err := videoFileHeader.Open()
	if err != nil {
		log.Println("VideoPublishHandler Open File Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	objectName := fmt.Sprintf("%d/%d_%s", userId, time.Now().Unix(), videoFileHeader.Filename)

	// 上传视频文件
	err = service.VideoPublish(objectName, videoFile)
	if err != nil {
		log.Println("VideoPublish Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	playUrl := service.VideoUploadUrlPrefix + objectName
	coverUrl := playUrl + service.VideoCoverSuffix

	if err = repository.GetVideoRepository().Create(c, &repository.Video{
		VideoId:  util.GenerateId(),
		UserId:   userId,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
		Title:    title,
	}); err != nil {
		log.Println("GetVideoRepository().Create Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
	})
}

func VideoPublishedListHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		log.Println("VideoPublishedListHandler Token <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	loginUserId, err := util.ParseToken(token)
	if err != nil {
		log.Println("VideoPublishedListHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
		return
	}

	userIdStr := c.Query("user_id")
	userId, err := util.Str2Int64(userIdStr)
	if err != nil {
		log.Println("Str2Int64 Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	videoList, err := repository.GetVideoRepository().FindByUserId(c, userId)
	if err != nil {
		log.Println("VideoPublishedListHandler FindByUserId Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	videoDOs, err := domain.FillVideoList(c, videoList, loginUserId, false)

	if err != nil {
		log.Println("FillVideoList Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"video_list": videoDOs,
		},
	})
}

//***********
func VideoFeedHandler(c *gin.Context) {
	var (
		userId    int64
		err       error
		videoList []*repository.Video
	)

	token := c.Query("token")
	latestTimeStr := c.Query("latest_time")

	if token != "" {
		userId, err = util.ParseToken(token)
		if err != nil {
			log.Println("VideoPublishHandler ParseToken Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
			return
		}
	}

	if latestTimeStr == "" { // 没有传入 latest_time
		videoList, err = repository.GetVideoRepository().FindWithLimit(c, FeedLimit)
		if err != nil {
			log.Println("GetVideoRepository().FindWithLimit Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	} else { // 传入了 latest_time
		latestTime, err := util.Str2Int64(latestTimeStr)
		if err != nil {
			log.Println("Str2Int64 Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}

		videoList, err = repository.GetVideoRepository().FindByCreateTimeWithLimit(c, latestTime, FeedLimit)
		if err != nil {
			log.Println("GetVideoRepository().FindByCreateTimeWithLimit Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	}

	videoDOs, err := domain.FillVideoList(c, videoList, userId, false)
	if err != nil {
		log.Println("FillVideoList Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"video_list": videoDOs,
			"next_time":  getMostEarlyTime(videoList),
		},
	})

}

// 获取本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
func getMostEarlyTime(videos []*repository.Video) int64 {
	if len(videos) == 0 {
		return time.Now().Unix()
	}
	return videos[len(videos)-1].CreateTime
}

/*
	点赞操作：
		1： 点赞，直接在favorite表创建 userid与favoriteId进行关联的数据，同时更新video表中 FavoriteCount 属性值
		2： 取消点赞，直接删除favorite表中这条记录，同时更新video表中 FavoriteCount 属性值

	actionTYpe == 1 点赞
		1) 先查询有无favorite record,无则创建，有则delete_time = 0
	actionType == 2 取消点赞
		1) delete_time = now()
*/
func VideoFavoriteHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		log.Println("VideoFavoriteHandler Token <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
	}

	userId, err := util.ParseToken(token)
	if err != nil {
		log.Println("GetUserInfoHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	videoId := c.Query("video_id")
	actionType := c.Query("action_type")

	vid, err := util.Str2Int64(videoId)
	if err != nil {
		log.Println("VideoFavoriteHandler Str2Int64 Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	video, err := repository.GetVideoRepository().FindByVideoId(c, vid)
	if err != nil {
		log.Println("VideoFavoriteHandler FindByVideoId Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	if actionType == "1" {
		video.FavoriteCount += 1

		favorite, err := repository.GetFavoriteRepository().FindFavoriteRecord(c, userId, vid)
		log.Println("[aaaa]", favorite)

		if err != nil {
			favoriteId := util.GenerateId()
			favorite := &repository.Favorite{
				FavoriteId: favoriteId,
				UserId:     userId,
				VideoId:    vid,
			}
			//创建favorite
			err = repository.GetFavoriteRepository().Create(c, favorite)
		} else {
			//点赞： 更新delete_time = 0
			favorite.DeleteTime = 0
			err = repository.GetFavoriteRepository().UpdateFavorite(c, userId, vid, favorite)
		}
	} else if actionType == "2" {
		video.FavoriteCount -= 1
		// 取消点赞： 更新delete_time = now()
		err = repository.GetFavoriteRepository().DeleteByUserIdAndVideoId(c, userId, vid)
	}

	err = repository.GetVideoRepository().VideoFavoriteAdd(c, video, vid)
	if err != nil {
		log.Println("VideoFavoriteHandler VideoFavoriteAdd Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
	})
	return
}

/*
	返回点赞list
 	1. 从favorite表中找出 userid= ?? 符合所有的favorite记录存储在favoriteList中
	2. 遍历favoriteList得到videoList
*/
func VedioFavoriteListHandler(c *gin.Context) {
	var (
		userId       int64
		err          error
		favoriteList []*repository.Favorite
		videoIdList  []int64
		videoList    []*repository.Video
	)

	token := c.Query("token")
	if token != "" {
		userId, err = util.ParseToken(token)
		if err != nil {
			log.Println("VideoPublishHandler ParseToken Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
		}
	}

	// 得到favoriteList
	favoriteList, err = repository.GetFavoriteRepository().FindVideoListByUseId(c, userId)
	if err != nil {
		log.Println("GetVideoRepository().FindWithLimit Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
	} else {
		for i := 0; i < len(favoriteList); i++ {
			videoIdList = append(videoIdList, favoriteList[i].VideoId)
		}
		for i := 0; i < len(videoIdList); i++ {
			tmp, err := repository.GetVideoRepository().FindByVideoId(c, videoIdList[i])
			videoList = append(videoList, tmp)
			if err != nil {
				videoList = append(videoList, tmp)
			}
		}
	}

	videoDOs, err := domain.FillVideoList(c, videoList, userId, true)
	if err != nil {
		log.Println("FillVideoList Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
	}
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"video_list": videoDOs,
		},
	})
}
