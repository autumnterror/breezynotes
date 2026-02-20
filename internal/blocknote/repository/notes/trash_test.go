package notes

import (
	"context"
	"testing"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/stretchr/testify/assert"
)

func TestTrashCycle(t *testing.T) {
	t.Run("trash Cycle", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		idNote := uid.New()
		idUser := uid.New()
		idTag := uid.New()
		n := &domain.Note{
			Id:        idNote,
			Title:     "test",
			CreatedAt: time.Now().UTC().Unix(),
			UpdatedAt: time.Now().UTC().Unix(),
			// Tag: &domain2.Tag{
			// 	Id:     "test",
			// 	Title:  "test",
			// 	Color:  "test",
			// 	Emoji:  "test",
			// 	UserId: "test",
			// },
			Author: idUser,
			Editors: []string{
				"test1ed", "test2ed",
			},
			Readers: []string{
				"test1red", "test2red",
			},
			Blocks: []string{
				"test1bl", "test2bl",
			},
		}

		assert.NoError(t, a.Create(context.Background(), n))

		newTag := &domain.Tag{
			Id:     idTag,
			Title:  "newTagTRASH",
			Color:  "newColor",
			Emoji:  "newEmoji",
			UserId: idUser,
		}

		assert.NoError(t, tgs.Create(context.Background(), newTag))
		assert.NoError(t, a.AddTagToNote(context.Background(), idNote, newTag))
		if nts, err := a.GetNoteListByTag(context.Background(), idTag, idUser); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Ntps)) {
			log.Green("get by tag ", nts)
		}

		assert.NoError(t, a.ToTrash(context.Background(), idNote))

		_, err := a.Get(context.Background(), idNote, "")
		assert.Error(t, err)

		if nts, err := a.GetNotesFromTrash(context.Background(), idUser); assert.NoError(t, err) {
			log.Green("trash notes after create", nts)
			assert.Greater(t, len(nts.Ntps), 0)
			assert.Equal(t, newTag, nts.Ntps[0].Tag)
		}
		if nts, err := a.GetNotesFullFromTrash(context.Background(), idUser); assert.NoError(t, err) {
			assert.Equal(t, n, nts.Nts[0])
		}

		assert.NoError(t, a.FromTrash(context.Background(), idNote))

		_, err = a.Get(context.Background(), idNote, "")
		assert.NoError(t, err)

		assert.NoError(t, a.ToTrash(context.Background(), idNote))

		assert.NoError(t, a.CleanTrash(context.Background(), idUser))
		assert.Error(t, a.FromTrash(context.Background(), idNote))
		if nts, err := a.GetNotesFromTrash(context.Background(), idUser); assert.NoError(t, err) {
			log.Green("trash notes after clean", nts)
		}
	})
}

func TestTrashCycleBad(t *testing.T) {
	t.Run("trash Cycle Bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)

		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})
		id := uid.New()

		assert.Error(t, a.ToTrash(context.Background(), id))
		assert.Error(t, a.FromTrash(context.Background(), id))
	})
}
