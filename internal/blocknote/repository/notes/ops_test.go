package notes

import (
	"context"

	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/textblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"
	"github.com/autumnterror/utils_go/pkg/utils/format"

	"github.com/autumnterror/utils_go/pkg/utils/uid"

	"testing"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"

	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestWithBlocks(t *testing.T) {

	t.Run("test with blocks", func(t *testing.T) {
		block.RegisterBlock("text", &textblock.Driver{})
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)
		idNote := uid.New()
		idUser := uid.New()
		t.Cleanup(func() {
			assert.NoError(t, a.delete(context.Background(), idNote))
			_, err := a.Get(context.Background(), idNote, idUser)
			assert.Error(t, err)

			nts, err := a.getAllByUser(context.Background(), idNote)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(nts.Nts))

			assert.NoError(t, m.Disconnect())
		})
		assert.NoError(t, a.Create(context.Background(), &domain2.Note{
			Id:        idNote,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			//Tag: &domain2.Tag{
			//	Id:     "test_tag",
			//	Title:  "test",
			//	Color:  "test",
			//	Emoji:  "test",
			//	UserId: "test",
			//},
			Author: idUser,
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

		assert.NoError(t, b.CreateBlock(context.Background(), &domain2.Block{
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
		assert.NoError(t, b.CreateBlock(context.Background(), &domain2.Block{
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

		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after insert blocks ", n)
			assert.Equal(t, 2, len(n.Blocks))
		}

		assert.NoError(t, a.DeleteBlock(context.Background(), idNote, idBlock1))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after delete block ", n)
			assert.Equal(t, 1, len(n.Blocks))
		}

		idTag := uid.New()
		newTag := &domain2.Tag{
			Id:     idTag,
			Title:  "newTag",
			Color:  "newColor",
			Emoji:  "newEmoji",
			UserId: idUser,
		}

		assert.NoError(t, tgs.Create(context.Background(), newTag))

		assert.NoError(t, a.AddTagToNote(context.Background(), idNote, newTag))

		if nts, err := a.GetNoteListByTag(context.Background(), idTag, idUser); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Ntps)) {
			log.Green("get by tag ", nts)
		}

		if nts, err := a.GetNoteListByUser(context.Background(), idUser); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Ntps)) {
			log.Green("get by user ", nts)
			if assert.Greater(t, len(nts.Ntps), 0) {
				assert.Equal(t, "test3 test4", nts.Ntps[0].FirstBlock)
			}
		}

	})
}

func TestCrudGood(t *testing.T) {

	t.Run("crud good", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)
		idNote := uid.New()
		idUser := uid.New()
		idTag := uid.New()
		t.Cleanup(func() {
			assert.NoError(t, a.delete(context.Background(), idNote))
			_, err := a.Get(context.Background(), idNote, idUser)
			assert.Error(t, err)

			nts, err := a.getAllByUser(context.Background(), idNote)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(nts.Nts))

			assert.NoError(t, tgs.Delete(context.Background(), idTag))

			assert.NoError(t, m.Disconnect())
		})

		assert.NoError(t, a.Create(context.Background(), &domain2.Note{
			Id:        idNote,
			Title:     "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			//Tag: &domain2.Tag{
			//	Id:     "test_tag",
			//	Title:  "test",
			//	Color:  "test",
			//	Emoji:  "test",
			//	UserId: "test",
			//},
			Author: idUser,
			Editors: []string{
				"test1ed", "test2ed",
			},
			Readers: []string{
				"test1red", "test2red",
			},
			Blocks: []string{
				"test1", "test2",
			},
			IsBlog:   false,
			IsPublic: false,
		}))

		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after create ", n)
		}

		if nts, err := a.getAllByUser(context.Background(), idUser); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Nts)) {
			log.Green("get all by user after create ", nts)
		}

		assert.NoError(t, a.UpdateTitle(context.Background(), idNote, "new_title"))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after update title ", n)
		}

		assert.NoError(t, a.UpdateUpdatedAt(context.Background(), idNote))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after update updated ", n)
		}

		newTag := &domain2.Tag{
			Id:     idTag,
			Title:  "newTag",
			Color:  "newColor",
			Emoji:  "newEmoji",
			UserId: idUser,
		}

		assert.NoError(t, tgs.Create(context.Background(), newTag))

		assert.NoError(t, a.AddTagToNote(context.Background(), idNote, newTag))
		if nts, err := a.GetNoteListByTag(context.Background(), idTag, idUser); assert.NoError(t, err) && assert.NotEqual(t, 0, len(nts.Ntps)) {
			log.Green("get by tag ", format.Struct(nts))
		}
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after add tag ", n)
			assert.Equal(t, newTag, n.Tag)
		}

		assert.NoError(t, a.RemoveTagFromNote(context.Background(), idNote, idUser))
		if nts, err := a.GetNoteListByTag(context.Background(), idTag, idUser); assert.NoError(t, err) && assert.Equal(t, 0, len(nts.Ntps)) {
			log.Green("get by tag ", format.Struct(nts))
		}
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get after rm tag ", n)
			assert.Nil(t, n.Tag)
		}
		assert.NoError(t, a.ShareNote(context.Background(), idNote, "neweditor", domain2.EditorRole))
		assert.NoError(t, a.ShareNote(context.Background(), idNote, "newreader", domain2.ReaderRole))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
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

		assert.NoError(t, a.ShareNote(context.Background(), idNote, "neweditor", domain2.ReaderRole))
		assert.NoError(t, a.ShareNote(context.Background(), idNote, "newreader", domain2.EditorRole))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get ChangeUserRole ", n)
			assert.Contains(t, n.Editors, "newreader")
			assert.Contains(t, n.Readers, "neweditor")
		}

		assert.NoError(t, a.UpdateBlog(context.Background(), idNote, true))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get UpdateBlogOrPublic ", n)
			assert.Equal(t, n.IsBlog, true)
			assert.Equal(t, n.IsPublic, false)
		}
		assert.NoError(t, a.UpdatePublic(context.Background(), idNote, true))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get UpdateBlogOrPublic ", n)
			assert.Equal(t, n.IsBlog, true)
			assert.Equal(t, n.IsPublic, true)
		}

		assert.NoError(t, a.DeleteRole(context.Background(), idNote, "neweditor"))
		assert.NoError(t, a.DeleteRole(context.Background(), idNote, "newreader"))
		if n, err := a.Get(context.Background(), idNote, idUser); assert.NoError(t, err) {
			log.Green("get DeleteRole ", n)
			assert.NotContains(t, n.Editors, "newreader")
			assert.NotContains(t, n.Editors, "neweditor")
			assert.NotContains(t, n.Readers, "newreader")
			assert.NotContains(t, n.Readers, "neweditor")
		}

	})
}

