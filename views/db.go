package views

import (
	"fmt"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"

	"github.com/autumnterror/breezynotes/pkg/log"
	"google.golang.org/protobuf/types/known/structpb"
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
	nn := brzrpc.Note{Blocks: n.Blocks, Editors: n.Editors, Readers: n.Readers}
	if n.Blocks == nil {
		nn.Blocks = []string{}
	}
	if n.Editors == nil {
		nn.Editors = []string{}
	}
	if n.Readers == nil {
		nn.Readers = []string{}
	}
	return &NoteDb{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       ToTagDb(n.Tag),
		Author:    n.Author,
		Editors:   nn.Editors,
		Readers:   nn.Readers,
		Blocks:    nn.Blocks,
	}
}

func FromNoteDb(n *NoteDb) *brzrpc.Note {
	if n == nil {
		return nil
	}
	nn := brzrpc.Note{Blocks: n.Blocks, Editors: n.Editors, Readers: n.Readers}
	if n.Blocks == nil {
		nn.Blocks = []string{}
	}
	if n.Editors == nil {
		nn.Editors = []string{}
	}
	if n.Readers == nil {
		nn.Readers = []string{}
	}
	return &brzrpc.Note{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       FromTagDb(n.Tag),
		Author:    n.Author,
		Editors:   nn.Editors,
		Readers:   nn.Readers,
		Blocks:    nn.Blocks,
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

func normalize(v any) any {
	switch x := v.(type) {

	case bson.M:
		m := make(map[string]any, len(x))
		for k, v2 := range x {
			m[k] = normalize(v2)
		}
		return m
	case bson.D:
		m := make(map[string]any, len(x))
		for _, e := range x {
			m[e.Key] = normalize(e.Value)
		}
		return m
	case bson.A:
		s := make([]any, len(x))
		for i, v2 := range x {
			s[i] = normalize(v2)
		}
		return s

	case map[string]any:
		for k, v2 := range x {
			x[k] = normalize(v2)
		}
		return x

	case []any:
		for i, v2 := range x {
			x[i] = normalize(v2)
		}
		return x

	default:
		return v
	}
}

// FromBlockDb on data field u can insert only base type. If u want ur struct convert to map[string]any (check models_test.go)
func FromBlockDb(b *BlockDb) *brzrpc.Block {
	normalized := normalize(b.Data)

	m, ok := normalized.(map[string]any)
	if !ok {
		log.Warn("mongo.FromBlockDb", "", fmt.Errorf("data is not a map[string]any after normalize, got %T", normalized))
		m = nil
	}

	s, err := structpb.NewStruct(m)
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
