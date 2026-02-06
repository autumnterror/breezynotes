package headerblock

import (
	"context"
	"encoding/json"
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain/domainblocks"
	blockpkg "github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type Driver struct{}

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domainblocks.FromUnifiedToHeaderBlock(block)
	if err != nil {
		return ""
	}
	if b.Data == nil {
		return ""
	}
	if b.Data.TextData == nil {
		return ""
	}
	return b.Data.TextData.PlainText()
}
func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domainblocks.FromUnifiedToHeaderBlock(block)
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
	case "change_level":
		return changeLevel(b, raw)
	default:
		return nil, domainblocks.ErrUnsupportedType
	}
}

func (tb *Driver) Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error) {
	const op = "headerblock.create"
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, format.Error(op, err)
	}

	b := &brzrpc.Block{Data: s}

	return b, nil
}
func (tb *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	const op = "listblock.ChangeType"
	b, err := domainblocks.FromUnifiedToHeaderBlock(block)
	if err != nil {
		return format.Error(op, err)
	}

	var textData *text.Data
	var plainText string
	if b.Data != nil {
		textData = b.Data.TextData
		if b.Data.TextData != nil {
			plainText = b.Data.TextData.PlainText()
		}
	}
	newData, err := blockpkg.ChangeTypeUnif(textData, plainText, newType, 0, 0)
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
