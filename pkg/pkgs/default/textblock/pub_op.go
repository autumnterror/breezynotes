package textblock

import (
	"context"
	"encoding/json"
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
)

type Driver struct{}

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := FromUnified(block)
	if err != nil {
		return ""
	}
	return b.PlainText()
}
func (tb *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	return nil
}

func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) error {
	b, err := FromUnified(block)
	if err != nil {
		return errors.New("bad block")
	}
	switch op {
	case "apply_style":
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		var req struct {
			start int
			end   int
			style string
		}
		if err := json.Unmarshal(raw, &req); err != nil {
			return err
		}
		if err := b.ApplyStyle(req.start, req.end, req.style); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unsupported type") //TODO dic errors
	}
}

func (tb *Driver) Create(ctx context.Context, data map[string]any) (*views.BlockDb, error) {
	b := &views.BlockDb{Data: data}
	return b, nil
}
