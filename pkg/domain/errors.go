package domain

import "errors"

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrBadRequest      = errors.New("bad request")
)
