package domain

import (
	"context"
	"douyin-app/repository"
	"log"
	"time"
)

type Comment struct {
	Id         int64   `json:"id"`
	User       *Author `json:"user"`
	Content    string  `json:"content"`
	CreateDate string  `json:"create_date"` //评论发布日期，格式 mm-dd
}

func FillComment(comment *repository.Comment, user *repository.User) *Comment {
	timeStr := time.Unix(comment.CreateTime, 0).Format("2006-01-02 15:04:05")

	return &Comment{
		Id: comment.CommentId,
		User: &Author{
			Id:            user.UserId,
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true, //这里作者是自己,所以是true
		},
		Content:    comment.Content,
		CreateDate: timeStr[5:10], //mm-dd
	}
}

func FillCommentList(ctx context.Context, comments []*repository.Comment, userId int64) ([]*Comment, error) {
	userIds := GetUserIdsFromCommentList(comments)
	authors, err := getAuthorsFromIds(ctx, userIds, userId)
	log.Println("ids", userIds)
	log.Println("authors", authors)
	if err != nil {
		return nil, err
	}
	commentDOs := make([]*Comment, len(comments))
	for i := range comments {
		commentDOs[i] = &Comment{
			Id:         comments[i].CommentId,
			User:       authors[i],
			Content:    comments[i].Content,
			CreateDate: time.Unix(comments[i].CreateTime, 0).Format("2006-01-02 15:04:05")[5:10], //mm-dd
		}
	}
	return commentDOs, nil
}

func GetUserIdsFromCommentList(comments []*repository.Comment) []int64 {
	ids := make([]int64, len(comments))
	for i := range comments {
		ids[i] = comments[i].UserId
	}
	return ids
}
