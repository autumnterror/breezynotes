package domainblocks

import (
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
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

	m := s.AsMap()
	tb.Data = &TextData{}

	rawTextData, ok := m["text_data"]
	if !ok || rawTextData == nil {
		return tb, nil
	}

	textDataMap, ok := rawTextData.(map[string]any)
	if !ok {
		return tb, format.Error(op, errors.New("text_data is not an object"))
	}

	data, err := text.NewDataFromMap(textDataMap)
	if err != nil {
		return tb, format.Error(op, err)
	}

	tb.Data = &TextData{TextData: data}
	return tb, nil
}

type TextBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"note_id" json:"note_id"`

	CreatedAt int64 `bson:"created_at" json:"created_at"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`

	IsUsed bool `bson:"is_used"`

	Data *TextData `bson:"data" json:"data"`
}

type TextData struct {
	TextData *text.Data `json:"text_data" bson:"text_data"`
}

func (ld *TextData) ToMap() map[string]any {
	if ld == nil {
		return nil
	}

	textmap := map[string]any{}

	if ld.TextData != nil {
		textMap := ld.TextData.ToMap()
		if textMap != nil {
			textmap["text_data"] = textMap
		}
	}

	return textmap
}
