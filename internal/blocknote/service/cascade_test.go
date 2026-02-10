package service

import (
	"testing"
)

func TestCascadeBlocks(t *testing.T) {
	//block.RegisterBlock("text", &textblock.Driver{})
	//m := mongo.MustConnect(config.Test())
	//b := blocks.NewApi(m.Blocks())
	//tgs := tags.NewApi(m.Tags())
	//n := notes.NewApi(m.Notes(), m.Trash(), b)
	//s := NewNoteService(config.Test(), mongotx.NewTxRunner(m.C), n, b, tgs)
	//
	//nId := uid.New()
	//idUser := uid.New()
	//assert.NoError(t, s.CreateNote(ctx, &domain2.Note{
	//	Id:        nId,
	//	Title:     "test",
	//	CreatedAt: 0,
	//	UpdatedAt: 0,
	//	Tag:       nil,
	//	Author:    idUser,
	//	Editors:   nil,
	//	Readers:   nil,
	//	Blocks:    nil,
	//}))
	//
	//blockId1, err := s.CreateBlock(ctx, "text", nId, map[string]any{
	//	"text": []any{
	//		map[string]any{"style": "default", "string": "123456789"},
	//		map[string]any{"style": "bald", "string": "01234"},
	//	},
	//}, 0, idUser)
	//assert.NoError(t, err)
	//blockId2, err := s.CreateBlock(ctx, "text", nId, map[string]any{
	//	"text": []any{
	//		map[string]any{"style": "default", "string": "123456789"},
	//		map[string]any{"style": "bald", "string": "01234"},
	//	},
	//}, 0, idUser)
	//assert.NoError(t, err)
	//
	//assert.NoError(t, s.ToTrash(ctx, nId, idUser))
	//assert.NoError(t, s.CleanTrash(ctx, idUser))
	//_, err = s.blk.Get(ctx, blockId1)
	//assert.ErrorIs(t, err, domain2.ErrNotFound)
	//_, err = s.blk.Get(ctx, blockId2)
	//assert.ErrorIs(t, err, domain2.ErrNotFound)
}

//var ctx = context.Background()
