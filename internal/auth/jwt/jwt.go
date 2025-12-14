package jwt

import (
	"github.com/autumnterror/breezynotes/internal/auth/config"
)

type WithConfig struct {
	cfg *config.Config
}

func NewWithConfig(cfg *config.Config) *WithConfig {
	return &WithConfig{cfg: cfg}
}
