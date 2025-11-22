package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/pkgs/default/textblock"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestOnText(t *testing.T) {
	pkgs.RegisterBlock("text", &textblock.Driver{})
	m := mongo.MustConnect(config.Test())
	a := NewApi(m)

	idnote := "testtextblock_note"

	var id string
	var err error

	t.Cleanup(func() {
		assert.NoError(t, a.Delete(context.TODO(), id))
		assert.NoError(t, m.Disconnect())
	})

	id, err = a.Create(context.TODO(), "text", idnote, map[string]any{
		"text": []any{
			map[string]any{"style": "default", "text": "123456789"},
			map[string]any{"style": "bald", "text": "01234"},
		},
	})
	assert.NoError(t, err)

	b, err := a.Get(context.TODO(), id)
	assert.NoError(t, err)
	log.Green("get block after create ", b)
	log.Green("data block after create ", b.Data.AsMap())

	bf, err := a.GetAsFirst(context.TODO(), id)
	assert.NoError(t, err)
	log.Green("get block as first ", bf)

	assert.NoError(t, a.OpBlock(context.TODO(), id, "apply_style", map[string]any{
		"start": 3,
		"end":   11,
		"style": "italic",
	}))

	b, err = a.Get(context.TODO(), id)
	assert.NoError(t, err)
	log.Green("get block after op ", b)
	log.Green("data block after op ", b.Data.AsMap())
}
