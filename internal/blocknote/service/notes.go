package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

func (s *NotesService) Create(ctx context.Context, n *domain.Note) error {
	const op = "service.CreateNote"
	if err := noteValidation(n); err != nil {
		return wrapServiceCheck(op, err)
	}
	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.Create(ctx, n)
	})

	return err
}

func (s *NotesService) GetNote(ctx context.Context, id string) (*domain.Note, error) {
	const op = "service.GetNote"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNote(ctx, id)
}

func (s *NotesService) GetNoteListByUser(ctx context.Context, id string) (*domain.NoteParts, error) {
	const op = "service.GetNoteListByUser"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNoteListByUser(ctx, id)
}

func (s *NotesService) GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain.NoteParts, error) {
	const op = "service.GetNoteListByTag"
	if err := idValidation(idTag); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNoteListByTag(ctx, idTag, idUser)
}

func (s *NotesService) AddTagToNote(ctx context.Context, noteId, tagId string) error {
	const op = "service.AddTagToNote"
	if err := idValidation(tagId); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(noteId); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		tag, err := s.tgs.Get(ctx, tagId)
		if err != nil {
			return nil, err
		}
		return nil, s.nts.AddTagToNote(ctx, noteId, tag)
	})

	return err
}

func (s *NotesService) RemoveTagFromNote(ctx context.Context, id string, tagId string) error {
	const op = "service.AddTagToNote"

	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(tagId); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.RemoveTagFromNote(ctx, id, tagId)
	})

	return err
}

func (s *NotesService) CreateBlock(ctx context.Context, _type, noteId string, data map[string]any, pos int) (string, error) {
	const op = "service.CreateBlock"

	if pos < 0 {
		return "", wrapServiceCheck(op, errors.New("pos < 0"))
	}

	res, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		blockId, err := s.blk.Create(ctx, _type, noteId, data)
		if err != nil {
			return nil, err
		}

		return blockId, s.nts.InsertBlock(ctx, noteId, blockId, pos)
	})

	if err != nil {
		return "", err
	}

	if resS, ok := res.(string); !ok {
		return "", wrapServiceCheck(op, errors.New("response type mismatch"))
	} else {
		return resS, nil
	}
}

func (s *NotesService) DeleteBlock(ctx context.Context, id, blockId string) error {
	const op = "service.AddTagToNote"

	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(blockId); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if err := s.blk.Delete(ctx, blockId); err != nil {
			return nil, err
		}
		return nil, s.nts.DeleteBlock(ctx, id, blockId)
	})

	return err
}

func (s *NotesService) ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error {
	const op = "service.AddTagToNote"

	if err := idValidation(noteID); err != nil {
		return wrapServiceCheck(op, err)
	}

	if oldOrder == newOrder {
		return nil
	}
	if oldOrder < 0 || newOrder < 0 {
		return wrapServiceCheck(op, errors.New("order < 0"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.ChangeBlockOrder(ctx, noteID, oldOrder, newOrder)
	})

	return err
}

func (s *NotesService) UpdateTitle(ctx context.Context, id string, nTitle string) error {
	const op = "service.AddTagToNote"

	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	if stringEmpty(nTitle) {
		return wrapServiceCheck(op, errors.New("title empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.UpdateTitle(ctx, id, nTitle)
	})

	return err
}
