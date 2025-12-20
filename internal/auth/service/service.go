package service

import (
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/repository"
)

type AuthService struct {
	tx     TxRunner
	repos  repository.Provider
	tokens jwt.WithConfigRepo
	cfg    *config.Config
}

func NewAuthService(
	tx TxRunner,
	repos repository.Provider,
	tokens jwt.WithConfigRepo,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		tx:     tx,
		repos:  repos,
		tokens: tokens,
		cfg:    cfg,
	}
}
