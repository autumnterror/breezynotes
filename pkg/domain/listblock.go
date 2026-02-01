package domain

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/text"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

func FromUnifiedToListBlock(b *brzrpc.Block) (*ListBlock, error) {
	const op = "listblock.FromUnifiedToListBlock"

	lb := ListBlock{
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

	listData, err := NewListDataFromMap(s.AsMap())
	if err != nil {
		return &lb, format.Error(op, err)
	}

	lb.Data = listData
	return &lb, nil
}

func (lb *ListBlock) ToUnified() (*brzrpc.Block, error) {
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

type ListBlock struct {
	Id     string `bson:"_id" json:"id"`
	Type   string `bson:"type" json:"type"`
	NoteId string `bson:"noteId" json:"noteId"`

	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool  `bson:"is_used"`

	Data *ListData `bson:"data" json:"data"`
}

type ListData struct {
	TextData *text.Data `json:"text_data" bson:"text_data"`
	Level    uint       `json:"level" bson:"level"`
	Type     string     `json:"type" bson:"type"`
	Value    int        `json:"value" bson:"value"`
}

func NewListDataFromMap(obj map[string]any) (*ListData, error) {
	const op = "listblock.NewListDataFromMap"
	if obj == nil {
		return nil, nil
	}

	var ld ListData

	if level, ok := obj["level"].(float64); ok {
		ld.Level = uint(level)
	}
	if typeStr, ok := obj["type"].(string); ok {
		ld.Type = typeStr
	}
	if value, ok := obj["value"].(float64); ok {
		ld.Value = int(value)
	}

	if rawTextData, ok := obj["text_data"].(map[string]any); ok {
		data, err := text.NewDataFromMap(rawTextData)
		if err != nil {
			return nil, format.Error(op, err)
		}
		ld.TextData = data
	}

	return &ld, nil
}

func (ld *ListData) ToMap() map[string]any {
	if ld == nil {
		return nil
	}

	listMap := map[string]any{
		"level": ld.Level,
		"type":  ld.Type,
		"value": ld.Value,
	}

	if ld.TextData != nil {
		textMap := ld.TextData.ToMap()
		if textMap != nil {
			listMap["text_data"] = textMap
		}
	}

	return listMap
}
