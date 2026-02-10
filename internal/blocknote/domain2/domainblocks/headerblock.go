package domainblocks

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

func FromUnifiedToHeaderBlock(b *brzrpc.Block) (*HeaderBlock, error) {
	const op = "listblock.FromUnifiedToHeaderBlock"

	lb := HeaderBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
	}

	s := b.GetData()
	if s == nil {
		return &lb, nil
	}

	headerData, err := NewHeaderDataFromMap(s.AsMap())
	if err != nil {
		return &lb, format.Error(op, err)
	}

	lb.Data = headerData
	return &lb, nil
}

func (lb *HeaderBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "listblock.ToUnified"

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

type HeaderBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Data *HeaderData `bson:"data" json:"data"`
}

type HeaderData struct {
	TextData *text.Data `json:"text_data" bson:"text_data"`
	Level    uint       `json:"level" bson:"level"`
}

func NewHeaderDataFromMap(obj map[string]any) (*HeaderData, error) {
	const op = "listblock.NewHeaderDataFromMap"
	if obj == nil {
		return nil, nil
	}

	var hd HeaderData

	if level, ok := obj["level"].(float64); ok {
		hd.Level = uint(level)
	}

	if rawTextData, ok := obj["text_data"].(map[string]any); ok {
		data, err := text.NewDataFromMap(rawTextData)
		if err != nil {
			return nil, format.Error(op, err)
		}
		hd.TextData = data
	}

	return &hd, nil
}

func (hd *HeaderData) ToMap() map[string]any {
	if hd == nil {
		return nil
	}

	headerMap := map[string]any{
		"level": hd.Level,
	}

	if hd.TextData != nil {
		textMap := hd.TextData.ToMap()
		if textMap != nil {
			headerMap["text_data"] = textMap
		}
	}

	return headerMap
}
