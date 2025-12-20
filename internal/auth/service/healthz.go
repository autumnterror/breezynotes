package service

import (
	"context"
)

func (s *AuthService) Health(ctx context.Context) error {
	const op = "service.Health"

	repo, err := s.healthRepo(ctx)
	if err != nil {
		return wrapServiceCheck(op, err)
	}
	return repo.Healthz(ctx)
}
