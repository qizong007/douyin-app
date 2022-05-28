package service

import (
	"context"
	"douyin-app/repository"
	"github.com/gin-gonic/gin"
	"log"
)

func RelationAction(c *gin.Context, userId int64, toUserId int64, actionType string) error {
	switch actionType {
	case "1":
		return Follow(c, userId, toUserId)
	case "2":
		return UnFollow(c, userId, toUserId)
	}
	return nil
}

func Follow(ctx context.Context, userId, toUserId int64) error {
	if err := repository.GetFollowRepository().IsFollow(ctx, userId, toUserId); err == nil {
		return err
	}
	// relation表新增relation
	if err := repository.GetFollowRepository().Create(ctx, userId, toUserId); err != nil {
		log.Println(err)
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

func UnFollow(ctx context.Context, userId, toUserId int64) error {
	if err := repository.GetFollowRepository().IsFollow(ctx, userId, toUserId); err != nil {
		return err
	}
	// relation表新增relation
	if err := repository.GetFollowRepository().Delete(ctx, userId, toUserId); err != nil {
		//log.Println(err)
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
