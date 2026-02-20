package domainblocks

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type QuoteBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Data *QuoteData `bson:"data" json:"data"`
}

type QuoteData struct {
	Text string `json:"text" bson:"text"`
}

func NewQuoteDataFromMap(obj map[string]any) (*QuoteData, error) {
	if obj == nil {
		return nil, nil
	}

	var qd QuoteData

	if text, ok := obj["text"].(string); ok {
		qd.Text = text
	}

	return &qd, nil
}

func (qd *QuoteData) ToMap() map[string]any {
	if qd == nil {
		return nil
	}
	return map[string]any{
		"text": qd.Text,
	}
}

func FromUnifiedToQuoteBlock(b *brzrpc.Block) (*QuoteBlock, error) {
	const op = "quoteblock.FromUnifiedToQuoteBlock"

	qb := &QuoteBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
	}

	s := b.GetData()
	if s == nil {
		return qb, nil
	}

	quoteData, err := NewQuoteDataFromMap(s.AsMap())
	if err != nil {
		return qb, format.Error(op, err)
	}

	qb.Data = quoteData
	return qb, nil
}

func (qb *QuoteBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "quoteblock.ToUnified"

	u := &brzrpc.Block{
		Id:        qb.Id,
		Type:      qb.Type,
		NoteId:    qb.NoteId,
		CreatedAt: qb.CreatedAt,
		UpdatedAt: qb.UpdatedAt,
		IsUsed:    qb.IsUsed,
		Data:      nil,
	}

	dataMap := qb.Data.ToMap()
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
