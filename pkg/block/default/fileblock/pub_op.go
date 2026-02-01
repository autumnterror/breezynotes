package fileblock

import (
	"context"
	"encoding/json"
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/autumnterror/breezynotes/pkg/text"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type Driver struct{}

// GetAsFirst returns the source URL of the file.
func (d *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domain.FromUnifiedToFileBlock(block)
	if err != nil || b.Data == nil {
		return ""
	}
	return b.Data.Src
}

// Op executes a specific operation on the block, like changing the source.
func (d *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domain.FromUnifiedToFileBlock(block)
	if err != nil {
		return nil, errors.New("bad block: not a file block")
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch op {
	case "change_src":
		return changeSrc(b, raw)
	default:
		return nil, domain.ErrUnsupportedType
	}
}

// Create creates a new FileBlock.
func (d *Driver) Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error) {
	const op = "fileblock.create"
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, format.Error(op, err)
	}
	return &brzrpc.Block{Data: s}, nil
}

// ChangeType converts the FileBlock to another block type.
func (d *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	const op = "fileblock.ChangeType"
	b, err := domain.FromUnifiedToFileBlock(block)
	if err != nil {
		return format.Error(op, err)
	}

	var fileSrc string
	if b.Data != nil {
		fileSrc = b.Data.Src
	}

	ntxtd := text.Data{Text: []text.Part{{
		Style:  "default",
		String: fileSrc,
	}}}

	var newData map[string]any
	switch newType {
	case domain.TextBlockType:
		newData = ntxtd.ToMap()
	case domain.ListBlockToDoType:
		nd := domain.ListData{TextData: &ntxtd, Type: domain.ListBlockToDoType}
		newData = nd.ToMap()
	case domain.ListBlockUnorderedType:
		nd := domain.ListData{TextData: &ntxtd, Type: domain.ListBlockUnorderedType}
		newData = nd.ToMap()
	case domain.ListBlockOrderedType:
		nd := domain.ListData{TextData: &ntxtd, Type: domain.ListBlockOrderedType, Value: 1}
		newData = nd.ToMap()
	case domain.CodeBlockType:
		nd := domain.CodeData{Text: fileSrc, Lang: "undefined"}
		newData = nd.ToMap()
	case domain.HeaderBlockType1, domain.HeaderBlockType2, domain.HeaderBlockType3:
		level := uint(1)
		if newType == domain.HeaderBlockType2 {
			level = 2
		}
		if newType == domain.HeaderBlockType3 {
			level = 3
		}
		nd := domain.HeaderData{TextData: &ntxtd, Level: level}
		newData = nd.ToMap()
	case domain.FileBlockType:
		nd := domain.FileData{Src: fileSrc}
		newData = nd.ToMap()
	case domain.LinkBlockType:
		nd := domain.LinkData{Text: fileSrc}
		newData = nd.ToMap()
	case domain.ImgBlockType:
		nd := domain.ImgData{Src: fileSrc, Alt: ""}
		newData = nd.ToMap()
	case domain.QuoteBlockType:
		nd := domain.QuoteData{Text: fileSrc}
		newData = nd.ToMap()
	default:
		return domain.ErrUnsupportedType
	}

	s, err := structpb.NewStruct(newData)
	if err != nil {
		return format.Error(op, err)
	}
	block.Data = s
	return nil
}
