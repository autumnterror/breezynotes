package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
)

func (s *BN) CleanTrash(ctx context.Context, uid string) error {
	const op = "service.CleanTrash"
	if err := idValidation(uid); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.CleanTrash(ctx, uid)
	})

	return err
}
func (s *BN) GetNotesFromTrash(ctx context.Context, uid string) (*domain.NoteParts, error) {
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
func (s *BN) ToTrash(ctx context.Context, idNote, idUser string) error {
	const op = "service.ToTrash"
	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote)
		if err != nil {
			return nil, err
		}
		if n.Author != idUser || !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.ToTrash(ctx, idNote)
	})

	return err
}
func (s *BN) FromTrash(ctx context.Context, idNote, idUser string) error {
	const op = "service.FromTrash"

	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote)
		if err != nil {
			return nil, err
		}
		if n.Author != idUser || !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.FromTrash(ctx, idNote)
	})

	return err
}
func (s *BN) FindOnTrash(ctx context.Context, idNote, idUser string) (*domain.Note, error) {
	const op = "service.FindOnTrash"

	if err := idValidation(idNote); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	n, err := s.nts.Get(ctx, idNote)
	if err != nil {
		return nil, err
	}
	if n.Author != idUser || !alg.IsIn(idUser, n.Editors) {
		return nil, domain.ErrUnauthorized
	}

	return s.nts.FindOnTrash(ctx, idNote)
}
