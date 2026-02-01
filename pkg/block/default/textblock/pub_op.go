package textblock

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
)

type Driver struct{}

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domain.FromUnifiedToTextBlock(block)
	if err != nil {
		return ""
	}
	if b.Data == nil {
		return ""
	}
	return b.Data.PlainText()
}

func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domain.FromUnifiedToTextBlock(block)
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

	b, err := domain.FromUnifiedToTextBlock(block)
	if err != nil {
		return format.Error(op, err)
	}

	var newData map[string]any
	switch newType {
	case domain.TextBlockType:
		return nil
	case domain.ListBlockToDoType:
		nd := domain.ListData{
			TextData: b.Data,
			Level:    0,
			Type:     domain.ListBlockToDoType,
			Value:    0,
		}
		newData = nd.ToMap()
	case domain.ListBlockUnorderedType:
		nd := domain.ListData{
			TextData: b.Data,
			Level:    0,
			Type:     domain.ListBlockUnorderedType,
			Value:    0,
		}
		newData = nd.ToMap()
	case domain.ListBlockOrderedType:
		nd := domain.ListData{
			TextData: b.Data,
			Level:    0,
			Type:     domain.ListBlockOrderedType,
			Value:    1,
		}
		newData = nd.ToMap()
	case domain.CodeBlockType:
		var nd domain.CodeData
		if b.Data != nil {
			nd = domain.CodeData{
				Text: b.Data.PlainText(),
				Lang: "undefined",
			}
		} else {
			nd = domain.CodeData{
				Text: "",
				Lang: "undefined",
			}
		}
		newData = nd.ToMap()
	case domain.HeaderBlockType1:
		nd := domain.HeaderData{
			TextData: b.Data,
			Level:    1,
		}
		newData = nd.ToMap()
	case domain.HeaderBlockType2:
		nd := domain.HeaderData{
			TextData: b.Data,
			Level:    2,
		}
		newData = nd.ToMap()
	case domain.HeaderBlockType3:
		nd := domain.HeaderData{
			TextData: b.Data,
			Level:    3,
		}
		newData = nd.ToMap()
	case domain.FileBlockType:
		nd := domain.FileData{
			Src: "",
		}
		newData = nd.ToMap()
	case domain.LinkBlockType:
		var nd domain.LinkData
		if b.Data != nil {
			nd = domain.LinkData{
				Text: b.Data.PlainText(),
			}
		} else {
			nd = domain.LinkData{
				Text: "",
			}
		}
		newData = nd.ToMap()
	case domain.ImgBlockType:
		var nd domain.ImgData
		if b.Data != nil {
			nd = domain.ImgData{
				Src: "",
				Alt: b.Data.PlainText(),
			}
		} else {
			nd = domain.ImgData{
				Src: "",
				Alt: "",
			}
		}
		newData = nd.ToMap()
	case domain.QuoteBlockType:
		var nd domain.QuoteData
		if b.Data != nil {
			nd = domain.QuoteData{
				Text: b.Data.PlainText(),
			}
		} else {
			nd = domain.QuoteData{
				Text: "",
			}
		}
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
