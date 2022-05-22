package repository

import "context"

type Video struct {
	VideoId       int64  `json:"video_id" gorm:"primaryKey"`
	UserId        int64  `json:"user_id" gorm:"index"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	CreateTime    int64  `json:"create_time" gorm:"autoCreateTime;index"`
	ModifyTime    int64  `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime    int64  `json:"delete_time"`
}

type IVideoRepository interface {
	Create(context.Context, *Video) error
	FindByUserId(context.Context, int64) ([]*Video, error)
	FindWithLimit(context.Context, int) ([]*Video, error)
	FindByCreateTimeWithLimit(context.Context, int64, int) ([]*Video, error)
}
type VideoRepository struct{}

func (r *VideoRepository) Create(ctx context.Context, video *Video) error {
	return DB.WithContext(ctx).Create(&video).Error
}

func (r *VideoRepository) FindByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Where("user_id = ? and delete_time = 0", userId).Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) FindWithLimit(ctx context.Context, limit int) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Limit(limit).Where("delete_time = 0").Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) FindByCreateTimeWithLimit(ctx context.Context, createTime int64, limit int) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Limit(limit).Where("create_time < ? and delete_time = 0", createTime).Find(&videos).Error
	return videos, err
}
