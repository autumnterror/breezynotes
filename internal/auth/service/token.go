package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"time"
)

func (s *AuthService) GenerateAccessToken(ctx context.Context, id string) (*domain.Token, error) {
	const op = "service.GenerateAccessToken"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	repo, err := s.tokenRepo()
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	token, err := repo.GenerateToken(id, jwt.TokenTypeAccess)
	if err != nil {
		return nil, wrapTokenError(op, err)
	}

	return &domain.Token{
		Value: token,
		Exp:   time.Now().UTC().Add(s.cfg.AccessTokenLifeTime),
	}, nil
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, id string) (*domain.Token, error) {
	const op = "service.GenerateRefreshToken"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	repo, err := s.tokenRepo()
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	token, err := repo.GenerateToken(id, jwt.TokenTypeRefresh)
	if err != nil {
		return nil, wrapTokenError(op, err)
	}

	return &domain.Token{
		Value: token,
		Exp:   time.Now().UTC().Add(s.cfg.RefreshTokenLifeTime),
	}, nil
}

func (s *AuthService) GenerateTokens(ctx context.Context, id string) (*domain.Tokens, error) {
	const op = "service.GenerateTokens"

	access, err := s.GenerateAccessToken(ctx, id)
	if err != nil {
		return nil, format.Error(op, err)
	}

	refresh, err := s.GenerateRefreshToken(ctx, id)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &domain.Tokens{
		Access:  *access,
		Refresh: *refresh,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, token string) (*domain.Token, error) {
	const op = "service.Refresh"
	if stringEmpty(token) {
		return nil, wrapServiceCheck(op, errors.New("token is empty"))
	}

	repo, err := s.tokenRepo()
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	newAccess, err := repo.Refresh(token)
	if err != nil {
		return nil, wrapTokenError(op, err)
	}

	return &domain.Token{
		Value: newAccess,
		Exp:   time.Now().UTC().Add(s.cfg.AccessTokenLifeTime),
	}, nil
}

func (s *AuthService) CheckToken(ctx context.Context, token string) error {
	const op = "service.CheckToken"
	_, _, err := s.parseToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetIdFromToken(ctx context.Context, token string) (string, error) {
	const op = "service.GetIdFromToken"

	raw, repo, err := s.parseAccessToken(ctx, token)
	if err != nil {
		return "", format.Error(op, err)
	}

	id, err := repo.GetIdFromToken(raw)
	if err != nil {
		return "", wrapTokenError(op, err)
	}

	if err := idValidation(id); err != nil {
		return "", wrapServiceCheck(op, err)
	}

	return id, nil
}

func (s *AuthService) GetUserDataFromToken(ctx context.Context, token string) (*domain.User, error) {
	const op = "service.GetUserDataFromToken"

	id, err := s.GetIdFromToken(ctx, token)
	if err != nil {
		return nil, format.Error(op, err)
	}

	repo, err := s.userRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	user, err := repo.GetInfo(ctx, id)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return user, nil
}

func (s *AuthService) parseToken(ctx context.Context, token string) (*jwtlib.Token, jwt.WithConfigRepo, error) {
	const op = "service.parseToken"
	if stringEmpty(token) {
		return nil, nil, wrapServiceCheck(op, errors.New("token is empty"))
	}

	repo, err := s.tokenRepo()
	if err != nil {
		return nil, nil, wrapServiceCheck(op, err)
	}

	raw, err := repo.VerifyToken(token)
	if err != nil {
		return nil, nil, wrapTokenError(op, err)
	}

	return raw, repo, nil
}

func (s *AuthService) parseAccessToken(ctx context.Context, token string) (*jwtlib.Token, jwt.WithConfigRepo, error) {
	const op = "service.parseAccessToken"

	raw, repo, err := s.parseToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	tp, err := repo.GetTypeFromToken(raw)
	if err != nil {
		return nil, nil, wrapTokenError(op, err)
	}
	if tp != jwt.TokenTypeAccess {
		return nil, nil, format.Error(op, domain.ErrTokenWrongType)
	}

	return raw, repo, nil
}

func wrapTokenError(op string, err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return format.Error(op, domain.ErrTokenExpired)
	case errors.Is(err, jwt.ErrWrongType):
		return format.Error(op, domain.ErrTokenWrongType)
	default:
		return format.Error(op, fmt.Errorf("%w: %v", domain.ErrTokenInvalid, err))
	}
}
