package service

import (
	"context"
	"douyin-app/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func Execute(c *gin.Context, userId int64, toUserId int64, actionType string) bool {
	switch actionType {
	case "1":
		return Follow(c, userId, toUserId)
	case "2":
		return UnFollow(c, userId, toUserId)
	}
	return false
}

func Follow(ctx context.Context, userId, toUserId int64) bool {
	if err := repository.GetRelationRepository().IsFollow(ctx, userId, toUserId); err == nil {
		return true
	}
	if err := repository.DB.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		// relation表新增relation
		if err := repository.GetRelationRepository().Create(ctx, userId, toUserId); err != nil {
			log.Println(err)
			return err
		}
		// to_user的User对象“粉丝总数” follower_count + 1
		if err := repository.GetRelationRepository().AddFollowerCount(ctx, toUserId); err != nil {
			return err
		}
		// user的User对象“关注总数” follow_count + 1
		if err := repository.GetRelationRepository().AddFollowCount(ctx, userId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}

func UnFollow(ctx context.Context, userId, toUserId int64) bool {
	if err := repository.GetRelationRepository().IsFollow(ctx, userId, toUserId); err != nil {
		return true
	}
	if err := repository.DB.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		// relation表新增relation
		if err := repository.GetRelationRepository().Delete(ctx, userId, toUserId); err != nil {
			log.Println(err)
			return err
		}
		// to_user的User对象“粉丝总数” follower_count + 1
		if err := repository.GetRelationRepository().ReduceFollowerCount(ctx, toUserId); err != nil {
			return err
		}
		// user的User对象“关注总数” follow_count + 1
		if err := repository.GetRelationRepository().ReduceFollowCount(ctx, userId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}
