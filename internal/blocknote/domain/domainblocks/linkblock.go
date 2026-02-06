package domainblocks

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type LinkBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Data *LinkData `bson:"data" json:"data"`
}

type LinkData struct {
	Text string `json:"text" bson:"text"`
	Url  string `json:"url" bson:"url"`
}

func NewLinkDataFromMap(obj map[string]any) (*LinkData, error) {
	if obj == nil {
		return nil, nil
	}

	var ld LinkData

	if text, ok := obj["text"].(string); ok {
		ld.Text = text
	}
	if url, ok := obj["url"].(string); ok {
		ld.Url = url
	}
	return &ld, nil
}

func (ld *LinkData) ToMap() map[string]any {
	if ld == nil {
		return nil
	}
	return map[string]any{
		"text": ld.Text,
		"url":  ld.Url,
	}
}

func FromUnifiedToLinkBlock(b *brzrpc.Block) (*LinkBlock, error) {
	const op = "linkblock.FromUnifiedToLinkBlock"

	lb := &LinkBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
	}

	s := b.GetData()
	if s == nil {
		return lb, nil
	}

	linkData, err := NewLinkDataFromMap(s.AsMap())
	if err != nil {
		return lb, format.Error(op, err)
	}

	lb.Data = linkData
	return lb, nil
}

func (lb *LinkBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "linkblock.ToUnified"

	u := &brzrpc.Block{
		Id:        lb.Id,
		Type:      lb.Type,
		NoteId:    lb.NoteId,
		CreatedAt: lb.CreatedAt,
		UpdatedAt: lb.UpdatedAt,
		IsUsed:    lb.IsUsed,
		Data:      nil,
	}

	dataMap := lb.Data.ToMap()
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
