package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

func (s *NotesService) CleanTrash(ctx context.Context, uid string) error {
	const op = "service.CleanTrash"
	if err := idValidation(uid); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.CleanTrash(ctx, uid)
	})

	return err
}
func (s *NotesService) GetNotesFromTrash(ctx context.Context, uid string) (*domain.NoteParts, error) {
	const op = "service.GetNotesFromTrash"
	if err := idValidation(uid); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	res, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return s.nts.GetNotesFromTrash(ctx, uid)
	})
	if err != nil {
		return nil, err
	}

	if resS, ok := res.(*domain.NoteParts); !ok {
		return nil, wrapServiceCheck(op, errors.New("response type mismatch"))
	} else {
		return resS, nil
	}
}
func (s *NotesService) ToTrash(ctx context.Context, id string) error {
	const op = "service.ToTrash"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.ToTrash(ctx, id)
	})

	return err
}
func (s *NotesService) FromTrash(ctx context.Context, id string) error {
	const op = "service.FromTrash"

	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.FromTrash(ctx, id)
	})

	return err
}
func (s *NotesService) FindOnTrash(ctx context.Context, id string) (*domain.Note, error) {
	const op = "service.FindOnTrash"

	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return s.nts.FindOnTrash(ctx, id)
}
