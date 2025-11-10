package textblock

import (
	"encoding/json"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
)

func FromUnified(u *brzrpc.Block) (TextBlock, error) {
	const op = "textblock.FromUnified"
	tb := TextBlock{
		Id:        u.GetId(),
		Type:      u.GetType(),
		NoteId:    u.GetNoteId(),
		Order:     int(u.GetOrder()),
		CreatedAt: u.GetCreatedAt(),
		UpdatedAt: u.GetUpdatedAt(),
		IsUsed:    u.GetIsUsed(),
		Text:      nil,
	}

	raw, err := json.Marshal(u.Data)
	if err != nil {
		return tb, format.Error(op, err)
	}
	if err := json.Unmarshal(raw, &tb.Text); err != nil {
		return tb, format.Error(op, err)
	}
	return tb, nil
}

func (tb *TextBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "textblock.ToUnified"
	u := brzrpc.Block{
		Id:        tb.Id,
		Type:      tb.Type,
		NoteId:    tb.NoteId,
		Order:     int32(tb.Order),
		CreatedAt: tb.CreatedAt,
		UpdatedAt: tb.UpdatedAt,
		IsUsed:    tb.IsUsed,
		Data:      nil,
	}
	raw, err := json.Marshal(tb.Text)
	if err != nil {
		return &u, format.Error(op, err)
	}
	if err := json.Unmarshal(raw, &u.Data); err != nil {
		return &u, format.Error(op, err)
	}
	return &u, nil
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
	Style string
	Text  string
}
