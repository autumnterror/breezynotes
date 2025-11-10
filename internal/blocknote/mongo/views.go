package mongo

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
	Id        string          `bson:"_id"`
	Title     string          `bson:"title"`
	CreatedAt int64           `bson:"created_at"`
	UpdatedAt int64           `bson:"updated_at"`
	Tag       *TagDb          `bson:"tag"`
	Author    string          `bson:"author"`
	Editors   []string        `bson:"editors"`
	Readers   []string        `bson:"readers"`
	Blocks    []string        `bson:"blocks"`
	Status    brzrpc.Statuses `bson:"status"`
}

type TagDb struct {
	Id     string `bson:"_id"`
	Title  string `bson:"title"`
	Color  string `bson:"color"`
	Emoji  string `bson:"emoji"`
	UserId string `bson:"user_id"`
}

func ToNoteDb(n *brzrpc.Note) *NoteDb {
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
	return &TagDb{
		Id:     t.Id,
		Title:  t.Title,
		Color:  t.Color,
		Emoji:  t.Emoji,
		UserId: t.UserId,
	}
}

func FromTagDb(t *TagDb) *brzrpc.Tag {
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
	return &BlockDb{}
}
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