func TestCrudNotExist(t *testing.T) {

	t.Run("crud bad", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)
		id := uid.New()
		t.Cleanup(func() {
			assert.NoError(t, m.Disconnect())
		})

		if _, err := a.Get(context.Background(), id, ""); assert.ErrorIs(t, err, domain2.ErrNotFound) {
		}
		if n, err := a.GetNoteListByUser(context.Background(), id); assert.NoError(t, err) && assert.Equal(t, 0, len(n.Ntps)) {
		}
		if n, err := a.GetNoteListByTag(context.Background(), id, id); assert.NoError(t, err) && assert.Equal(t, 0, len(n.Ntps)) {
		}
		assert.ErrorIs(t, a.UpdateTitle(context.Background(), id, "new_title"), domain2.ErrNotFound)
		assert.ErrorIs(t, a.UpdateUpdatedAt(context.Background(), id), domain2.ErrNotFound)
		assert.Error(t, a.delete(context.Background(), id))
		assert.ErrorIs(t, a.ShareNote(context.Background(), id, "neweditor", domain2.ReaderRole), domain2.ErrNotFound)
		assert.ErrorIs(t, a.DeleteRole(context.Background(), id, "newreader"), domain2.ErrNotFound)
	})
}

func TestBlockOrder(t *testing.T) {

	t.Run("test block order", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)

		id := uid.New()

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

		assert.NoError(t, a.Create(context.Background(), &domain2.Note{
			Id: id,
			//Tag:    &domain2.Tag{},
			Blocks: start,
		}))
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 4, 3))

		if note, err := a.Get(context.Background(), id, ""); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted1, note.Blocks)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 5, 1))
		if note, err := a.Get(context.Background(), id, ""); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted2, note.Blocks)
		}
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 0, 6))
		if note, err := a.Get(context.Background(), id, ""); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted3, note.Blocks)
		}
	})
}
func TestBlockOrder2(t *testing.T) {

	t.Run("test block order 2 el", func(t *testing.T) {
		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)
		id := uid.New()

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

		assert.NoError(t, a.Create(context.Background(), &domain2.Note{
			Id: id,
			//Tag:    &domain2.Tag{},
			Blocks: start,
		}))
		assert.NoError(t, a.ChangeBlockOrder(context.Background(), id, 1, 0))

		if note, err := a.Get(context.Background(), id, ""); assert.NoError(t, err) {
			log.Println(note.Blocks)
			assert.Equal(t, wanted1, note.Blocks)
		}
	})
}
