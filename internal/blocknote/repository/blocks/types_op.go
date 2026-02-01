package blocks

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/pkg/block"

	"time"

	"github.com/autumnterror/utils_go/pkg/log"

	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
)

// Create universal func that get block, calls func blocks.createBlock by field _type.
// NEED TO REGISTER TYPE BEFORE USE.
// method create id
func (a *API) Create(ctx context.Context, _type, noteId string, data map[string]any) (string, error) {
	const op = "blocks.Create"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	if block.BlockRegistry[_type] == nil {
		return "", domain.ErrTypeNotDefined
	}
	b, err := block.BlockRegistry[_type].Create(ctx, data)
	if err != nil {
		return "", format.Error(op, err)
	}

	b.Id = uid.New()
	b.NoteId = noteId
	b.Type = _type
	b.CreatedAt = time.Now().UTC().Unix()
	b.UpdatedAt = time.Now().UTC().Unix()
	b.IsUsed = false

	if err := a.createBlock(ctx, domain.ToBlockDb(b)); err != nil {
		return "", format.Error(op, err)
	}

	return b.Id, nil
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

	b, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if block.BlockRegistry[b.Type] == nil {
		return domain.ErrTypeNotDefined
	}
	newData, err := block.BlockRegistry[b.Type].Op(ctx, domain.FromBlockDb(b), opName, data)
	if err != nil {
		return format.Error(op, err)
	}

	if newData == nil {
		return nil
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

	b, err := a.Get(ctx, id)
	if err != nil {
		return "", format.Error(op, err)
	}
	if block.BlockRegistry[b.Type] == nil {
		return "", domain.ErrTypeNotDefined
	}

	return block.BlockRegistry[b.Type].GetAsFirst(ctx, domain.FromBlockDb(b)), nil
}

// ChangeType universal func that get block, calls func ChangeType by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) ChangeType(ctx context.Context, id, newType string) error {
	const op = "blocks.ChangeType"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	b, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if block.BlockRegistry[b.Type] == nil {
		return domain.ErrTypeNotDefined
	}
	err = block.BlockRegistry[b.Type].ChangeType(ctx, domain.FromBlockDb(b), newType)
	if err != nil {
		return format.Error(op, err)
	}
	if err := a.updateType(ctx, id, newType); err != nil {
		return format.Error(op, err)
	}
	if err := a.updateData(ctx, id, b.Data); err != nil {
		return format.Error(op, err)
	}
	return nil
}
