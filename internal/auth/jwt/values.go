package jwt

import "errors"

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrWrongType    = errors.New("wrong type of token")
)

const (
	TokenTypeAccess  = "ACCESS"
	TokenTypeRefresh = "REFRESH"
)
