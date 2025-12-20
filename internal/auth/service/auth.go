package service

import (
	"context"
	"errors"
)

func (s *AuthService) Auth(ctx context.Context, email, login, pw string) (string, error) {
	const op = "service.Auth"
	if stringEmpty(email) && stringEmpty(login) {
		return "", wrapServiceCheck(op, errors.New("email and login is empty"))
	}
	if stringEmpty(pw) {
		return "", wrapServiceCheck(op, errors.New("pw is empty"))
	}

	repo, err := s.authRepo(ctx)
	if err != nil {
		return "", wrapServiceCheck(op, err)
	}

	return repo.Authentication(ctx, email, login, pw)
}
