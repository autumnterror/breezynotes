package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

func (s *BlocksService) Get(ctx context.Context, id string) (*domain.Block, error) {
	const op = "service.GetBlock"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.bls.Get(ctx, id)
}

func (s *BlocksService) Delete(ctx context.Context, id string) error {
	const op = "service.DeleteBlock"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.bls.Delete(ctx, id)
	})

	return err
}

func (s *BlocksService) Create(ctx context.Context, _type, noteId string, data map[string]any) (string, error) {
	const op = "service.CreateBlock"

	if stringEmpty(_type) {
		return "", wrapServiceCheck(op, errors.New("_type is empty"))
	}
	if idValidation(noteId) != nil {
		return "", wrapServiceCheck(op, errors.New("bad id"))
	}

	res, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return s.bls.Create(ctx, _type, noteId, data)
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

func (s *BlocksService) OpBlock(ctx context.Context, id, opName string, data map[string]any) error {
	const op = "service.OpBlock"

	if stringEmpty(opName) {
		return wrapServiceCheck(op, errors.New("opName is empty"))
	}
	if idValidation(id) != nil {
		return wrapServiceCheck(op, errors.New("bad id"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.bls.OpBlock(ctx, id, opName, data)
	})

	return err
}

func (s *BlocksService) ChangeType(ctx context.Context, id, newType string) error {
	const op = "service.ChangeTypeBlock"
	if stringEmpty(newType) {
		return wrapServiceCheck(op, errors.New("newType is empty"))
	}
	if idValidation(id) != nil {
		return wrapServiceCheck(op, errors.New("bad id"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.bls.ChangeType(ctx, id, newType)
	})

	return err
}
