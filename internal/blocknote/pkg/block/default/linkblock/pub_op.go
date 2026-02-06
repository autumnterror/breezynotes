package linkblock

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

// GetAsFirst returns the text of the link.
func (d *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domainblocks.FromUnifiedToLinkBlock(block)
	if err != nil || b.Data == nil {
		return ""
	}
	return b.Data.Text
}

// Op executes a specific operation on the link block.
func (d *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domainblocks.FromUnifiedToLinkBlock(block)
	if err != nil {
		return nil, errors.New("bad block: not a link block")
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch op {
	case "change_text":
		return changeText(b, raw)
	case "change_url":
		return changeUrl(b, raw)
	default:
		return nil, domainblocks.ErrUnsupportedType
	}
}

// Create creates a new LinkBlock.
func (d *Driver) Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error) {
	const op = "linkblock.create"
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, format.Error(op, err)
	}
	return &brzrpc.Block{Data: s}, nil
}

// ChangeType converts the LinkBlock to another block type.
func (d *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	const op = "linkblock.ChangeType"
	b, err := domainblocks.FromUnifiedToLinkBlock(block)
	if err != nil {
		return format.Error(op, err)
	}

	var plainText string
	if b.Data != nil {
		plainText = b.Data.Text
	}
	ntxtd := &text.Data{Text: []text.Part{{
		Style:  "default",
		String: plainText,
	}}}

	newData, err := blockpkg.ChangeTypeUnif(ntxtd, plainText, newType, 0, 0)
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
