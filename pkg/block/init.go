package block

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
)

type BlockRepo interface {
	Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error)
	GetAsFirst(ctx context.Context, block *brzrpc.Block) string
	ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error
	Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error)
	//Render(ctx context.Context, block *domain.Block) (*domain.Block, error)
}

var BlockRegistry = make(map[string]BlockRepo)

// RegisterBlock register new type of block.
func RegisterBlock(blockType string, b BlockRepo) {
	BlockRegistry[blockType] = b
}

func GetRegisteredTypes() []string {
	var res []string
	for t := range BlockRegistry {
		res = append(res, t)
	}
	return res
}
