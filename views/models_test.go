package views

import (
	"github.com/autumnterror/breezynotes/pkg/log"
	"testing"
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
	bdb := &BlockDb{
		Id:        "test",
		Type:      "test",
		NoteId:    "test",
		Order:     1,
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
	log.Blue(*brzb)
	log.Blue(*ToBlockDb(brzb))
}
