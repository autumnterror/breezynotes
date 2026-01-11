package textblock

import (
	"fmt"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

func FromUnified(b *brzrpc.Block) (TextBlock, error) {
	const op = "textblock.FromUnified"

	tb := TextBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
		Text:      nil,
	}

	s := b.GetData()
	if s == nil {
		return tb, nil
	}

	m := s.AsMap()

	rawText, ok := m["text"]
	if !ok {
		return tb, nil
	}

	list, ok := rawText.([]any)
	if !ok {
		return tb, format.Error(op, fmt.Errorf(`field "text" has unexpected type %T, want []interface{}`, rawText))
	}

	texts := make([]TextData, 0, len(list))
	for i, v := range list {
		obj, ok := v.(map[string]any)
		if !ok {
			return tb, format.Error(op, fmt.Errorf("text[%d] has unexpected type %T, want map[string]interface{}", i, v))
		}

		style, _ := obj["style"].(string)
		text, _ := obj["text"].(string)

		texts = append(texts, TextData{
			Style: style,
			Text:  text,
		})
	}

	tb.Text = texts
	return tb, nil
}

func (tb *TextBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "textblock.ToUnified"

	u := &brzrpc.Block{
		Id:        tb.Id,
		Type:      tb.Type,
		NoteId:    tb.NoteId,
		Order:     int32(tb.Order),
		CreatedAt: tb.CreatedAt,
		UpdatedAt: tb.UpdatedAt,
		IsUsed:    tb.IsUsed,
		Data:      nil,
	}
	arr := make([]interface{}, 0, len(tb.Text))
	for _, t := range tb.Text {
		arr = append(arr, map[string]interface{}{
			"style": t.Style,
			"text":  t.Text,
		})
	}

	data := map[string]interface{}{
		"text": arr,
	}

	s, err := structpb.NewStruct(data)
	if err != nil {
		return u, format.Error(op, err)
	}

	u.Data = s
	return u, nil
}

type TextBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`
	Order  int    `bson:"order" json:"order"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Text []TextData `bson:"data" json:"text"`
}

type TextData struct {
	Style string `json:"style" bson:"style"`
	Text  string `json:"text" bson:"text"`
}
