package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/blocks"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO GetNoteListByUser
//
//	if np, err := a.GetNoteListByUser(context.TODO(), "test_auth_TestCrudGood"); assert.NoError(t, err) {
//		assert.Equal(t, len(np.Items), 1)
//		assert.Equal(t, np.Items[0]., 1)
//	}
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

			nts, err := a.GetAllByUser(context.TODO(), id)
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
				"test1bl", "test2bl",
			},
			Status: brzrpc.Statuses_IN_USE,
		}))

		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after create ", n)
		}

		if nts, err := a.GetAllByTag(context.TODO(), "test_tag"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Items)) {
			log.Green("get by tag ", nts)
		}

		if nts, err := a.GetAllByUser(context.TODO(), "test_auth"); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Items)) {
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

		_, err := a.Get(context.TODO(), id)
		assert.Error(t, err)

		assert.ErrorIs(t, a.UpdateTitle(context.TODO(), id, "new_title"), mongo.ErrNotFiend)
		assert.ErrorIs(t, a.UpdateUpdatedAt(context.TODO(), id), mongo.ErrNotFiend)
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
			Status: brzrpc.Statuses_IN_USE,
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
			log.Println(note)
			assert.Equal(t, note.Blocks, wanted1)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.TODO(), id, 5, 1))
		if note, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Println(note)
			assert.Equal(t, note.Blocks, wanted2)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.TODO(), id, 0, 6))
		if note, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Println(note)
			assert.Equal(t, note.Blocks, wanted3)
		}
	})
}
