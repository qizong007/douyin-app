package util

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func initValidate() {
	Validate = validator.New()
}

func InitUtil() {
	initIdGenerator()
	initJWTVal()
	initValidate()
}
