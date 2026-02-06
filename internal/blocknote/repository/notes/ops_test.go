package notes

import (
	"context"

	"github.com/autumnterror/utils_go/pkg/utils/uid"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/pkg/block"
	"github.com/autumnterror/breezynotes/pkg/block/default/textblock"

	"testing"

	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestWithBlocks(t *testing.T) {
	t.Parallel()
	t.Run("test with blocks", func(t *testing.T) {
		block.RegisterBlock("text", &textblock.Driver{})
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)
		idNote := "test_with_blocks_note"
		t.Cleanup(func() {
			assert.NoError(t, a.delete(context.Background(), idNote))
			_, err := a.Get(context.Background(), idNote)
			assert.Error(t, err)

			nts, err := a.getAllByUser(context.Background(), idNote)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(nts.Nts))

			assert.NoError(t, m.Disconnect())
		})
		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        idNote,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			Tag: &domain.Tag{
				Id:     "test_tag",
				Title:  "test",
				Color:  "test",
				Emoji:  "test",
				UserId: "test",
			},
			Author: "test_auth_with_blocks",
			Editors: []string{
				"test1ed", "test2ed",
			},
			Readers: []string{
				"test1red", "test2red",
			},
			Blocks: []string{},
		}))
		idBlock1 := uid.New()
		idBlock2 := uid.New()

		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:        idBlock1,
			Type:      "text",
			NoteId:    idNote,
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "test1"},
					map[string]any{"style": "bald", "string": " test2"},
				},
			},
		}))
		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:        idBlock2,
			Type:      "text",
			NoteId:    idNote,
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "test3"},
					map[string]any{"style": "bald", "string": " test4"},
				},
			},
		}))

		assert.NoError(t, a.InsertBlock(context.Background(), idNote, idBlock1, 0))
		assert.NoError(t, a.InsertBlock(context.Background(), idNote, idBlock2, 0))

		if n, err := a.Get(context.Background(), idNote); assert.NoError(t, err) {
			log.Green("get after insert blocks ", n)
			assert.Equal(t, 2, len(n.Blocks))
		}

		assert.NoError(t, a.DeleteBlock(context.Background(), idNote, idBlock1))
		if n, err := a.Get(context.Background(), idNote); assert.NoError(t, err) {
			log.Green("get after delete block ", n)
			assert.Equal(t, 1, len(n.Blocks))
		}

		if nts, err := a.GetNoteListByTag(context.Background(), "test_tag", "test_auth_with_blocks"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Ntps)) {
			log.Green("get by tag ", nts)
		}

		if nts, err := a.GetNoteListByUser(context.Background(), "test_auth_with_blocks"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Ntps)) {
			log.Green("get by user ", nts)
			if assert.Greater(t, len(nts.Ntps), 0) {
				assert.Equal(t, "test3 test4", nts.Ntps[0].FirstBlock)
			}
		}

	})
}

func TestCrudGood(t *testing.T) {
	t.Parallel()
	t.Run("crud good", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)
		id := "testIDGood"
		t.Cleanup(func() {
			assert.NoError(t, a.delete(context.Background(), id))
			_, err := a.Get(context.Background(), id)
			assert.Error(t, err)

			nts, err := a.getAllByUser(context.Background(), id)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(nts.Nts))

			assert.NoError(t, m.Disconnect())
		})

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        id,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			Tag: &domain.Tag{
				Id:     "test_tag",
				Title:  "test",
				Color:  "test",
				Emoji:  "test",
				UserId: "test",
			},
			Author: "test_auth_TestCrudGood",
			Editors: []string{
				"test1ed", "test2ed",
			},
			Readers: []string{
				"test1red", "test2red",
			},
			Blocks: []string{
				"test1", "test2",
			},
		}))

		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get after create ", n)
		}

		if nts, err := a.getAllByUser(context.Background(), "test_auth_TestCrudGood"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Nts)) {
			log.Green("get all by user after create ", nts)
		}

		assert.NoError(t, a.UpdateTitle(context.Background(), id, "new_title"))
		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get after update title ", n)
		}

		assert.NoError(t, a.UpdateUpdatedAt(context.Background(), id))
		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get after update updated ", n)
		}

		assert.NoError(t, a.AddTagToNote(context.Background(), id, &domain.Tag{
			Id:     "newIdTag",
			Title:  "newTag",
			Color:  "newColor",
			Emoji:  "newEmoji",
			UserId: "newUserId",
		}))
		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get after add tag ", n)
			assert.Equal(t, "newIdTag", n.Tag.Id)
		}

		assert.NoError(t, a.RemoveTagFromNote(context.Background(), id))
		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get after remove tag ", n)
			assert.Nil(t, n.Tag)
		}

		assert.NoError(t, a.ShareNote(context.Background(), id, "neweditor", domain.EditorRole))
		assert.NoError(t, a.ShareNote(context.Background(), id, "newreader", domain.ReaderRole))
		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get share note ", n)
			assert.Contains(t, n.Readers, "newreader")
			assert.Contains(t, n.Editors, "neweditor")
		}

		n, err := a.GetNoteListByUser(context.Background(), "neweditor")
		assert.NoError(t, err)
		assert.Len(t, n.Ntps, 1)

		n, err = a.GetNoteListByUser(context.Background(), "newreader")
		assert.NoError(t, err)
		assert.Len(t, n.Ntps, 1)

		assert.NoError(t, a.ChangeUserRole(context.Background(), id, "neweditor", domain.ReaderRole))
		assert.NoError(t, a.ChangeUserRole(context.Background(), id, "newreader", domain.EditorRole))
		if n, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Green("get ChangeUserRole ", n)
			assert.Contains(t, n.Editors, "newreader")
			assert.Contains(t, n.Readers, "neweditor")
		}
	})
}

