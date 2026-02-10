package domain2

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrTypeNotDefined = errors.New("need to register type")
	ErrAlreadyUsed    = errors.New("block already in use")
	ErrBadRequest     = errors.New("bad fields")
	ErrUnauthorized   = errors.New("you dont have permission")
)
