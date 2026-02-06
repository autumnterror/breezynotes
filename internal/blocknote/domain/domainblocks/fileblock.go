package domainblocks

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/structpb"
)

func FromUnifiedToFileBlock(b *brzrpc.Block) (*FileBlock, error) {
	const op = "FileBlock.FromUnifiedToFileBlock"

	fb := &FileBlock{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
	}

	s := b.GetData()
	if s == nil {
		return fb, nil
	}

	fileData, err := NewFileDataFromMap(s.AsMap())
	if err != nil {
		return fb, format.Error(op, err)
	}

	fb.Data = fileData
	return fb, nil
}

func (fb *FileBlock) ToUnified() (*brzrpc.Block, error) {
	const op = "FileBlock.ToUnified"

	u := &brzrpc.Block{
		Id:        fb.Id,
		Type:      fb.Type,
		NoteId:    fb.NoteId,
		CreatedAt: fb.CreatedAt,
		UpdatedAt: fb.UpdatedAt,
		IsUsed:    fb.IsUsed,
		Data:      nil,
	}

	dataMap := fb.Data.ToMap()
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

type FileBlock struct {
	Id        string    `bson:"_id" json:"id"`
	Type      string    `bson:"type" json:"type"`
	NoteId    string    `bson:"noteId" json:"noteId"`
	CreatedAt int64     `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64     `bson:"updatedAt" json:"updatedAt"`
	IsUsed    bool      `bson:"is_used"`
	Data      *FileData `bson:"data" json:"data"`
}

type FileData struct {
	Src string `bson:"src" json:"src"`
}

func NewFileDataFromMap(obj map[string]any) (*FileData, error) {
	if obj == nil {
		return nil, nil
	}

	var fd FileData

	if src, ok := obj["src"].(string); ok {
		fd.Src = src
	}

	return &fd, nil
}

func (fd *FileData) ToMap() map[string]any {
	if fd == nil {
		return nil
	}
	return map[string]any{
		"src": fd.Src,
	}
}
