package textblock

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
	blockpkg "github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
)

type Driver struct{}

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domainblocks.FromUnifiedToTextBlock(block)
	if err != nil {
		return ""
	}
	if b.Data == nil {
		return ""
	}
	return b.Data.PlainText()
}

func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domainblocks.FromUnifiedToTextBlock(block)
	if err != nil {
		return nil, errors.New("bad block")
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	switch op {
	case "apply_style":
		return applyStyleOp(b, raw)
	case "insert_text":
		return insertTextOp(b, raw)
	case "delete_range":
		return deleteRangeOp(b, raw)
	default:
		return nil, errors.New("unsupported type")
	}
}

func (tb *Driver) Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error) {
	const op = "textblock.create"
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, format.Error(op, err)
	}

	b := &brzrpc.Block{Data: s}

	return b, nil
}

func (tb *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	const op = "textblock.ChangeType"

	b, err := domainblocks.FromUnifiedToTextBlock(block)
	if err != nil {
		return format.Error(op, err)
	}
	var plainText string
	if b.Data != nil {
		plainText = b.Data.PlainText()
	}
	newData, err := blockpkg.ChangeTypeUnif(b.Data, plainText, newType, 0, 0)
	if err != nil {
		return format.Error(op, err)
	}

	s, err := structpb.NewStruct(newData)
	if err != nil {
		return format.Error(op, err)
	}
	block.Data = s
	return nil
}
