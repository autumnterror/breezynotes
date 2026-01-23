package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/autumnterror/utils_go/pkg/utils/validate"
)

func (s *AuthService) Auth(ctx context.Context, email, login, pw string) (string, string, error) {
	const op = "service.Auth"
	if stringEmpty(email) && stringEmpty(login) {
		return "", "", wrapServiceCheck(op, errors.New("email and login is empty"))
	}
	if stringEmpty(pw) {
		return "", "", wrapServiceCheck(op, errors.New("pw is empty"))
	}

	repo, err := s.authRepo(ctx)
	if err != nil {
		return "", "", wrapServiceCheck(op, err)
	}
	repoJwt, err := s.tokenRepo()
	if err != nil {
		return "", "", wrapServiceCheck(op, err)
	}

	id, err := repo.Authentication(ctx, email, login, pw)
	if err != nil {
		return "", "", err
	}

	at, err := repoJwt.GenerateToken(id, domain.TokenTypeAccess)
	if err != nil {
		return "", "", err
	}
	rt, err := repoJwt.GenerateToken(id, domain.TokenTypeRefresh)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

func (s *AuthService) Reg(ctx context.Context, email, login, pw string) (string, string, error) {
	const op = "service.Reg"
	if stringEmpty(email) && stringEmpty(login) {
		return "", "", wrapServiceCheck(op, errors.New("email and login is empty"))
	}
	if stringEmpty(pw) {
		return "", "", wrapServiceCheck(op, errors.New("pw is empty"))
	}
	if !validate.Password(pw) {
		return "", "", wrapServiceCheck(op, errors.New("pw not in policy"))
	}
	repo, err := s.userRepo(ctx)
	if err != nil {
		return "", "", wrapServiceCheck(op, err)
	}
	repoJwt, err := s.tokenRepo()
	if err != nil {
		return "", "", wrapServiceCheck(op, err)
	}

	id := uid.New()
	if err := repo.Create(ctx, &domain.User{
		Id:       id,
		Login:    login,
		Email:    email,
		About:    "😌",
		Photo:    "images/default.png",
		Password: pw,
	}); err != nil {
		return "", "", err
	}

	at, err := repoJwt.GenerateToken(id, domain.TokenTypeAccess)
	if err != nil {
		return "", "", err
	}
	rt, err := repoJwt.GenerateToken(id, domain.TokenTypeRefresh)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

func (s *AuthService) ValidateTokens(ctx context.Context, at, rt string) (string, error) {
	const op = "service.Reg"
	if stringEmpty(at) || stringEmpty(rt) {
		return "", wrapServiceCheck(op, errors.New("one of tokens is empty"))
	}

	repoJwt, err := s.tokenRepo()
	if err != nil {
		return "", wrapServiceCheck(op, err)
	}

	_, err = repoJwt.VerifyToken(at)
	if err != nil {
		if errors.Is(err, domain.ErrTokenExpired) {
			rt, err := repoJwt.Refresh(rt)
			if err != nil {
				return "", err
			}
			return rt, nil
		} else {
			return "", err
		}
	}

	return "", nil
}
