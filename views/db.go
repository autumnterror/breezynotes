package views

import (
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

const (
	WaitTime = 3 * time.Second
)

type NoteDb struct {
	Id        string   `bson:"_id"`
	Title     string   `bson:"title"`
	CreatedAt int64    `bson:"created_at"`
	UpdatedAt int64    `bson:"updated_at"`
	Tag       *TagDb   `bson:"tag"`
	Author    string   `bson:"author"`
	Editors   []string `bson:"editors"`
	Readers   []string `bson:"readers"`
	Blocks    []string `bson:"blocks"`
	//TODO deprecated (status in table)
	Status brzrpc.Statuses `bson:"status"`
}

type TagDb struct {
	Id     string `bson:"_id"`
	Title  string `bson:"title"`
	Color  string `bson:"color"`
	Emoji  string `bson:"emoji"`
	UserId string `bson:"user_id"`
}

func ToNoteDb(n *brzrpc.Note) *NoteDb {
	if n == nil {
		return nil
	}
	return &NoteDb{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       ToTagDb(n.Tag),
		Author:    n.Author,
		Editors:   n.Editors,
		Readers:   n.Readers,
		Blocks:    n.Blocks,
		Status:    n.Status,
	}
}

func FromNoteDb(n *NoteDb) *brzrpc.Note {
	if n == nil {
		return nil
	}
	return &brzrpc.Note{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       FromTagDb(n.Tag),
		Author:    n.Author,
		Editors:   n.Editors,
		Readers:   n.Readers,
		Blocks:    n.Blocks,
		Status:    n.Status,
	}
}

func ToTagDb(t *brzrpc.Tag) *TagDb {
	if t == nil {
		return nil
	}
	return &TagDb{
		Id:     t.Id,
		Title:  t.Title,
		Color:  t.Color,
		Emoji:  t.Emoji,
		UserId: t.UserId,
	}
}

func FromTagDb(t *TagDb) *brzrpc.Tag {
	if t == nil {
		return nil
	}
	return &brzrpc.Tag{
		Id:     t.Id,
		Title:  t.Title,
		Color:  t.Color,
		Emoji:  t.Emoji,
		UserId: t.UserId,
	}
}

type BlockDb struct {
	Id     string `bson:"_id"`
	Type   string `bson:"type"`
	NoteId string `bson:"note_id"`

	// Order is deprecated
	Order int `bson:"order"`

	CreatedAt int64          `bson:"created_at"`
	UpdatedAt int64          `bson:"updated_at"`
	IsUsed    bool           `bson:"is_used"`
	Data      map[string]any `bson:"data"`
}

func ToBlockDb(b *brzrpc.Block) *BlockDb {
	return &BlockDb{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
		Order:     int(b.GetOrder()),
		CreatedAt: b.GetCreatedAt(),
		UpdatedAt: b.GetUpdatedAt(),
		IsUsed:    b.GetIsUsed(),
		Data:      b.GetData().AsMap(),
	}
}

// FromBlockDb on data field u can insert only base type. If u want ur struct convert to map[string]any
//
//	type B struct {
//		I int
//		F float32
//	}
//
//	func (b B) ToMap() map[string]interface{} {
//		return map[string]interface{}{
//			"I": b.I,
//			"F": b.F,
//		}
//	}
//
//	&BlockDb{
//			Id:        "test",
//			Type:      "test",
//			NoteId:    "test",
//			Order:     1,
//			CreatedAt: 2,
//			UpdatedAt: 3,
//			IsUsed:    true,
//			Data: map[string]any{
//				"test": B{
//					I: 1,
//					F: 1.0,
//				}.ToMap(),
//				"test2": B{
//					I: 2,
//					F: 2.0,
//				}.ToMap(),
//			},
func FromBlockDb(b *BlockDb) *brzrpc.Block {
	s, err := structpb.NewStruct(b.Data)
	if err != nil {
		log.Warn("mongo.FromBlockDb", "", err)
		s = nil
	}
	return &brzrpc.Block{
		Id:        b.Id,
		Type:      b.Type,
		NoteId:    b.NoteId,
		Order:     int32(b.Order),
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		IsUsed:    b.IsUsed,
		Data:      s,
	}
}
