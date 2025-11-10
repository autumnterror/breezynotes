package pkgs

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
)

type BlockRepo interface {
	Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) error
	GetAsFirst(ctx context.Context, block *brzrpc.Block) string
	ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error
}

var BlockRegistry = make(map[string]BlockRepo)

// RegisterBlock register new type of block.
func RegisterBlock(blockType string, b BlockRepo) {
	BlockRegistry[blockType] = b
}
