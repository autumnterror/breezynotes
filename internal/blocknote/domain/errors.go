package domain

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not fiend")
	ErrTypeNotDefined = errors.New("need to register type")
	ErrAlreadyUsed    = errors.New("block already in use")
	ErrBadRequest     = errors.New("bad fields")
)
