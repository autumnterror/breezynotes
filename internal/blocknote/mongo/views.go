package mongo

import (
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"time"
)

const (
	WaitTime = 3 * time.Second
)

type NoteDb struct {
	Id        string          `bson:"_id"`
	Title     string          `bson:"title"`
	CreatedAt int64           `bson:"createdAt"`
	UpdatedAt int64           `bson:"updatedAt"`
	Tag       *TagDb          `bson:"tag"`
	Author    string          `bson:"author"`
	Editors   []string        `bson:"editors"`
	Readers   []string        `bson:"readers"`
	Blocks    []string        `bson:"blocks"`
	Status    brzrpc.Statuses `bson:"status"`
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

type TagDb struct {
	Id     string `bson:"_id"`
	Title  string `bson:"title"`
	Color  string `bson:"color"`
	Emoji  string `bson:"emoji"`
	UserId string `bson:"userId"`
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
