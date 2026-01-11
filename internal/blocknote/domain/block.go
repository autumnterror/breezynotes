package domain

import (
	"fmt"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Block struct {
	Id     string `bson:"_id"`
	Type   string `bson:"type"`
	NoteId string `bson:"note_id"`

	CreatedAt int64          `bson:"created_at"`
	UpdatedAt int64          `bson:"updated_at"`
	IsUsed    bool           `bson:"is_used"`
	Data      map[string]any `bson:"data"`
}

type Blocks struct {
	Blks []*Block
}

func ToBlockDb(b *brzrpc.Block) *Block {
	return &Block{
		Id:        b.GetId(),
		Type:      b.GetType(),
		NoteId:    b.GetNoteId(),
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

// FromBlockDb on data field u can insert only base type.
// If u want ur struct convert to map[string]any (check models_test.go)
func FromBlockDb(b *Block) *brzrpc.Block {
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
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		IsUsed:    b.IsUsed,
		Data:      s,
	}
}

func ToBlocksDb(b *brzrpc.Blocks) *Blocks {
	if b == nil {
		return nil
	}

	var blks []*Block
	for _, blk := range b.GetItems() {
		blks = append(blks, ToBlockDb(blk))
	}

	return &Blocks{
		Blks: blks,
	}
}

func FromBlocksDb(b *Blocks) *brzrpc.Blocks {
	if b == nil {
		return nil
	}

	var blks []*brzrpc.Block
	for _, blk := range b.Blks {
		blks = append(blks, FromBlockDb(blk))
	}

	return &brzrpc.Blocks{
		Items: blks,
	}
}
