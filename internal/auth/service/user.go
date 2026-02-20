package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/utils_go/pkg/utils/validate"
)

func (s *AuthService) CreateAdmin(ctx context.Context) error {
	const op = "service.CreateAdmin"

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		_, err = repo.CreateAdmin(ctx)
		return err
	})
}

func (s *AuthService) UpdatePassword(ctx context.Context, id, oldPassword string, newPassword string) error {
	const op = "service.UpdatePassword"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	if !validate.Password(newPassword) {
		return wrapServiceCheck(op, errors.New("new password not in policy"))
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return err
		}
		repoAuth, err := s.authRepo(ctx)
		if err != nil {
			return err
		}

		info, err := repo.GetInfo(ctx, id)
		if err != nil {
			return err
		}
		if _, err := repoAuth.Authentication(ctx, "", info.Login, oldPassword); err != nil {
			return domain.ErrNotFound
		}

		return repo.UpdatePassword(ctx, id, newPassword)
	})
}
func (s *AuthService) UpdatePhoto(ctx context.Context, id, np string) error {
	const op = "service.UpdatePhoto"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(np) {
		np = "images/default.png"
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return err
		}
		return repo.UpdatePhoto(ctx, id, np)
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
			return err
		}
		return repo.UpdateEmail(ctx, id, email)
	})
}
func (s *AuthService) UpdateAbout(ctx context.Context, id, about string) error {
	const op = "service.UpdateAbout"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return err
		}
		return repo.UpdateAbout(ctx, id, about)
	})
}
func (s *AuthService) Delete(ctx context.Context, id string) error {
	const op = "service.delete"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.userRepo(ctx)
		if err != nil {
			return err
		}
		return repo.Delete(ctx, id)
	})
}

func (s *AuthService) GetIdFromLogin(ctx context.Context, login string) (string, error) {
	const op = "service.Create"
	if stringEmpty(login) {
		return "", wrapServiceCheck(op, errors.New("login is empty"))
	}

	repo, err := s.userRepo(ctx)
	if err != nil {
		return "", wrapServiceCheck(op, err)
	}
	return repo.GetIdFromLogin(ctx, login)
}

func (s *AuthService) GetUserDataFromToken(ctx context.Context, token string) (*domain.User, error) {
	const op = "service.GetUserDataFromToken"

	id, err := s.GetIdFromToken(ctx, token)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	repo, err := s.userRepo(ctx)
	if err != nil {
		return nil, err
	}

	user, err := repo.GetInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
