package imgblock

import (
	"context"
	"encoding/json"
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	blockpkg "github.com/autumnterror/breezynotes/pkg/block"
	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/autumnterror/breezynotes/pkg/text"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type Driver struct{}

// GetAsFirst returns the alt text of the image.
func (d *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domain.FromUnifiedToImgBlock(block)
	if err != nil || b.Data == nil {
		return ""
	}
	return b.Data.Alt
}

// Op executes a specific operation on the image block.
func (d *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domain.FromUnifiedToImgBlock(block)
	if err != nil {
		return nil, errors.New("bad block: not an image block")
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch op {
	case "change_src":
		return changeSrc(b, raw)
	case "change_alt":
		return changeAlt(b, raw)
	default:
		return nil, domain.ErrUnsupportedType
	}
}

// Create creates a new ImgBlock.
func (d *Driver) Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error) {
	const op = "imgblock.create"
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, format.Error(op, err)
	}
	return &brzrpc.Block{Data: s}, nil
}

// ChangeType converts the ImgBlock to another block type.
func (d *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	const op = "imgblock.ChangeType"
	b, err := domain.FromUnifiedToImgBlock(block)
	if err != nil {
		return format.Error(op, err)
	}

	var plainText string
	if b.Data != nil {
		plainText = b.Data.Alt
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
