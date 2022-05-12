package repository

import (
	"context"
)

type User struct {
	Id         int64  `json:"id" gorm:"primaryKey"`
	UserId     int64  `json:"user_id" gorm:"index"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime int64  `json:"modify_time" gorm:"autoUpdateTime"`
	DeleteTime int64  `json:"delete_time"`
}

type IUserRepository interface {
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	DeleteByUserId(context.Context, int64) error
	FindByUserId(context.Context, int64) (*User, error)
	FindByUserIds(context.Context, []int64) ([]*User, error)
	FindByUsername(context.Context, string) (*User, error)
}

type UserRepository struct{}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return DB.WithContext(ctx).Create(&user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return DB.WithContext(ctx).Where("delete_time = 0").Updates(&user).Error
}

func (r *UserRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	return DB.WithContext(ctx).Where("user_id = ? and delete_time = 0", userId).Delete(User{}).Error
}

func (r *UserRepository) FindByUserId(ctx context.Context, userId int64) (user *User, err error) {
	err = DB.WithContext(ctx).Where("user_id = ? and delete_time = 0", userId).First(&user).Error
	return user, err
}

func (r *UserRepository) FindByUserIds(ctx context.Context, userIdList []int64) ([]*User, error) {
	users := make([]*User, 0)
	err := DB.WithContext(ctx).Where("user_id in (?) and delete_time = 0", userIdList).Find(&users).Error
	return users, err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (user *User, err error) {
	err = DB.WithContext(ctx).Where("username = ? and delete_time = 0", username).First(&user).Error
	return user, err
}
