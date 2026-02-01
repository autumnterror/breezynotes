package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrashCycle(t *testing.T) {
	t.Parallel()
	t.Run("trash Cycle", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		id := "testIDTrashCycle"
		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        id,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			Tag: &domain.Tag{
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

		assert.NoError(t, a.ToTrash(context.Background(), id))

		_, err := a.Get(context.Background(), id)
		assert.Error(t, err)

		if nts, err := a.GetNotesFromTrash(context.Background(), "testAuthor"); assert.NoError(t, err) {
			log.Green("trash notes after create", nts)
		}

		assert.NoError(t, a.FromTrash(context.Background(), id))

		_, err = a.Get(context.Background(), id)
		assert.NoError(t, err)

		assert.NoError(t, a.ToTrash(context.Background(), id))

		assert.NoError(t, a.CleanTrash(context.Background(), "testAuthor"))
		assert.Error(t, a.FromTrash(context.Background(), id))
		if nts, err := a.GetNotesFromTrash(context.Background(), "testAuthor"); assert.NoError(t, err) {
			log.Green("trash notes after clean", nts)
		}
	})
}

func TestTrashCycleBad(t *testing.T) {
	t.Parallel()
	t.Run("trash Cycle Bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		a := NewApi(m.Notes(), m.Trash(), b)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		id := "testIDTrashCycleBad"

		assert.Error(t, a.ToTrash(context.Background(), id))
		assert.Error(t, a.FromTrash(context.Background(), id))
	})
}
