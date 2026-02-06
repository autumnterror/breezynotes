package service

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo/mongotx"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/textblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCascadeBlocks(t *testing.T) {
	block.RegisterBlock("text", &textblock.Driver{})
	m := mongo.MustConnect(config.Test())
	b := blocks.NewApi(m.Blocks())
	tgs := tags.NewApi(m.Tags())
	n := notes.NewApi(m.Notes(), m.Trash(), b)
	s := NewNoteService(config.Test(), mongotx.NewTxRunner(m.C), n, b, tgs)

	nId := uid.New()
	idUser := uid.New()
	assert.NoError(t, s.CreateNote(ctx, &domain.Note{
		Id:        nId,
		Title:     "test",
		CreatedAt: 0,
		UpdatedAt: 0,
		Tag:       nil,
		Author:    idUser,
		Editors:   nil,
		Readers:   nil,
		Blocks:    nil,
	}))

	blockId1, err := s.CreateBlock(ctx, "text", nId, map[string]any{
		"text": []any{
			map[string]any{"style": "default", "string": "123456789"},
			map[string]any{"style": "bald", "string": "01234"},
		},
	}, 0, idUser)
	assert.NoError(t, err)
	blockId2, err := s.CreateBlock(ctx, "text", nId, map[string]any{
		"text": []any{
			map[string]any{"style": "default", "string": "123456789"},
			map[string]any{"style": "bald", "string": "01234"},
		},
	}, 0, idUser)
	assert.NoError(t, err)

	assert.NoError(t, s.ToTrash(ctx, nId, idUser))
	assert.NoError(t, s.CleanTrash(ctx, idUser))
	_, err = s.blk.Get(ctx, blockId1)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	_, err = s.blk.Get(ctx, blockId2)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

var ctx = context.Background()
