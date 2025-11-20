package tags

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrudGood(t *testing.T) {
	t.Parallel()
	t.Run("crud good", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		a := NewApi(m)

		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Tag{
			Id:     "test_id_good",
			Title:  "test_title",
			Color:  "test_color",
			Emoji:  "test_emoji",
			UserId: "test_userid",
		}))

		if tgs, err := a.GetAllById(context.TODO(), "test_userid"); assert.NoError(t, err) {
			log.Green("get after create", tgs.Items)
		}

		assert.NoError(t, a.UpdateColor(context.TODO(), "test_id_good", "new_color"))
		assert.NoError(t, a.UpdateTitle(context.TODO(), "test_id_good", "new_title"))
		assert.NoError(t, a.UpdateEmoji(context.TODO(), "test_id_good", "new_emoji"))

		if tgs, err := a.GetAllById(context.TODO(), "test_userid"); assert.NoError(t, err) {
			log.Green("get after update", tgs.Items)
		}

		_, err := a.Get(context.TODO(), "test_id_good")
		assert.NoError(t, err)

		t.Cleanup(func() {
			assert.NoError(t, a.Delete(context.TODO(), "test_id_good"))
			if tgs, err := a.GetAllById(context.TODO(), "test_userid"); assert.NoError(t, err) {
				log.Green("get after delete", tgs.Items)
			}
			assert.NoError(t, m.Disconnect())
		})
	})
}

func TestBad(t *testing.T) {
	t.Parallel()
	m := mongo.MustConnect(config.Test())
	a := NewApi(m)
	t.Cleanup(func() {
		assert.NoError(t, m.Disconnect())
	})

	t.Run("BAD create repeat id", func(t *testing.T) {
		assert.NoError(t, a.Create(context.TODO(), &brzrpc.Tag{
			Id:     "test_id",
			Title:  "test_title",
			Color:  "test_color",
			Emoji:  "test_emoji",
			UserId: "test_userid",
		}))
		assert.Error(t, a.Create(context.TODO(), &brzrpc.Tag{
			Id:     "test_id",
			Title:  "test_title",
			Color:  "test_color",
			Emoji:  "test_emoji",
			UserId: "test_userid",
		}))

		assert.NoError(t, a.Delete(context.TODO(), "test_id"))
	})

	t.Run("BAD update ErrNotFound", func(t *testing.T) {
		assert.ErrorIs(t, a.UpdateColor(context.TODO(), "test_id", "new_color"), mongo.ErrNotFound)
		assert.ErrorIs(t, a.UpdateTitle(context.TODO(), "test_id", "new_title"), mongo.ErrNotFound)
		assert.ErrorIs(t, a.UpdateEmoji(context.TODO(), "test_id", "new_emoji"), mongo.ErrNotFound)
	})
}
