package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"

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

	id := "test_block_"

	t.Cleanup(func() {
		assert.NoError(t, a.Delete(context.Background(), id))
		assert.NoError(t, m.Disconnect())
	})

	assert.NoError(t, a.CreateBlock(context.Background(), &domain.Block{
		Id:        id,
		Type:      "text",
		NoteId:    "",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data: map[string]any{
			"text": []any{
				map[string]any{"style": "default", "string": "123456789"},
				map[string]any{"style": "bald", "string": "01234"},
			},
		},
	}))

	b, err := a.Get(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block after create ", b)

	newData := map[string]any{
		"text": []any{
			map[string]any{"style": "default", "string": "test"},
			map[string]any{"style": "bald", "string": "test"},
		},
	}

	assert.NoError(t, a.UpdateData(context.Background(), id, newData))

	b, err = a.Get(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block after update data ", b)
	assert.Equal(t, newData, domain.ToBlockDb(domain.FromBlockDb(b)).Data)

	newType := "test"
	assert.NoError(t, a.UpdateType(context.Background(), id, newType))
	b, err = a.Get(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block after update type ", b)
	assert.Equal(t, newType, b.Type)

	newUsed := true
	assert.NoError(t, a.UpdateUsed(context.Background(), id, newUsed))
	b, err = a.Get(context.Background(), id)
	assert.NoError(t, err)
	log.Green("get block after update used ", b)
	assert.Equal(t, newUsed, b.IsUsed)

	id2 := "test_block_2"
	id3 := "test_block_3"

	assert.NoError(t, a.CreateBlock(context.Background(), &domain.Block{
		Id:        id2,
		Type:      "text",
		NoteId:    "",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data: map[string]any{
			"text": []any{
				map[string]any{"style": "default", "string": "123456789"},
				map[string]any{"style": "bald", "string": "01234"},
			},
		},
	}))

	assert.NoError(t, a.CreateBlock(context.Background(), &domain.Block{
		Id:        id3,
		Type:      "text",
		NoteId:    "",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data: map[string]any{
			"text": []any{
				map[string]any{"style": "default", "string": "123456789"},
				map[string]any{"style": "bald", "string": "01234"},
			},
		},
	}))

	assert.NoError(t, a.DeleteMany(context.Background(), []string{id2, id3}))

	_, err = a.Get(context.Background(), id2)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	_, err = a.Get(context.Background(), id3)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
