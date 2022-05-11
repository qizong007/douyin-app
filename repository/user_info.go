package repository

import (
	"douyin-app/util"
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	ID         int    //逻辑id,自增
	UserId     int64  `json:"user_id"` //业务id
	Username   string `json:"username"`
	Password   string
	CreateTime int64
	ModifyTime int64
	DeleteTime int64
}

//在创建前设置时间
func (user *UserInfo) BeforeCreate(*gorm.DB) error {
	user.CreateTime = time.Now().Unix()
	return nil
}

//判断用户名在数据库中是否存在,存在则返回true
func ExistUserByName(username string) bool {
	var user UserInfo
	DB.Select("id").Where("username = ?", username).First(&user)
	if user.ID > 0 {
		return true
	}
	return false
}

func CreateUserInfo(username string, password string) (ID int, userId int64, err error) {
	user := UserInfo{
		UserId:   util.GenerateId(),
		Username: username,
		Password: password,
	}

	err = DB.Create(&user).Error
	if err != nil {
		return 0, 0, err
	}
	return user.ID, user.UserId, nil
}
