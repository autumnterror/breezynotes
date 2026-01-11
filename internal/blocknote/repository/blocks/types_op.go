package blocks

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"time"

	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
)

// Create universal func that get block, calls func blocks.createBlock by field _type.
// NEED TO REGISTER TYPE BEFORE USE.
// method create id
func (a *API) Create(ctx context.Context, _type, noteId string, data map[string]any) (string, error) {
	const op = "blocks.Create"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	if pkgs.BlockRegistry[_type] == nil {
		return "", domain.ErrTypeNotDefined
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

// OpBlock universal func that get block, calls func op by field _type and set field data of block after op func. NEED TO REGISTER TYPE BEFORE USE
func (a *API) OpBlock(ctx context.Context, id, opName string, data map[string]any) error {
	const op = "blocks.OpBlock"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	if err := a.updateUsed(ctx, id, true); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrAlreadyUsed
		}
		return format.Error(op, err)
	}
	defer func() {
		if err := a.updateUsed(ctx, id, false); err != nil {
			log.Error(op, "cant set isUse to false", err)
		}
	}()

	block, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if pkgs.BlockRegistry[block.Type] == nil {
		return domain.ErrTypeNotDefined
	}
	newData, err := pkgs.BlockRegistry[block.Type].Op(ctx, block, opName, data)
	if err != nil {
		return format.Error(op, err)
	}

	if err := a.updateData(ctx, id, newData); err != nil {
		return format.Error(op, err)
	}

	return nil
}

// GetAsFirst universal func that get block, calls func GetAsFirst by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) GetAsFirst(ctx context.Context, id string) (string, error) {
	const op = "blocks.GetAsFirst"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	block, err := a.Get(ctx, id)
	if err != nil {
		return "", format.Error(op, err)
	}
	if pkgs.BlockRegistry[block.Type] == nil {
		return "", domain.ErrTypeNotDefined
	}

	return pkgs.BlockRegistry[block.Type].GetAsFirst(ctx, block), nil
}

// ChangeType universal func that get block, calls func ChangeType by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) ChangeType(ctx context.Context, id, newType string) error {
	const op = "blocks.ChangeType"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	if err := a.updateUsed(ctx, id, true); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrAlreadyUsed
		}
		return format.Error(op, err)
	}
	defer func() {
		if err := a.updateUsed(ctx, id, false); err != nil {
			log.Error(op, "cant set isUse to false", err)
		}
	}()
	block, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if pkgs.BlockRegistry[block.Type] == nil {
		return domain.ErrTypeNotDefined
	}
	return pkgs.BlockRegistry[block.Type].ChangeType(ctx, block, newType)
}
