package domain

import (
	"errors"
	"github.com/autumnterror/breezynotes/pkg/text"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

func (tb *TextBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "textblock.ToUnified"

	u := &brzrpc.Block{
		Id:        tb.Id,
		Type:      tb.Type,
		NoteId:    tb.NoteId,
		CreatedAt: tb.CreatedAt,
		UpdatedAt: tb.UpdatedAt,
		IsUsed:    tb.IsUsed,
		Data:      nil,
	}

	dataMap := tb.Data.ToMap()
	if dataMap == nil {
		return u, nil
	}

	s, err := structpb.NewStruct(dataMap)
	if err != nil {
		return u, format.Error(op, err)
	}

	u.Data = s
	return u, nil
}

func FromUnifiedToTextBlock(b *brzrpc.Block) (*TextBlock, error) {
	const op = "textblock.FromUnifiedToTextBlock"
	if b == nil {
		return nil, errors.New("block is nil")
	}
	tb := &TextBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
		Data:      nil,
	}

	s := b.GetData()
	if s == nil {
		return tb, nil
	}

	data, err := text.NewDataFromMap(s.AsMap())
	if err != nil {
		return tb, format.Error(op, err)
	}

	tb.Data = data
	return tb, nil
}

type TextBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Data *text.Data `bson:"data" json:"data"`
}
