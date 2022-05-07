package repository

type UserInfo struct {
	Id int
	//ServiceId int
	Name       string
	Password   string
	CreateTime int64
	ModifyTime int64
	DeleteTime int64
}
