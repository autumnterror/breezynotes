package textblock

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
)

type Driver struct{}

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	log.Println(block)
	b, err := FromUnified(block)
	if err != nil {
		return ""
	}
	return b.PlainText()
}
func (tb *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	return nil
}

func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := FromUnified(block)
	if err != nil {
		return nil, errors.New("bad block")
	}
	switch op {
	case "apply_style":
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		var req struct {
			Start int    `json:"start"`
			End   int    `json:"end"`
			Style string `json:"style"`
		}
		if err := json.Unmarshal(raw, &req); err != nil {
			return nil, err
		}
		if err := b.ApplyStyle(req.Start, req.End, req.Style); err != nil {
			return nil, err
		}
		nb, err := b.ToUnified()
		if err != nil {
			return nil, err
		}

		return nb.GetData().AsMap(), nil
	default:
		return nil, errors.New("unsupported type") //TODO dic errors
	}
}

func (tb *Driver) Create(ctx context.Context, data map[string]any) (*views.BlockDb, error) {
	b := &views.BlockDb{Data: data}
	return b, nil
}