func TestCrudNotExist(t *testing.T) {
	t.Parallel()
	t.Run("crud bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)
		id := "testIDNotExist"
		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})

		if _, err := a.Get(context.Background(), id); assert.ErrorIs(t, err, domain.ErrNotFound) {
		}
		if n, err := a.GetNoteListByUser(context.Background(), id); assert.NoError(t, err) && assert.Equal(t, 0, len(n.Ntps)) {
		}
		if n, err := a.GetNoteListByTag(context.Background(), id, id); assert.NoError(t, err) && assert.Equal(t, 0, len(n.Ntps)) {
		}
		assert.ErrorIs(t, a.UpdateTitle(context.Background(), id, "new_title"), domain.ErrNotFound)
		assert.ErrorIs(t, a.UpdateUpdatedAt(context.Background(), id), domain.ErrNotFound)
		assert.Error(t, a.delete(context.Background(), id))
		assert.ErrorIs(t, a.ChangeUserRole(context.Background(), id, "neweditor", domain.ReaderRole), domain.ErrNotFound)
		assert.ErrorIs(t, a.ShareNote(context.Background(), id, "newreader", domain.ReaderRole), domain.ErrNotFound)
	})
}

func TestBlockOrder(t *testing.T) {
	t.Parallel()
	t.Run("test block order", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)

		id := "testblockorder"

		t.Cleanup(func() {
			assert.NoError(t, a.delete(context.Background(), id))
			assert.NoError(t, m.Disconnect())
		})
		start := []string{
			"test1bl", "test2bl", "test3bl", "test4bl", "test5bl", "test6bl",
		}
		wanted1 := []string{
			"test1bl", "test2bl", "test3bl", "test5bl", "test4bl", "test6bl",
		}
		wanted2 := []string{
			"test1bl", "test6bl", "test2bl", "test3bl", "test5bl", "test4bl",
		}
		wanted3 := []string{
			"test6bl", "test2bl", "test3bl", "test5bl", "test4bl", "test1bl",
		}

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:     id,
			Tag:    &domain.Tag{},
			Blocks: start,
		}))
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 4, 3))

		if note, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted1, note.Blocks)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 5, 1))
		if note, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted2, note.Blocks)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 0, 6))
		if note, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted3, note.Blocks)
		}
	})
}
func TestBlockOrder2(t *testing.T) {
	t.Parallel()
	t.Run("test block order 2 el", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)
		id := "testblockorder2"

		t.Cleanup(func() {
			assert.NoError(t, a.delete(context.Background(), id))
			assert.NoError(t, m.Disconnect())
		})
		start := []string{
			"test1bl", "test2bl",
		}
		wanted1 := []string{
			"test2bl", "test1bl",
		}

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:     id,
			Tag:    &domain.Tag{},
			Blocks: start,
		}))
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 1, 0))

		if note, err := a.Get(context.Background(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted1, note.Blocks)
		}
	})
}
