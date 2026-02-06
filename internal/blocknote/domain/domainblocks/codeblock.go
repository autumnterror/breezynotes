package domainblocks

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type CodeBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Data *CodeData `bson:"data" json:"data"`
}

type CodeData struct {
	Text string `json:"text" bson:"text"`
	Lang string `json:"lang" bson:"lang"`
}

func NewCodeDataFromMap(obj map[string]any) (*CodeData, error) {
	if obj == nil {
		return nil, nil
	}

	var cd CodeData

	if text, ok := obj["text"].(string); ok {
		cd.Text = text
	}
	if lang, ok := obj["lang"].(string); ok {
		cd.Lang = lang
	}

	return &cd, nil
}

func (cd *CodeData) ToMap() map[string]any {
	if cd == nil {
		return nil
	}
	return map[string]any{
		"text": cd.Text,
		"lang": cd.Lang,
	}
}

func FromUnifiedToCodeBlock(b *brzrpc.Block) (*CodeBlock, error) {
	const op = "codeblock.FromUnifiedToCodeBlock"

	cb := &CodeBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
	}

	s := b.GetData()
	if s == nil {
		return cb, nil
	}

	codeData, err := NewCodeDataFromMap(s.AsMap())
	if err != nil {
		return cb, format.Error(op, err)
	}

	cb.Data = codeData
	return cb, nil
}

func (cb *CodeBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "codeblock.ToUnified"

	u := &brzrpc.Block{
		Id:        cb.Id,
		Type:      cb.Type,
		NoteId:    cb.NoteId,
		CreatedAt: cb.CreatedAt,
		UpdatedAt: cb.UpdatedAt,
		IsUsed:    cb.IsUsed,
		Data:      nil,
	}

	dataMap := cb.Data.ToMap()
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
