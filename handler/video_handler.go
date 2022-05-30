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
	}

	title := c.PostForm("title")
	if title == "" {
		log.Println("VideoPublishHandler Title <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
	}

	userId, err := util.ParseToken(token)
	if err != nil {
		log.Println("VideoPublishHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	// 读取视频文件数据
	videoFileHeader, err := c.FormFile("data")
	if err != nil {
		log.Println("VideoPublishHandler FormFile Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
	}

	videoFile, err := videoFileHeader.Open()
	if err != nil {
		log.Println("VideoPublishHandler Open File Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
	}

	objectName := fmt.Sprintf("%d/%d_%s", userId, time.Now().Unix(), videoFileHeader.Filename)

	// 上传视频文件
	err = service.VideoPublish(objectName, videoFile)
	if err != nil {
		log.Println("VideoPublish Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
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
	}

	userId, err := util.ParseToken(token)
	if err != nil {
		log.Println("VideoPublishedListHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
	}

	videoList, err := repository.GetVideoRepository().FindByUserId(c, userId)
	if err != nil {
		log.Println("VideoPublishedListHandler FindByUserId Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
	}

	videoDOs, err := domain.FillVideoList(c, videoList, userId, false)
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
		}
	}

	if latestTimeStr == "" { // 没有传入 latest_time
		videoList, err = repository.GetVideoRepository().FindWithLimit(c, FeedLimit)
		if err != nil {
			log.Println("GetVideoRepository().FindWithLimit Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
		}
	} else { // 传入了 latest_time
		latestTime, err := util.Str2Int64(latestTimeStr)
		if err != nil {
			log.Println("Str2Int64 Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
		}

		videoList, err = repository.GetVideoRepository().FindByCreateTimeWithLimit(c, latestTime, FeedLimit)
		if err != nil {
			log.Println("GetVideoRepository().FindByCreateTimeWithLimit Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
		}
	}

	videoDOs, err := domain.FillVideoList(c, videoList, userId, false)
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
	log.Println("[-1]:", videoId)

	actionType := c.Query("action_type")

	vid, err := util.Str2Int64(videoId)
	log.Println("[0]:", vid)

	video, err := repository.GetVideoRepository().FindByVideoId(c, vid)
	//video.FavoriteCount += 1
	log.Println("[1]:", video.VideoId)
	/*
		点赞操作：
			1： 点赞，直接在favorite表创建 userid与favoriteId进行关联的数据，同时更新video表中 FavoriteCount 属性值
			2： 取消点赞，直接删除favorite表中这条记录，同时更新video表中 FavoriteCount 属性值
	*/
	if actionType == "1" {
		video.FavoriteCount += 1
		favoriteId := util.GenerateId()
		favorite := &repository.Favorite{
			FavoriteId: favoriteId,
			UserId:     userId,
			VideoId:    vid,
		}
		// 创建 favorite信息
		err = repository.GetFavoriteRepository().Create(c, favorite)
	} else if actionType == "2" {
		video.FavoriteCount -= 1
		// 删除favorite信息
		err = repository.GetFavoriteRepository().DeleteByUserIdAndVideoId(c, userId, vid)
		if err != nil {
			log.Println("aaaaaaaaa", err)
		}
	}

	log.Printf("favorite count %d", video.FavoriteCount)

	err = repository.GetVideoRepository().VideoFavoriteAdd(c, video, vid)
	if err != nil {
		log.Println("aaaAFASDFASDFaaaaaa", err)
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
