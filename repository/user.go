package repository

import (
	"douyin-app/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         int64  `json:"id"`      //逻辑id,自增
	UserId     int64  `json:"user_id"` //业务id
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreateTime int64  `json:"create_time"`
	ModifyTime int64  `json:"modify_time"`
	DeleteTime int64  `json:"delete_time"`
}

//在创建前设置时间
func (user *User) BeforeCreate(*gorm.DB) error {
	user.CreateTime = time.Now().Unix()
	return nil
}

//判断用户名在数据库中是否存在,存在则返回true
func ExistUserByName(username string) bool {
	var user User
	err := DB.Select("id").Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func CreateUserInfo(username string, password string) (ID int64, UserId int64, err error) {
	user := User{
		UserId:   util.GenerateId(),
		Username: username,
		Password: password,
	}

	err = DB.Create(&user).Error
	if err != nil {
		return 0, 0, err
	}
	return user.Id, user.UserId, nil
}

//通过username获取用户信息
//并验证用户的密码是否正确
func VerifyPassword(username string, password string) (ID int64, UserId int64, ok bool) {
	var user User
	DB.Where("username = ?", username).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, 0, false
	}
	return user.Id, user.UserId, true
}
