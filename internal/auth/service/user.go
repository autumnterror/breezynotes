package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
)

func (s *AuthService) GetInfo(ctx context.Context, id string) (*domain.User, error) {
	const op = "service.GetInfo"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.userRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetInfo(ctx, id)
}

func (s *AuthService) Create(ctx context.Context, u *domain.User) error {
	const op = "service.Create"
	if err := userValidation(u); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.Create(ctx, u)
	})
}
func (s *AuthService) UpdatePhoto(ctx context.Context, id, np string) error {
	const op = "service.UpdatePhoto"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(np) {
		return wrapServiceCheck(op, errors.New("new photo is empty"))
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdatePhoto(ctx, id, np)
	})
}
func (s *AuthService) UpdatePassword(ctx context.Context, id, newPassword string) error {
	const op = "service.UpdatePassword"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(newPassword) {
		return wrapServiceCheck(op, errors.New("new password is empty"))
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdatePassword(ctx, id, newPassword)
	})
}
func (s *AuthService) UpdateEmail(ctx context.Context, id, email string) error {
	const op = "service.UpdateEmail"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(email) {
		return wrapServiceCheck(op, errors.New("new email is empty"))
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateEmail(ctx, id, email)
	})
}
func (s *AuthService) UpdateAbout(ctx context.Context, id, about string) error {
	const op = "service.UpdateAbout"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(about) {
		return wrapServiceCheck(op, errors.New("new about is empty"))
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateAbout(ctx, id, about)
	})
}
func (s *AuthService) Delete(ctx context.Context, id string) error {
	const op = "service.Delete"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.Delete(ctx, id)
	})
}
