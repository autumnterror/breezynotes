package blocks

import (
	"context"
	"errors"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
	"github.com/autumnterror/breezynotes/views"
)

// Create universal func that get block, calls func blocks.createBlock by field _type.
// NEED TO REGISTER TYPE BEFORE USE.
// method create id
func (a *API) Create(ctx context.Context, _type, noteId string, data map[string]any) (string, error) {
	const op = "blocks.Create"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if pkgs.BlockRegistry[_type] == nil {
		return "", ErrTypeNotDefined
	}
	block, err := pkgs.BlockRegistry[_type].Create(ctx, data)
	if err != nil {
		return "", format.Error(op, err)
	}
	block.Id = uid.New()
	block.NoteId = noteId
	block.Type = _type
	block.CreatedAt = time.Now().UTC().Unix()
	block.UpdatedAt = time.Now().UTC().Unix()
	block.IsUsed = false

	if err := a.createBlock(ctx, block); err != nil {
		return "", format.Error(op, err)
	}

	return block.Id, nil
}

// Render universal func that calls func blocks.Render by field _type.
// NEED TO REGISTER TYPE BEFORE USE.
//func (a *API) Render(ctx context.Context, id, _type string) (*brzrpc.Block, error) {
//	const op = "blocks.Render"
//	ctx, done := context.WithTimeout(ctx, views.WaitTime)
//	defer done()
//
//	if pkgs.BlockRegistry[_type] == nil {
//		return nil, ErrTypeNotDefined
//	}
//	block, err := pkgs.BlockRegistry[_type].Render(ctx, id, _type)
//	if err != nil {
//		return nil, format.Error(op, err)
//	}
//
//	return "", nil
//}

// OpBlock universal func that get block, calls func op by field _type and set field data of block after op func. NEED TO REGISTER TYPE BEFORE USE
func (a *API) OpBlock(ctx context.Context, id, opName string, data map[string]any) error {
	const op = "blocks.OpBlock"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if err := a.UpdateUsed(ctx, id, true); err != nil {
		if errors.Is(err, mongo.ErrNotFound) {
			return ErrAlreadyUsed
		}
		return format.Error(op, err)
	}
	defer func() {
		if err := a.UpdateUsed(ctx, id, false); err != nil {
			log.Error(op, "cant set isUse to false", err)
		}
	}()

	block, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if pkgs.BlockRegistry[block.GetType()] == nil {
		return ErrTypeNotDefined
	}
	if err := pkgs.BlockRegistry[block.GetType()].Op(ctx, block, opName, data); err != nil {
		return format.Error(op, err)
	}

	if err := a.updateData(ctx, id, block.GetData().AsMap()); err != nil {
		return format.Error(op, err)
	}

	return nil
}

// GetAsFirst universal func that get block, calls func GetAsFirst by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) GetAsFirst(ctx context.Context, id string) (string, error) {
	const op = "blocks.OpBlock"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	block, err := a.Get(ctx, id)
	if err != nil {
		return "", format.Error(op, err)
	}
	if pkgs.BlockRegistry[block.GetType()] == nil {
		return "", ErrTypeNotDefined
	}

	return pkgs.BlockRegistry[block.GetType()].GetAsFirst(ctx, block), nil
}

// ChangeType universal func that get block, calls func ChangeType by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) ChangeType(ctx context.Context, id, newType string) error {
	const op = "blocks.ChangeType"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if err := a.UpdateUsed(ctx, id, true); err != nil {
		if errors.Is(err, mongo.ErrNotFound) {
			return ErrAlreadyUsed
		}
		return format.Error(op, err)
	}
	defer func() {
		if err := a.UpdateUsed(ctx, id, false); err != nil {
			log.Error(op, "cant set isUse to false", err)
		}
	}()
	block, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if pkgs.BlockRegistry[block.GetType()] == nil {
		return ErrTypeNotDefined
	}
	return pkgs.BlockRegistry[block.GetType()].ChangeType(ctx, block, newType)
}
