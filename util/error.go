package util

const (
	Success = iota

	ErrorAuth
	NoAuth

	InternalServerError
	ParamError
)

var ErrCode2Msg = map[int]string{
	Success: "success",

	ErrorAuth: "the token is expired or invalid",
	NoAuth:    "token not received",

	InternalServerError: "internal server error",
	ParamError:          "something wrong with param...",
}
