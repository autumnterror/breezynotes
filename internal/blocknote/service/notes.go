package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
)

func (s *BN) CreateNote(ctx context.Context, n *domain.Note) error {
	const op = "service.CreateNote"
	if err := noteValidation(n); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.Create(ctx, n)
	})

	return err
}

func (s *BN) GetNote(ctx context.Context, idNote, idUser string) (*domain.NoteWithBlocks, error) {
	const op = "service.Get"
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

	if n.Author != idUser || !alg.IsIn(idUser, n.Editors) || !alg.IsIn(idUser, n.Readers) {
		return nil, domain.ErrUnauthorized
	}

	blks, err := s.blk.GetMany(ctx, n.Blocks)
	if err != nil {
		return nil, err
	}

	return &domain.NoteWithBlocks{
		Id:        n.Id,
		Title:     n.Title,
		Blocks:    blks.Blks,
		Author:    n.Author,
		Readers:   n.Readers,
		Editors:   n.Editors,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       n.Tag,
	}, nil
}

func (s *BN) GetNoteListByUser(ctx context.Context, idUser string) (*domain.NoteParts, error) {
	const op = "service.GetNoteListByUser"
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNoteListByUser(ctx, idUser)
}

func (s *BN) GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain.NoteParts, error) {
	const op = "service.GetNoteListByTag"
	if err := idValidation(idTag); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNoteListByTag(ctx, idTag, idUser)
}

func (s *BN) AddTagToNote(ctx context.Context, noteId, tagId, idUser string) error {
	const op = "service.AddTagToNote"
	if err := idValidation(tagId); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(noteId); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, noteId)
		if err != nil {
			return nil, err
		}
		if n.Author != idUser || !alg.IsIn(idUser, n.Editors) || !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		tag, err := s.tgs.Get(ctx, tagId)
		if err != nil {
			return nil, err
		}
		return nil, s.nts.AddTagToNote(ctx, noteId, tag)
	})

	return err
}

func (s *BN) RemoveTagFromNote(ctx context.Context, noteId string, tagId, idUser string) error {
	const op = "service.AddTagToNote"

	if err := idValidation(noteId); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(tagId); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, noteId)
		if err != nil {
			return nil, err
		}
		if n.Author != idUser || !alg.IsIn(idUser, n.Editors) || !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.RemoveTagFromNote(ctx, noteId, tagId)
	})

	return err
}

func (s *BN) UpdateTitleNote(ctx context.Context, idNote, idUser, nTitle string) error {
	const op = "service.UpdateTitleNote"

	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	if stringEmpty(nTitle) {
		return wrapServiceCheck(op, errors.New("title empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote)
		if err != nil {
			return nil, err
		}
		if n.Author != idUser || !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.UpdateTitle(ctx, idNote, nTitle)
	})

	return err
}
