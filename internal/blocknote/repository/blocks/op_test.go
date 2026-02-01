package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/pkg/block"
	"github.com/autumnterror/breezynotes/pkg/block/default/textblock"
	"github.com/autumnterror/utils_go/pkg/log"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestOnText(t *testing.T) {
	block.RegisterBlock("text", &textblock.Driver{})
	m := mongo.MustConnect(config.Test())
	a := NewApi(m.Blocks())

	idnote := "testtextblock_note"

	var id string
	var err error

	t.Cleanup(func() {
		assert.NoError(t, a.Delete(context.Background(), id))
		assert.NoError(t, m.Disconnect())
	})

	id, err = a.Create(context.Background(), "text", idnote, map[string]any{
		"text": []any{
			map[string]any{"style": "default", "text": "123456789"},
			map[string]any{"style": "bald", "text": "01234"},
		},
	})
	assert.NoError(t, err)

	b, err := a.Get(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block after create ", b)
	log.Green("data block after create ", b.Data)

	bf, err := a.GetAsFirst(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block as first ", bf)

	assert.NoError(t, a.OpBlock(context.Background(), id, "apply_style", map[string]any{
		"start": 3,
		"end":   11,
		"style": "italic",
	}))

	b, err = a.Get(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block after op ", b)
	log.Green("data block after op ", b.Data)
}
