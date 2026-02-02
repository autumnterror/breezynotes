package codeblock

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

func (tb *Driver) GetAsFirst(ctx context.Context, block *brzrpc.Block) string {
	b, err := domain.FromUnifiedToCodeBlock(block)
	if err != nil {
		return ""
	}
	if b.Data == nil {
		return ""
	}
	return b.Data.Text
}
func (tb *Driver) Op(ctx context.Context, block *brzrpc.Block, op string, data map[string]any) (map[string]any, error) {
	b, err := domain.FromUnifiedToCodeBlock(block)
	if err != nil {
		return nil, errors.New("bad block")
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch op {
	case "change_text":
		return changeText(b, raw)
	case "analyse_lang":
		return analyseLang(b)
	default:
		return nil, domain.ErrUnsupportedType
	}
}

func (tb *Driver) Create(ctx context.Context, data map[string]any) (*brzrpc.Block, error) {
	const op = "codeblock.create"
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, format.Error(op, err)
	}

	b := &brzrpc.Block{Data: s}

	cb, err := domain.FromUnifiedToCodeBlock(b)
	if err != nil {
		return nil, err
	}

	newData, err := analyseLang(cb)
	if err != nil {
		return nil, err
	}

	s, err = structpb.NewStruct(newData)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.Block{Data: s}, nil
}
func (tb *Driver) ChangeType(ctx context.Context, block *brzrpc.Block, newType string) error {
	const op = "codeblock.ChangeType"
	b, err := domain.FromUnifiedToCodeBlock(block)
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
