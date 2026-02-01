package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
)

func (s *BN) ShareNote(ctx context.Context, noteId, userId, userIdToShare, role string) error {
	const op = "service.ShareNote"

	if stringEmpty(role) {
		return wrapServiceCheck(op, errors.New("role is empty"))
	}
	if idValidation(noteId) != nil {
		return wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(userId) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}
	if idValidation(userIdToShare) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}
	if userId == userIdToShare {
		return wrapServiceCheck(op, errors.New("can't share for yourself"))
	}

	switch role {
	case domain.EditorRole:
	case domain.ReaderRole:
	default:
		return wrapServiceCheck(op, errors.New("role undefined"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, noteId); err != nil {
			return nil, err
		} else if n.Author != userId && !alg.IsIn(userId, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.ShareNote(ctx, noteId, userIdToShare, role)
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *BN) ChangeUserRole(ctx context.Context, noteId, userId, userIdToChange, newRole string) error {
	const op = "service.ChangeUserRole"
	if stringEmpty(newRole) {
		return wrapServiceCheck(op, errors.New("newRole is empty"))
	}
	if idValidation(noteId) != nil {
		return wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(userId) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}
	if idValidation(userIdToChange) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}
	if userId == userIdToChange {
		return wrapServiceCheck(op, errors.New("can't change for yourself"))
	}

	switch newRole {
	case domain.EditorRole:
	case domain.ReaderRole:
	default:
		return wrapServiceCheck(op, errors.New("role undefined"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, noteId); err != nil {
			return nil, err
		} else if n.Author != userId && !alg.IsIn(userId, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.ChangeUserRole(ctx, noteId, userIdToChange, newRole)
	})

	if err != nil {
		return err
	}
	return nil
}
