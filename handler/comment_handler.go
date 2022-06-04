package handler

import (
	"douyin-app/domain"
	"douyin-app/middleware"
	"douyin-app/repository"
	"douyin-app/service"
	"douyin-app/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strings"
)

const (
	PublishCommentType = "1"
	DeleteCommentType  = "2"
	MaxCommentLen      = 500
)

func CommentHandler(c *gin.Context) {
	//获取从JWTMiddleware解析好的userId
	userId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	actionType := c.Query("action_type")
	//根据actionType做功能的拆分
	switch actionType {
	case PublishCommentType:
		PublishCommentHandler(c, userId)
	case DeleteCommentType:
		DeleteCommentHandler(c)
	default:
		log.Println("CommentHandler Wrong ActionType")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
	}
}

func PublishCommentHandler(c *gin.Context, userId int64) {
	commentText := c.Query("comment_text")
	//去首尾空格并限制comment_text长度
	commentText = strings.TrimSpace(commentText)
	if len(commentText) > MaxCommentLen {
		log.Println("PublishCommentHandler CommentText Is Too Long ")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.CommentTooLong,
		})
		return
	}
	if len(commentText) == 0 {
		log.Println("PublishCommentHandler CommentText Is Null ")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.CommentIsEmpty,
		})
		return
	}

	//解析videoId
	reqVideoId := c.Query("video_id")
	videoId, err := util.Str2Int64(reqVideoId)
	if err != nil {
		log.Println("PublishCommentHandler ParseVideoId failed ")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}
	//判断videoId是否存在
	_, err = repository.GetVideoRepository().FindByVideoId(c, videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("PublishCommentHandler VideoId Not Exist")
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.VideoNotExist,
			})
			return
		}
		log.Println("PublishCommentHandler GetVideoRepository().FindByVideoId Failed")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	comment, user, err := service.PublishComment(c, userId, videoId, commentText)
	if err != nil {
		if errors.Is(err, util.ErrSensitiveComment) {
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.SensitiveComment,
			})
			return
		} else {
			log.Println("PublishCommentHandler Err=", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	}

	commentDO := domain.FillComment(comment, user)
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"comment": &commentDO,
		},
	})
}

func DeleteCommentHandler(c *gin.Context) {

	reqCommentId := c.Query("comment_id")
	commentId, err := util.Str2Int64(reqCommentId)
	if err != nil {
		log.Println("CommentHandler ParseCommentId Failed ")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	err = service.DeleteComment(c, commentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.CommentNotExist,
			})
			return
		} else {
			log.Println("service.DeleteComment Err=", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"comment": nil,
		},
	})
}

func CommentListHandler(c *gin.Context) {
	//获取从JWTMiddleware解析好的userId
	userId, err := middleware.GetUserId(c)
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}
	reqVideoId := c.Query("video_id")
	videoId, err := util.Str2Int64(reqVideoId)
	if err != nil {
		log.Println("CommentListHandler ParseVideoId failed ")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}
	//判断videoId是否存在
	_, err = repository.GetVideoRepository().FindByVideoId(c, videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("CommentListHandler VideoId Not Exist")
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.VideoNotExist,
			})
			return
		} else {
			log.Println("CommentListHandler GetVideoRepository().FindByVideoId Failed")
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.InternalServerError,
			})
			return
		}
	}

	comments, err := repository.GetCommentRepository().FindByVideoId(c, videoId)
	if err != nil {
		log.Println("CommentListHandler GetCommentRepository().FindByVideoId Failed")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	commentDOs, err := domain.FillCommentList(c, comments, userId)
	if err != nil {
		log.Println("CommentListHandler domain.FillCommentList Failed")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"comment_list": commentDOs,
		},
	})
}
