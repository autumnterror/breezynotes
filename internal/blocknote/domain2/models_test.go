package domain2

import (
	"testing"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
)

type B struct {
	I int
	F float32
}

func (b B) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"I": b.I,
		"F": b.F,
	}
}

func TestBlockDb(t *testing.T) {
	bdb := &Block{
		Id:        "test",
		Type:      "test",
		NoteId:    "test",
		CreatedAt: 2,
		UpdatedAt: 3,
		IsUsed:    true,
		Data: map[string]any{
			"test": B{
				I: 1,
				F: 1.0,
			}.ToMap(),
			"test2": B{
				I: 2,
				F: 2.0,
			}.ToMap(),
		},
	}
	brzb := FromBlockDb(bdb)
	log.Blue(format.Struct(*brzb))
	log.Blue(format.Struct(*ToBlockDb(brzb)))
}

func TestBlockDb2(t *testing.T) {
	bdb := &Block{
		Id:        "test",
		Type:      "test",
		NoteId:    "test",
		CreatedAt: 2,
		UpdatedAt: 3,
		IsUsed:    true,
		Data: map[string]any{
			"text": []any{
				map[string]any{"style": "default", "text": "hello its example!"},
				map[string]any{"style": "default", "text": "hello its example!"},
			},
		},
	}
	brzb := FromBlockDb(bdb)
	log.Blue(format.Struct(*brzb))
	log.Blue(format.Struct(*ToBlockDb(brzb)))
}
