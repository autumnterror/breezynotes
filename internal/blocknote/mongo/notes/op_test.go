package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCrudGood(t *testing.T) {
	t.Parallel()
	t.Run("crud good", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		a := NewApi(m)
		id := "testIDGood"
		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.TODO(), id))
			_, err := a.Get(context.TODO(), id)
			assert.Error(t, err)
			assert.NoError(t, m.Disconnect())
		})

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
			Author: "test",
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

		assert.NoError(t, a.UpdateTitle(context.TODO(), id, "new_title"))
		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after update title ", n)
		}

		time.Sleep(1 * time.Second)

		assert.NoError(t, a.UpdateUpdatedAt(context.TODO(), id))
		if n, err := a.Get(context.TODO(), id); assert.NoError(t, err) {
			log.Green("get after update updated ", n)
		}
	})
}

func TestCrudNotExist(t *testing.T) {
	t.Parallel()
	t.Run("crud bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		a := NewApi(m)
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
		a := NewApi(m)

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

		assert.NoError(t, a.FromTrash(context.TODO(), id))

		_, err = a.Get(context.TODO(), id)
		assert.NoError(t, err)

		assert.NoError(t, a.ToTrash(context.TODO(), id))

		assert.NoError(t, a.CleanTrash(context.TODO(), "testAuthor"))
		assert.Error(t, a.FromTrash(context.TODO(), id))
	})
}

func TestTrashCycleBad(t *testing.T) {
	t.Parallel()
	t.Run("Trash Cycle Bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		a := NewApi(m)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		id := "testIDTrashCycleBad"

		assert.Error(t, a.ToTrash(context.TODO(), id))
		assert.Error(t, a.FromTrash(context.TODO(), id))
	})
}
