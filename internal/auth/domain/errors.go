package domain

import "errors"

var (
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrWrongInput        = errors.New("wrong input")
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("obj already exist")
	ErrForeignKey        = errors.New("sub obj dont exist")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenWrongType    = errors.New("token wrong type")
	ErrTokenInvalid      = errors.New("token invalid")
)
