package textblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
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
