package handler

import (
	"douyin-app/domain"
	"douyin-app/middleware"
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
	//获取从JWTMiddleware解析好的userId
	userId := middleware.GetUserId(c)

	title := c.PostForm("title")
	if title == "" {
		log.Println("VideoPublishHandler Title <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
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
	//获取从JWTMiddleware解析好的userId
	loginUserId := middleware.GetUserId(c)

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

func VideoFeedHandler(c *gin.Context) {
	var (
		userId    int64
		err       error
		videoList []*repository.Video
	)

	latestTimeStr := c.Query("latest_time")

	//获取从JWTMiddleware解析好的userId
	userId = middleware.GetUserId(c)

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
