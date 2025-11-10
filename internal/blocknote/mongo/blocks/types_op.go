package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
)

// TODO CreateBlock

// OpBlock universal func that get block, calls func op by field _type and set field data of block after op func. NEED TO REGISTER TYPE BEFORE USE
func (a *API) OpBlock(ctx context.Context, id, opName string, data map[string]any) error {
	const op = "blocks.OpBlock"
	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

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

	if err := a.UpdateData(ctx, id, block.GetData().AsMap()); err != nil {
		return format.Error(op, err)
	}

	return nil
}

// GetAsFirst universal func that get block, calls func GetAsFirst by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) GetAsFirst(ctx context.Context, id string) (string, error) {
	const op = "blocks.OpBlock"
	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
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
	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	block, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if pkgs.BlockRegistry[block.GetType()] == nil {
		return ErrTypeNotDefined
	}
	return pkgs.BlockRegistry[block.GetType()].ChangeType(ctx, block, newType)
}
