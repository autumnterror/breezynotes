package tags

import (
	"context"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"

	"testing"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	mongo2 "github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestCrudGood(t *testing.T) {
	t.Parallel()
	t.Run("crud good", func(t *testing.T) {
		m := mongo2.MustConnect(config.Test())
		a := NewApi(m.Tags(), m.NoteTags())

		assert.NoError(t, a.Create(context.Background(), &domain2.Tag{
			Id:     "test_id_good",
			Title:  "test_title",
			Color:  "test_color",
			Emoji:  "test_emoji",
			UserId: "test_userid",
		}))

		if tgs, err := a.GetAllById(context.Background(), "test_userid"); assert.NoError(t, err) {
			log.Green("get after create", tgs.Tgs)
		}

		assert.NoError(t, a.UpdateColor(context.Background(), "test_id_good", "new_color"))
		assert.NoError(t, a.UpdateTitle(context.Background(), "test_id_good", "new_title"))
		assert.NoError(t, a.UpdateEmoji(context.Background(), "test_id_good", "new_emoji"))

		if tgs, err := a.GetAllById(context.Background(), "test_userid"); assert.NoError(t, err) {
			log.Green("get after update", tgs.Tgs)
		}

		_, err := a.Get(context.Background(), "test_id_good")
		assert.NoError(t, err)

		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.Background(), "test_id_good"))
			if tgs, err := a.GetAllById(context.Background(), "test_userid"); assert.NoError(t, err) {
				log.Green("get after delete", tgs.Tgs)
			}
			assert.NoError(t, m.Disconnect())
		})
	})
}

func TestBad(t *testing.T) {
	t.Parallel()
	m := mongo2.MustConnect(config.Test())
	a := NewApi(m.Tags(), m.NoteTags())
	t.Cleanup(func() {
		assert.NoError(t, m.Disconnect())
	})

	t.Run("BAD create repeat id", func(t *testing.T) {
		assert.NoError(t, a.Create(context.Background(), &domain2.Tag{
			Id:     "test_id",
			Title:  "test_title",
			Color:  "test_color",
			Emoji:  "test_emoji",
			UserId: "test_userid",
		}))
		assert.Error(t, a.Create(context.Background(), &domain2.Tag{
			Id:     "test_id",
			Title:  "test_title",
			Color:  "test_color",
			Emoji:  "test_emoji",
			UserId: "test_userid",
		}))

		assert.NoError(t, a.Delete(context.Background(), "test_id"))
	})

	t.Run("BAD ErrNotFound", func(t *testing.T) {
		assert.ErrorIs(t, a.UpdateColor(context.Background(), "test_id", "new_color"), domain2.ErrNotFound)
		assert.ErrorIs(t, a.UpdateTitle(context.Background(), "test_id", "new_title"), domain2.ErrNotFound)
		assert.ErrorIs(t, a.UpdateEmoji(context.Background(), "test_id", "new_emoji"), domain2.ErrNotFound)
		assert.ErrorIs(t, a.Delete(context.Background(), "test_id"), domain2.ErrNotFound)
		_, err := a.Get(context.Background(), "test_id")
		assert.ErrorIs(t, err, domain2.ErrNotFound)
	})

}
