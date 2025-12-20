package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/repository"
)

func (s *AuthService) runInTx(ctx context.Context, op string, fn func(ctx context.Context) error) error {
	if s.tx == nil {
		return wrapServiceCheck(op, errors.New("tx runner is nil"))
	}
	return s.tx.RunInTx(ctx, fn)
}

func (s *AuthService) userRepo(ctx context.Context) (repository.UserRepo, error) {
	repo, err := s.repoGetter(ctx, func(p repository.Provider) any {
		return p.User(ctx)
	})
	if err != nil {
		return nil, err
	}
	res, _ := repo.(repository.UserRepo)
	if res == nil {
		return nil, errors.New("user repository is nil")
	}
	return res, nil
}

func (s *AuthService) authRepo(ctx context.Context) (repository.AuthRepo, error) {
	repo, err := s.repoGetter(ctx, func(p repository.Provider) any {
		return p.Auth(ctx)
	})
	if err != nil {
		return nil, err
	}
	res, _ := repo.(repository.AuthRepo)
	if res == nil {
		return nil, errors.New("auth repository is nil")
	}
	return res, nil
}

func (s *AuthService) healthRepo(ctx context.Context) (repository.HealthRepo, error) {
	repo, err := s.repoGetter(ctx, func(p repository.Provider) any {
		return p.Health(ctx)
	})
	if err != nil {
		return nil, err
	}
	res, _ := repo.(repository.HealthRepo)
	if res == nil {
		return nil, errors.New("health repository is nil")
	}
	return res, nil
}

func (s *AuthService) repoGetter(ctx context.Context, getter func(repository.Provider) any) (any, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	if getter == nil {
		return nil, errors.New("repository provider is nil")
	}
	return getter(s.repos), nil
}

func (s *AuthService) tokenRepo() (jwt.WithConfigRepo, error) {
	if s.tokens == nil {
		return nil, errors.New("token repository is nil")
	}
	if s.cfg == nil {
		return nil, errors.New("config is nil")
	}
	return s.tokens, nil
}
