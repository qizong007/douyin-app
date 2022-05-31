package repository

import (
	"context"
	"douyin-app/conf"
	"douyin-app/util"
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"testing"
)

const (
	confPath = "../conf/conf.yaml"
)

func TestUserRepo(t *testing.T) {
	conf.InitConf(confPath)
	initMySQL()
	initRepository()
	util.InitIdGenerator()

	ctx := context.Background()
	userId := util.GenerateId()
	newUser := &User{
		UserId:   userId,
		Username: "qizong007",
		Password: "123456",
	}

	if err := GetUserRepository().Create(ctx, newUser); err != nil {
		log.Println(err)
		return
	}
	log.Println(newUser)

	user, err := GetUserRepository().FindByUserId(ctx, userId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(user)

	user.Password = "6666"
	if err = GetUserRepository().Update(ctx, user); err != nil {
		log.Println(err)
		return
	}
	log.Println(user)

	if err = GetUserRepository().DeleteByUserId(ctx, userId); err != nil {
		log.Println(err)
		return
	}

	_, err = GetUserRepository().FindByUserId(ctx, userId)
	assert.Equal(t, true, errors.Is(err, gorm.ErrRecordNotFound))

	id1 := util.GenerateId()
	id2 := util.GenerateId()
	id3 := util.GenerateId()

	user1 := &User{
		UserId:   id1,
		Username: "user1",
		Password: "123456",
	}

	user2 := &User{
		UserId:   id2,
		Username: "user2",
		Password: "123456",
	}

	user3 := &User{
		UserId:   id3,
		Username: "user3",
		Password: "123456",
	}

	_ = GetUserRepository().Create(ctx, user1)
	_ = GetUserRepository().Create(ctx, user2)
	_ = GetUserRepository().Create(ctx, user3)

	users, err := GetUserRepository().FindByUserIds(ctx, []int64{id1, id2, id3})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(users)

	user, err = GetUserRepository().FindByUsername(ctx, "user2")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(user)
}

//这个方法大多数情况可以达到更好的性能
func Benchmark_FindByUserIds(b *testing.B) {
	conf.InitConf(confPath)
	initMySQL()
	initRepository()
	ctx := context.Background()
	var ids []int64
	for k := 0; k < 5; k++ {
		ids = append(ids, 6774665417138176)
	}
	for k := 0; k < 5; k++ {
		ids = append(ids, 6774736186019840)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUserRepository().FindByUserIds(ctx, ids)
	}
}

func Benchmark_myFindByIds(b *testing.B) {
	conf.InitConf(confPath)
	initMySQL()
	initRepository()
	ctx := context.Background()
	var ids []int64
	for k := 0; k < 5; k++ {
		ids = append(ids, 6774665417138176)
	}
	for k := 0; k < 5; k++ {
		ids = append(ids, 6774736186019840)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFindByUserIds(ctx, ids)
	}
}

//low performance
func testFindByUserIds(ctx context.Context, userIdList []int64) ([]*User, error) {
	users := make([]*User, len(userIdList))

	for i := range userIdList {
		if i-1 > 0 && userIdList[i] == userIdList[i-1] {
			users[i] = users[i-1]
		} else {
			err := DB.WithContext(ctx).Where("user_id = (?) and delete_time = 0", userIdList[i]).Find(&users[i]).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return users, nil
}
