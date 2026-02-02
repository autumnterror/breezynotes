package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) GetIdFromToken(ctx context.Context, token string) (string, error) {
	const op = "service.GetIdFromToken"

	raw, repo, err := s.parseAccessToken(ctx, token)
	if err != nil {
		return "", format.Error(op, err)
	}

	id, err := repo.GetIdFromToken(raw)
	if err != nil {
		return "", err
	}

	if err := idValidation(id); err != nil {
		return "", wrapServiceCheck(op, err)
	}

	return id, nil
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
		return nil, nil, err
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
		return nil, nil, err
	}
	if tp != domain.TokenTypeAccess {
		return nil, nil, wrapServiceCheck(op, domain.ErrTokenWrongType)
	}

	return raw, repo, nil
}

//func (s *AuthService) GenerateAccessToken(ctx context.Context, id string) (*domain.Token, error) {
//	const op = "service.GenerateAccessToken"
//	if err := idValidation(id); err != nil {
//		return nil, wrapServiceCheck(op, err)
//	}
//
//	repo, err := s.tokenRepo()
//	if err != nil {
//		return nil, wrapServiceCheck(op, err)
//	}
//
//	token, err := repo.GenerateToken(id, domain.TokenTypeAccess)
//	if err != nil {
//		return nil, err
//	}
//
//	return &domain.Token{
//		Value: token,
//		Exp:   time.Now().UTC().Add(s.cfg.AccessTokenLifeTime),
//	}, nil
//}
//
//func (s *AuthService) GenerateRefreshToken(ctx context.Context, id string) (*domain.Token, error) {
//	const op = "service.GenerateRefreshToken"
//	if err := idValidation(id); err != nil {
//		return nil, wrapServiceCheck(op, err)
//	}
//
//	repo, err := s.tokenRepo()
//	if err != nil {
//		return nil, wrapServiceCheck(op, err)
//	}
//
//	token, err := repo.GenerateToken(id, domain.TokenTypeRefresh)
//	if err != nil {
//		return nil, err
//	}
//
//	return &domain.Token{
//		Value: token,
//		Exp:   time.Now().UTC().Add(s.cfg.RefreshTokenLifeTime),
//	}, nil
//}
//
//func (s *AuthService) GenerateTokens(ctx context.Context, id string) (*domain.Tokens, error) {
//	const op = "service.GenerateTokens"
//
//	access, err := s.GenerateAccessToken(ctx, id)
//	if err != nil {
//		return nil, format.Error(op, err)
//	}
//
//	refresh, err := s.GenerateRefreshToken(ctx, id)
//	if err != nil {
//		return nil, format.Error(op, err)
//	}
//
//	return &domain.Tokens{
//		Access:  *access,
//		Refresh: *refresh,
//	}, nil
//}
//
//func (s *AuthService) Refresh(ctx context.Context, token string) (*domain.Token, error) {
//	const op = "service.Refresh"
//	if stringEmpty(token) {
//		return nil, wrapServiceCheck(op, errors.New("token is empty"))
//	}
//
//	repo, err := s.tokenRepo()
//	if err != nil {
//		return nil, wrapServiceCheck(op, err)
//	}
//
//	newAccess, err := repo.Refresh(token)
//	if err != nil {
//		return nil, err
//	}
//
//	return &domain.Token{
//		Value: newAccess,
//		Exp:   time.Now().UTC().Add(s.cfg.AccessTokenLifeTime),
//	}, nil
//}
//
//func (s *AuthService) CheckToken(ctx context.Context, token string) error {
//	const op = "service.CheckToken"
//	_, _, err := s.parseToken(ctx, token)
//	if err != nil {
//		return err
//	}
//	return nil
//}
