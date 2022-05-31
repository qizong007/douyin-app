package util

import "errors"

const (
	Success = iota

	WrongAuth
	NoAuth

	UserExisted
	UserNotExist
	WrongPassword
	VideoNotExist
	CommentNotExist

	SensitiveComment
	IsFollow
	NotFollow
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

	IsFollow:        "该关注记录已存在",
	NotFollow:       "该关注记录尚未存在",
	VideoNotExist:   "video is not exist",
	CommentNotExist: "commentId is not exist",

	SensitiveComment:    "Comment contains sensitive words",
	InternalServerError: "internal server error",
	ParamError:          "something wrong with param...",
}

var (
	ErrWrongAuth        = errors.New(ErrCode2Msg[WrongAuth])
	ErrNoAuth           = errors.New(ErrCode2Msg[NoAuth])
	ErrUserExisted      = errors.New(ErrCode2Msg[UserExisted])
	ErrWrongPassword    = errors.New(ErrCode2Msg[WrongPassword])
	ErrIsFollow         = errors.New(ErrCode2Msg[IsFollow])
	ErrNotFollow        = errors.New(ErrCode2Msg[NotFollow])
	ErrSensitiveComment = errors.New(ErrCode2Msg[SensitiveComment])
	ErrCommentNotExist  = errors.New(ErrCode2Msg[CommentNotExist])
)
