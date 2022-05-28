package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Relation struct {
	ID         int64 `json:"id" gorm:"primaryKey"`
	UserID     int64 `json:"user_id" `
	FollowerID int64 `json:"follower_id"`
	CreateTime int64 `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime int64 `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime int64 `json:"delete_time"`
}

type IRelationRepository interface {
	IsFollow(context.Context, int64, int64) error
	Create(context.Context, int64, int64) error
	AddFollowerCount(context.Context, int64) error
	AddFollowCount(context.Context, int64) error
	Delete(context.Context, int64, int64) error
	ReduceFollowerCount(context.Context, int64) error
	ReduceFollowCount(context.Context, int64) error
}

type RelationRepository struct{}

func (r *RelationRepository) IsFollow(ctx context.Context, userId, toUserId int64) error {
	return DB.WithContext(ctx).First(&Relation{}, "user_id = ? and follower_id = ? and delete_time = 0", toUserId, userId).Error
}

func (r *RelationRepository) Create(ctx context.Context, userId, toUserId int64) error {
	return DB.WithContext(ctx).Create(&Relation{UserID: toUserId, FollowerID: userId}).Error
}

func (r *RelationRepository) AddFollowerCount(ctx context.Context, toUserId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", toUserId).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error
}

func (r *RelationRepository) AddFollowCount(ctx context.Context, userId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", userId).
		UpdateColumn("follow_count", gorm.Expr("follow_count + 1")).Error
}

func (r *RelationRepository) Delete(ctx context.Context, userId, toUserId int64) error {
	return DB.WithContext(ctx).Model(&Relation{}).Where("user_id=? AND follower_id=?", toUserId, userId).
		Update("delete_time", time.Now().Unix()).Error
}

func (r *RelationRepository) ReduceFollowerCount(ctx context.Context, toUserId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", toUserId).
		UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error
}

func (r *RelationRepository) ReduceFollowCount(ctx context.Context, userId int64) error {
	return DB.WithContext(ctx).Model(&User{}).Where("user_id = ? and delete_time = 0", userId).
		UpdateColumn("follow_count", gorm.Expr("follow_count - 1")).Error
}
