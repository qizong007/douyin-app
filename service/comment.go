package service

import (
	"douyin-app/repository"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
)

func PublishComment(c *gin.Context, userId int64, videoId int64, content string) (comment *repository.Comment, user *repository.User, err error) {
	//TODO  检查content敏感词

	comment = &repository.Comment{
		CommentId: util.GenerateId(),
		UserId:    userId,
		VideoId:   videoId,
		Content:   content,
	}
	err = repository.GetCommentRepository().Create(c, comment)
	if err != nil {
		return nil, nil, err
	}

	err = repository.GetCommentRepository().AddCommentCount(c, comment)
	if err != nil {
		return nil, nil, err
	}
	//此处userId是肯定存在的
	user, err = repository.GetUserRepository().FindByUserId(c, userId)
	if err != nil {
		return nil, nil, err
	}

	return comment, user, nil
}

func DeleteComment(c *gin.Context, commentId int64) (err error) {

	rowsAffected, err := repository.GetCommentRepository().DeleteByCommentId(c, commentId)
	if rowsAffected == 0 {
		return util.ErrCommentNotExist
	}
	if err != nil {
		return err
	}

	err = repository.GetCommentRepository().AddCommentCount(c, &repository.Comment{CommentId: commentId})
	if err != nil {
		return err
	}
	return nil
}
