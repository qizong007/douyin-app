package service

import (
	"context"
	"douyin-app/domain"
	"douyin-app/repository"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	FollowAction   = "1" // 1-关注
	UnFollowAction = "2" // 2-取消关注
)

func RelationAction(c *gin.Context, userId int64, toUserId int64, actionType string) error {
	switch actionType {
	case FollowAction:
		return Follow(c, userId, toUserId)
	case UnFollowAction:
		return UnFollow(c, userId, toUserId)
	}
	return nil
}

// Follow 关注
func Follow(ctx context.Context, userId, toUserId int64) error {
	// error不为空说明没有该关注记录，继续下面的关注操作
	if err := repository.GetFollowRepository().FindByUserId(ctx, userId, toUserId); err == nil {
		return err
	}
	// relation表新增relation
	if err := repository.GetFollowRepository().Create(ctx, userId, toUserId); err != nil {
		return err
	}
	// to_user的User对象“粉丝总数” follower_count + 1
	if err := repository.GetFollowRepository().AddFollowerCount(ctx, toUserId); err != nil {
		return err
	}
	// user的User对象“关注总数” follow_count + 1
	if err := repository.GetFollowRepository().AddFollowCount(ctx, userId); err != nil {
		return err
	}
	return nil
}

// UnFollow 取消关注
func UnFollow(ctx context.Context, userId, toUserId int64) error {
	// error为空说明已经有该关注记录，继续下面的取消关注操作
	if err := repository.GetFollowRepository().FindByUserId(ctx, userId, toUserId); err != nil {
		return err
	}
	// relation表新增relation
	if err := repository.GetFollowRepository().Delete(ctx, userId, toUserId); err != nil {
		return err
	}
	// to_user的User对象“粉丝总数” follower_count + 1
	if err := repository.GetFollowRepository().ReduceFollowerCount(ctx, toUserId); err != nil {
		return err
	}
	// user的User对象“关注总数” follow_count + 1
	if err := repository.GetFollowRepository().ReduceFollowCount(ctx, userId); err != nil {
		return err
	}
	return nil
}

// GetFollowList 获得关注列表
func GetFollowList(ctx context.Context, userId int64) []*domain.UserInfo {
	followList, err := repository.GetFollowRepository().FindByFromUserId(ctx, userId)
	if err != nil {
		log.Println("FindByFromUserId Failed", err)
		return nil
	}
	userIds := domain.GetToUserIdsFromFollowList(followList)
	userList, err := domain.GetUserInfosFromIds(ctx, userIds)
	if err != nil {
		log.Println("GetUserInfosFromIds Failed", err)
		return nil
	}
	return userList
}

// GetFollowerList 获得粉丝列表
func GetFollowerList(ctx context.Context, userId int64) []*domain.UserInfo {
	followList, err := repository.GetFollowRepository().FindByToUserId(ctx, userId)
	if err != nil {
		log.Println("FindByToUserId Failed", err)
		return nil
	}
	userIds := domain.GetFromUserIdsFromFollowList(followList)
	userList, err := domain.GetUserInfosFromIds(ctx, userIds)
	if err != nil {
		log.Println("GetUserInfosFromIds Failed", err)
		return nil
	}
	return userList
}
