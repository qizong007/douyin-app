package util

import "errors"

const (
	Success = iota

	WrongAuth
	NoAuth

	UserExisted
	UserNotExist
	WrongPassword

	InternalServerError
	ParamError
)

var ErrCode2Msg = map[int]string{
	Success: "success",

	WrongAuth: "the token is expired or invalid",
	NoAuth:    "token not received",

	UserExisted:   "user already exists",
	UserNotExist:  "username is not exist",
	WrongPassword: "password is wrong",

	InternalServerError: "internal server error",
	ParamError:          "something wrong with param...",
}

var (
	ErrWrongAuth     = errors.New(ErrCode2Msg[WrongAuth])
	ErrNoAuth        = errors.New(ErrCode2Msg[NoAuth])
	ErrUserExisted   = errors.New(ErrCode2Msg[UserExisted])
	ErrWrongPassword = errors.New(ErrCode2Msg[WrongPassword])
)
