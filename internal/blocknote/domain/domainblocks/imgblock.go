package domainblocks

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

func FromUnifiedToImgBlock(b *brzrpc.Block) (*ImgBlock, error) {
	const op = "ImgBlock.FromUnifiedToImgBlock"

	ib := &ImgBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
	}

	s := b.GetData()
	if s == nil {
		return ib, nil
	}

	imgData, err := NewImgDataFromMap(s.AsMap())
	if err != nil {
		return ib, format.Error(op, err)
	}

	ib.Data = imgData
	return ib, nil
}

func (ib *ImgBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "ImgBlock.ToUnified"

	u := &brzrpc.Block{
		Id:        ib.Id,
		Type:      ib.Type,
		NoteId:    ib.NoteId,
		CreatedAt: ib.CreatedAt,
		UpdatedAt: ib.UpdatedAt,
		IsUsed:    ib.IsUsed,
		Data:      nil,
	}

	dataMap := ib.Data.ToMap()
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

type ImgBlock struct {
	Id        string   `bson:"_id" json:"id"`
	Type      string   `bson:"type" json:"type"`
	NoteId    string   `bson:"noteId" json:"noteId"`
	CreatedAt int64    `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64    `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool     `bson:"is_used"`
	Data      *ImgData `bson:"data" json:"data"`
}

type ImgData struct {
	Src string `bson:"src" json:"src"`
	Alt string `bson:"alt" json:"alt"`
}

func NewImgDataFromMap(obj map[string]any) (*ImgData, error) {
	if obj == nil {
		return nil, nil
	}

	var id ImgData

	if src, ok := obj["src"].(string); ok {
		id.Src = src
	}
	if alt, ok := obj["alt"].(string); ok {
		id.Alt = alt
	}

	return &id, nil
}

func (id *ImgData) ToMap() map[string]any {
	if id == nil {
		return nil
	}
	return map[string]any{
		"src": id.Src,
		"alt": id.Alt,
	}
}
