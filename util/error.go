package util

import "errors"

const (
	Success = iota

	AuthWrong
	NoAuth

	UserExist
	UserNotExist
	WrongPassword

	InternalServerError
	ParamError
)

var ErrCode2Msg = map[int]string{
	Success: "success",

	AuthWrong: "the token is expired or invalid",
	NoAuth:    "token not received",

	UserExist:     "user already exists",
	UserNotExist:  "username is not exist",
	WrongPassword: "password is wrong",

	InternalServerError: "internal server error",
	ParamError:          "something wrong with param...",
}

var (
	ErrUserExist = errors.New(ErrCode2Msg[UserExist])
)
