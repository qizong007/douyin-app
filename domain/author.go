package domain

import "douyin-app/repository"

type Author struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func FillAuthor(user *repository.User) *Author {
	author := &Author{
		Id:            user.UserId,
		Name:          user.Username,
		FollowCount:   0,    // TODO
		FollowerCount: 0,    // TODO
		IsFollow:      true, // TODO
	}
	return author
}
