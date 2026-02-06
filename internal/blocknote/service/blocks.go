package service

import (
	"context"
	"errors"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain/domainblocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/autumnterror/utils_go/pkg/utils/uid"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
)

func (s *BN) GetRegisteredBlocks(ctx context.Context) []string {
	return block.GetRegisteredTypes()
}

func (s *BN) ChangeBlockOrder(ctx context.Context, idNote, idUser string, oldOrder, newOrder int) error {
	const op = "service.ChangeBlockOrder"

	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	if oldOrder == newOrder {
		return nil
	}
	if oldOrder < 0 || newOrder < 0 {
		return wrapServiceCheck(op, errors.New("order < 0"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote)
		if err != nil {
			return nil, domain.ErrNotFound
		}

		if n.Author != idUser && !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.ChangeBlockOrder(ctx, idNote, oldOrder, newOrder)
	})

	return err
}

func (s *BN) GetBlock(ctx context.Context, idBlock, idNote, idUser string) (*domain.Block, error) {
	const op = "service.GetBlock"
	if err := idValidation(idBlock); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := idValidation(idNote); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	n, err := s.nts.Get(ctx, idNote)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
		return nil, domain.ErrUnauthorized
	}

	return s.blk.Get(ctx, idBlock)
}

func (s *BN) GetBlocks(ctx context.Context, ids []string) (*domain.Blocks, error) {
	const op = "service.GetBlocks"

	for _, id := range ids {
		if err := idValidation(id); err != nil {
			return nil, wrapServiceCheck(op, err)
		}
	}

	return s.blk.GetMany(ctx, ids)
}

func (s *BN) DeleteBlock(ctx context.Context, noteId, blockId, idUser string) error {
	const op = "service.AddTagToNote"

	if idValidation(noteId) != nil {
		return wrapServiceCheck(op, errors.New("bad note noteId"))
	}
	if idValidation(blockId) != nil {
		return wrapServiceCheck(op, errors.New("bad block noteId"))
	}
	if idValidation(idUser) != nil {
		return wrapServiceCheck(op, errors.New("bad block noteId"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, noteId); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		if err := s.blk.Delete(ctx, blockId); err != nil {
			return nil, err
		}
		return nil, s.nts.DeleteBlock(ctx, noteId, blockId)
	})

	return err
}

func (s *BN) CreateBlock(ctx context.Context, _type, noteId string, data map[string]any, pos int, idUser string) (string, error) {
	const op = "service.CreateBlock"

	if pos < 0 {
		return "", wrapServiceCheck(op, errors.New("pos < 0"))
	}
	if idValidation(noteId) != nil {
		return "", wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(idUser) != nil {
		return "", wrapServiceCheck(op, errors.New("bad user id"))
	}

	if stringEmpty(_type) {
		return "", wrapServiceCheck(op, errors.New("type is empty"))
	}

	res, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, noteId); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		if block.Registry[_type] == nil {
			return "", domain.ErrTypeNotDefined
		}
		b, err := block.Registry[_type].Create(ctx, data)
		if err != nil {
			return "", format.Error(op, err)
		}

		b.Id = uid.New()
		b.NoteId = noteId
		b.Type = _type
		b.CreatedAt = time.Now().UTC().Unix()
		b.UpdatedAt = time.Now().UTC().Unix()
		b.IsUsed = false

		if err := s.blk.CreateBlock(ctx, domain.ToBlockDb(b)); err != nil {
			return "", format.Error(op, err)
		}

		return b.Id, s.nts.InsertBlock(ctx, noteId, b.Id, pos)
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

func (s *BN) OpBlock(ctx context.Context, id, opName string, data map[string]any, idNote, idUser string) error {
	const op = "service.OpBlock"

	if stringEmpty(opName) {
		return wrapServiceCheck(op, errors.New("op is empty"))
	}
	if idValidation(id) != nil {
		return wrapServiceCheck(op, errors.New("bad id"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, idNote); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		b, err := s.blk.Get(ctx, id)
		if err != nil {
			return nil, domain.ErrNotFound
		}

		if err := s.blk.UpdateUsed(ctx, id, true); err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.ErrAlreadyUsed
			}
			return nil, err
		}

		if block.Registry[b.Type] == nil {
			return nil, domain.ErrTypeNotDefined
		}

		newData, err := block.Registry[b.Type].Op(ctx, domain.FromBlockDb(b), opName, data)
		if err != nil {
			return nil, err
		}

		if newData == nil {
			return nil, nil
		}

		if err := s.blk.UpdateData(ctx, id, newData); err != nil {
			return nil, err
		}
		if err := s.blk.UpdateUsed(ctx, id, false); err != nil {
			return nil, err
		}
		if err := s.nts.UpdateUpdatedAt(ctx, idNote); err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func (s *BN) ChangeTypeBlock(ctx context.Context, idBlock, idNote, idUser, newType string) error {
	const op = "service.ChangeTypeBlock"
	if stringEmpty(newType) {
		return wrapServiceCheck(op, errors.New("type is empty"))
	}
	if idValidation(idBlock) != nil {
		return wrapServiceCheck(op, errors.New("bad idBlock"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, idNote); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		b, err := s.blk.Get(ctx, idBlock)
		if err != nil {
			return nil, format.Error(op, err)
		}

		if block.Registry[b.Type] == nil {
			return nil, domain.ErrTypeNotDefined
		}

		nb := domain.FromBlockDb(b)
		err = block.Registry[b.Type].ChangeType(ctx, nb, newType)
		if err != nil {
			return nil, format.Error(op, err)
		}

		switch newType {
		case domainblocks.ListBlockToDoType, domainblocks.ListBlockOrderedType, domainblocks.ListBlockUnorderedType:
			newType = "list"
		case domainblocks.HeaderBlockType1, domainblocks.HeaderBlockType2, domainblocks.HeaderBlockType3:
			newType = "header"
		}

		if err := s.blk.UpdateType(ctx, idBlock, newType); err != nil {
			return nil, format.Error(op, err)
		}
		if err := s.blk.UpdateData(ctx, idBlock, nb.Data.AsMap()); err != nil {
			return nil, format.Error(op, err)
		}
		return nil, nil
	})

	return err
}
