package textblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
)

type Driver struct{}

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	return ""
}
func (tb *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	return nil
}

func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) error {
	return nil
}

func (tb *Driver) Create(ctx context.Context, _type string, data map[string]any) (*views.BlockDb, error) {
	return nil, nil
}

func (tb *Driver) Render(ctx context.Context) (*brzrpc.Block, error) {
	return nil, nil
}
