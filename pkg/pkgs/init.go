package pkgs

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

type BlockRepo interface {
	Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error)
	GetAsFirst(ctx context.Context, block *brzrpc.Block) string
	ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error
	Create(ctx context.Context, data map[string]any) (*domain.Block, error)
	//Render(ctx context.Context, block *domain.Block) (*domain.Block, error)
}

var BlockRegistry = make(map[string]BlockRepo)

// RegisterBlock register new type of block.
func RegisterBlock(blockType string, b BlockRepo) {
	BlockRegistry[blockType] = b
}
