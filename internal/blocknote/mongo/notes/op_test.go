package notes

import (
	"context"
	"testing"

	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/pkgs/default/textblock"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/blocks"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/stretchr/testify/assert"
)

func TestWithBlocks(t *testing.T) {
	t.Parallel()
	t.Run("test with blocks", func(t *testing.T) {
		pkgs.RegisterBlock("text", &textblock.Driver{})
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)
		idNote := "test_with_blocks_note"
		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.TODO(), idNote))
			_, err := a.Get(context.TODO(), idNote)
			assert.Error(t, err)

			nts, err := a.getAllByUser(context.TODO(), idNote)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(nts.Items))

			assert.NoError(t, m.Disconnect())
		})
		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Note{
			Id:        idNote,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			Tag: &brzrpc.Tag{
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
		idBlock1, err := b.Create(context.TODO(), "text", idNote, map[string]any{
			"text": []any{
				map[string]any{"style": "default", "text": "test1"},
				map[string]any{"style": "bald", "text": " test2"},
			},
		})
		assert.NoError(t, err)
		idBlock2, err := b.Create(context.TODO(), "text", idNote, map[string]any{
			"text": []any{
				map[string]any{"style": "default", "text": "test3"},
				map[string]any{"style": "bald", "text": " test4"},
			},
		})

		assert.NoError(t, err)

		assert.NoError(t, a.InsertBlock(context.TODO(), idNote, idBlock1, 0))
		assert.NoError(t, a.InsertBlock(context.TODO(), idNote, idBlock2, 0))

		if n, err := a.Get(context.TODO(), idNote); assert.NoError(t, err) {
			log.Green("get after insert blocks ", n)
			assert.Equal(t, 2, len(n.Blocks))
		}

		if nts, err := a.GetNoteListByTag(context.TODO(), "test_tag", "test_auth_with_blocks"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Items)) {
			log.Green("get by tag ", nts)
		}

		if nts, err := a.GetNoteListByUser(context.TODO(), "test_auth_with_blocks"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Items)) {
			log.Green("get by user ", nts)
			if assert.Greater(t, len(nts.GetItems()), 0) {
				assert.Equal(t, "test3 test4", nts.GetItems()[0].GetFirstBlock())
			}
		}

	})
}

func TestCrudGood(t *testing.T) {
	t.Parallel()
	t.Run("crud good", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)
		id := "testIDGood"
		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.TODO(), id))
			_, err := a.Get(context.TODO(), id)
			assert.Error(t, err)

			nts, err := a.getAllByUser(context.TODO(), id)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(nts.Items))

			assert.NoError(t, m.Disconnect())
		})

		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Note{
			Id:        id,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			Tag: &brzrpc.Tag{
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

		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after create ", n)
		}

		if nts, err := a.getAllByUser(context.TODO(), "test_auth_TestCrudGood"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Items)) {
			log.Green("get all by user after create ", nts)
		}

		assert.NoError(t, a.UpdateTitle(context.TODO(), id, "new_title"))
		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after update title ", n)
		}

		assert.NoError(t, a.UpdateUpdatedAt(context.TODO(), id))
		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after update updated ", n)
		}

		idTag := "newIdTag"
		assert.NoError(t, a.AddTagToNote(context.TODO(), id, &brzrpc.Tag{
			Id:     "newIdTag",
			Title:  "newTag",
			Color:  "newColor",
			Emoji:  "newEmoji",
			UserId: "newUserId",
		}))
		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after add tag ", n)
			assert.Equal(t, "newIdTag", n.Tag.Id)
		}

		assert.NoError(t, a.RemoveTagFromNote(context.TODO(), id, idTag))
		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after remove tag ", n)
			assert.Nil(t, n.Tag)
		}
	})
}

func TestCrudNotExist(t *testing.T) {
	t.Parallel()
	t.Run("crud bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)
		id := "testIDNotExist"
		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})

		if _, err := a.Get(context.TODO(), id); assert.ErrorIs(t, err, mongo.ErrNotFound) {
		}
		if n, err := a.GetNoteListByUser(context.TODO(), id); assert.NoError(t, err) && assert.Equal(t, 0, len(n.GetItems())) {
		}
		if n, err := a.GetNoteListByTag(context.TODO(), id, id); assert.NoError(t, err) && assert.Equal(t, 0, len(n.GetItems())) {
		}
		assert.ErrorIs(t, a.UpdateTitle(context.TODO(), id, "new_title"), mongo.ErrNotFound)
		assert.ErrorIs(t, a.UpdateUpdatedAt(context.TODO(), id), mongo.ErrNotFound)
		assert.Error(t, a.Delete(context.TODO(), id))
	})
}

func TestTrashCycle(t *testing.T) {
	t.Parallel()
	t.Run("Trash Cycle", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		id := "testIDTrashCycle"
		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Note{
			Id:        id,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			Tag: &brzrpc.Tag{
				Id:     "test",
				Title:  "test",
				Color:  "test",
				Emoji:  "test",
				UserId: "test",
			},
			Author: "testAuthor",
			Editors: []string{
				"test1ed", "test2ed",
			},
			Readers: []string{
				"test1red", "test2red",
			},
			Blocks: []string{
				"test1bl", "test2bl",
			},
		}))

		assert.NoError(t, a.ToTrash(context.TODO(), id))

		_, err := a.Get(context.TODO(), id)
		assert.Error(t, err)

		if nts, err := a.GetNotesFromTrash(context.TODO(), "testAuthor"); assert.NoError(t, err) {
			log.Green("trash notes after create", nts)
		}

		assert.NoError(t, a.FromTrash(context.TODO(), id))

		_, err = a.Get(context.TODO(), id)
		assert.NoError(t, err)

		assert.NoError(t, a.ToTrash(context.TODO(), id))

		assert.NoError(t, a.CleanTrash(context.TODO(), "testAuthor"))
		assert.Error(t, a.FromTrash(context.TODO(), id))
		if nts, err := a.GetNotesFromTrash(context.TODO(), "testAuthor"); assert.NoError(t, err) {
			log.Green("trash notes after clean", nts)
		}
	})
}

func TestTrashCycleBad(t *testing.T) {
	t.Parallel()
	t.Run("Trash Cycle Bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		id := "testIDTrashCycleBad"

		assert.Error(t, a.ToTrash(context.TODO(), id))
		assert.Error(t, a.FromTrash(context.TODO(), id))
	})
}

func TestBlockOrder(t *testing.T) {
	t.Parallel()
	t.Run("test block order", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)

		id := "testblockorder"

		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.TODO(), id))
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

		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Note{
			Id:     id,
			Tag:    &brzrpc.Tag{},
			Blocks: start,
		}))
		assert.NoError(t, a.ChangeBlockOrder(context.TODO(), id, 4, 3))

		if note, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted1, note.Blocks)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.TODO(), id, 5, 1))
		if note, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted2, note.Blocks)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.TODO(), id, 0, 6))
		if note, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted3, note.Blocks)
		}
	})
}
func TestBlockOrder2(t *testing.T) {
	t.Parallel()
	t.Run("test block order 2 el", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m)
		a := NewApi(m, b)

		id := "testblockorder2"

		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.TODO(), id))
			assert.NoError(t, m.Disconnect())
		})
		start := []string{
			"test1bl", "test2bl",
		}
		wanted1 := []string{
			"test2bl", "test1bl",
		}

		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Note{
			Id:     id,
			Tag:    &brzrpc.Tag{},
			Blocks: start,
		}))
		assert.NoError(t, a.ChangeBlockOrder(context.TODO(), id, 1, 0))

		if note, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted1, note.Blocks)
		}

	})
}
