package service

import (
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
)

var (
	ErrBadServiceCheck = errors.New("bad service check")
)

func wrapServiceCheck(op string, err error) error {
	if err == nil {
		return nil
	}
	return format.Error(op, fmt.Errorf("%w: %v", ErrBadServiceCheck, err))
}

func idValidation(id string) error {
	if uid.Validate(id) {
		return nil
	}
	return errors.New("id not in uuid")
}

func stringEmpty(s string) bool {
	if s == "" {
		return true
	}
	return false
}

func userValidation(u *domain.User) error {
	if err := idValidation(u.Id); err != nil {
		return err
	}
	if stringEmpty(u.Login) {
		return errors.New("login is empty")
	}
	if stringEmpty(u.Email) {
		return errors.New("email is empty")
	}
	if stringEmpty(u.Password) {
		return errors.New("pw is empty")
	}
	if stringEmpty(u.Photo) {
		return errors.New("photo is empty")
	}

	return nil
}
