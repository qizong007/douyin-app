package repository

import (
	"context"
	"time"
)

type Favorite struct {
	FavoriteId int64 `json:"favorite_id" gorm:"primaryKey"`
	UserId     int64 `json:"user_id" gorm:"index"`
	VideoId    int64 `json:"video_id" gorm:"index"`
	CreateTime int64 `json:"create_time" gorm:"autoCreateTime;index"`
	ModifyTime int64 `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime int64 `json:"delete_time"`
}

type IFavoriteRepository interface {
	Create(context.Context, *Favorite) error
	Update(context.Context, *Favorite) error
	DeleteByUserIdAndVideoId(context.Context, int64, int64) error
	//FindByVideoIdAndUserId(context.Context, int64, int64) (*Favorite, error)
	FindVideoListByUseId(context.Context, int64) ([]*Favorite, error)
	FindFavoriteRecord(ctx context.Context, userId int64, videoId int64) (favorite *Favorite, err error)
	UpdateFavorite(ctx context.Context, userId int64, videoId int64, favorite *Favorite) (err error)
}

type FavoriteRepository struct{}

func (r *FavoriteRepository) FindVideoListByUseId(ctx context.Context, userId int64) ([]*Favorite, error) {
	favorite := make([]*Favorite, 0)
	err := DB.WithContext(ctx).Order("create_time desc").Where("user_id = ?", userId).Find(&favorite).Error
	return favorite, err
}

func (r *FavoriteRepository) FindFavoriteRecord(ctx context.Context, userId int64, videoId int64) (favorite *Favorite, err error) {
	err = DB.WithContext(ctx).Where("user_id = ? and video_id = ?", userId, videoId).First(&favorite).Error
	return favorite, err
}

func (r *FavoriteRepository) UpdateFavorite(ctx context.Context, userId int64, videoId int64, favorite *Favorite) (err error) {
	err = DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Update("delete_time", favorite.DeleteTime).Error
	return err
}

func (r *FavoriteRepository) Create(ctx context.Context, favorite *Favorite) error {
	return DB.WithContext(ctx).Create(&favorite).Error
}

func (r *FavoriteRepository) Update(ctx context.Context, favorite *Favorite) error {
	return DB.WithContext(ctx).Where("user_id = ? and ").Updates(&favorite).Error
}

func (r *FavoriteRepository) DeleteByUserIdAndVideoId(ctx context.Context, userId int64, videoId int64) error {
	return DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Update("delete_time", time.Now().Unix()).Error
}

//type IVideoRepository interface {
//	Create(context.Context, *Video) error
//	FindByUserId(context.Context, int64) ([]*Video, error)
//	FindByVideoId(context.Context, int64) (*Video, error)
//	FindWithLimit(context.Context, int) ([]*Video, error)
//	FindByCreateTimeWithLimit(context.Context, int64, int) ([]*Video, error)
//	VideoFavoriteAdd(context.Context, *Video, int64) error
//}
//type VideoRepository struct{}
//
//func (r *VideoRepository) Create(ctx context.Context, video *Video) error {
//	return DB.WithContext(ctx).Create(&video).Error
//}
//
//func (r *VideoRepository) FindByUserId(ctx context.Context, userId int64) ([]*Video, error) {
//	videos := make([]*Video, 0)
//	err := DB.WithContext(ctx).Order("create_time desc").Where("user_id = ? and delete_time = 0", userId).Find(&videos).Error
//	return videos, err
//}
//
//func (r *VideoRepository) FindByVideoId(ctx context.Context, videoId int64) (video *Video, err error) {
//	err = DB.WithContext(ctx).Where("video_id = ?", videoId).Find(&video).Error
//	return video, err
//}
//
//func (r *VideoRepository) FindWithLimit(ctx context.Context, limit int) ([]*Video, error) {
//	videos := make([]*Video, 0)
//	err := DB.WithContext(ctx).Order("create_time desc").Limit(limit).Where("delete_time = 0").Find(&videos).Error
//	return videos, err
//}
//
//func (r *VideoRepository) FindByCreateTimeWithLimit(ctx context.Context, createTime int64, limit int) ([]*Video, error) {
//	videos := make([]*Video, 0)
//	err := DB.WithContext(ctx).Order("create_time desc").Limit(limit).Where("create_time < ? and delete_time = 0", createTime).Find(&videos).Error
//	return videos, err
//}
//
//func (r *VideoRepository) VideoFavoriteAdd(ctx context.Context, video *Video, videoId int64) (err error) {
//	err = DB.WithContext(ctx).Where("video_id = ? ", videoId).Updates(video).Error
//	return err
//}
