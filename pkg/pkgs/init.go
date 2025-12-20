package pkgs

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/views"
)

type BlockRepo interface {
	Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error)
	GetAsFirst(ctx context.Context, block *brzrpc.Block) string
	ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error
	Create(ctx context.Context, data map[string]any) (*views.BlockDb, error)
	//Render(ctx context.Context, block *brzrpc.Block) (*brzrpc.Block, error)
}

var BlockRegistry = make(map[string]BlockRepo)

// RegisterBlock register new type of block.
func RegisterBlock(blockType string, b BlockRepo) {
	BlockRegistry[blockType] = b
}
