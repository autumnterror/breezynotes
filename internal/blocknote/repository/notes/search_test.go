package notes

import (
	"context"
	"testing"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/textblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("search by title/content + access roles + short prompt behavior + duplicates allowed", func(t *testing.T) {
		block.RegisterBlock("text", &textblock.Driver{})

		m := mongo.MustConnect(config.Test())
		b := blocks.NewApi(m.Blocks())
		tgs := tags.NewApi(m.Tags(), m.NoteTags())
		a := NewApi(m.Notes(), m.Trash(), m.NoteTags(), tgs, b)

		uAuthor := uid.New()
		uEditor := uid.New()
		uReader := uid.New()
		uStranger := uid.New()

		nTitle := uid.New()
		nContent := uid.New()
		nAccess := uid.New()
		nNoAccess := uid.New()

		blkTitle1 := uid.New()
		blkTitle2 := uid.New()
		blkContent1 := uid.New()
		blkAccess1 := uid.New()
		blkNoAccess1 := uid.New()

		t.Cleanup(func() {
			_ = a.delete(context.Background(), nTitle)
			_ = a.delete(context.Background(), nContent)
			_ = a.delete(context.Background(), nAccess)
			_ = a.delete(context.Background(), nNoAccess)
			assert.NoError(t, m.Disconnect())
		})

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        nTitle,
			Title:     "Trip to Perm - TESTCASE",
			Author:    uAuthor,
			Editors:   []string{},
			Readers:   []string{},
			Blocks:    []string{},
			CreatedAt: 0,
			UpdatedAt: 10,
			IsBlog:    false,
			IsPublic:  false,
		}))

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        nContent,
			Title:     "Totally unrelated title",
			Author:    uAuthor,
			Editors:   []string{},
			Readers:   []string{},
			Blocks:    []string{},
			CreatedAt: 0,
			UpdatedAt: 9,
			IsBlog:    false,
			IsPublic:  false,
		}))

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        nAccess,
			Title:     "Chelyabinsk TESTCASE for readers",
			Author:    uAuthor,
			Editors:   []string{uEditor},
			Readers:   []string{uReader},
			Blocks:    []string{},
			CreatedAt: 0,
			UpdatedAt: 8,
			IsBlog:    false,
			IsPublic:  false,
		}))

		assert.NoError(t, a.Create(context.Background(), &domain.Note{
			Id:        nNoAccess,
			Title:     "Hidden TESTCASE note",
			Author:    uid.New(),
			Editors:   []string{},
			Readers:   []string{},
			Blocks:    []string{},
			CreatedAt: 0,
			UpdatedAt: 7,
			IsBlog:    false,
			IsPublic:  false,
		}))

		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:     blkTitle1,
			Type:   "text",
			NoteId: nTitle,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "This block contains testcase in content too"},
				},
			},
		}))
		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:     blkTitle2,
			Type:   "text",
			NoteId: nTitle,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "Other text"},
				},
			},
		}))

		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:     blkContent1,
			Type:   "text",
			NoteId: nContent,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "only CONTENT has TestCase keyword"},
				},
			},
		}))

		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:     blkAccess1,
			Type:   "text",
			NoteId: nAccess,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "access note content testcase"},
				},
			},
		}))

		assert.NoError(t, b.CreateBlock(context.Background(), &domain.Block{
			Id:     blkNoAccess1,
			Type:   "text",
			NoteId: nNoAccess,
			Data: map[string]any{
				"text": []any{
					map[string]any{"style": "default", "string": "hidden testcase content"},
				},
			},
		}))

		assert.NoError(t, a.InsertBlock(context.Background(), nTitle, blkTitle1, 0))
		assert.NoError(t, a.InsertBlock(context.Background(), nTitle, blkTitle2, 1))

		assert.NoError(t, a.InsertBlock(context.Background(), nContent, blkContent1, 0))
		assert.NoError(t, a.InsertBlock(context.Background(), nAccess, blkAccess1, 0))
		assert.NoError(t, a.InsertBlock(context.Background(), nNoAccess, blkNoAccess1, 0))

		collect := func(ch <-chan *domain.NotePart) []*domain.NotePart {
			var res []*domain.NotePart
			for it := range ch {
				res = append(res, it)
			}
			return res
		}

		{
			ch := a.Search(context.Background(), uAuthor, "TeStCaSe")
			got := collect(ch)
			for _, r := range got {
				log.Green("search(author) result", r.Id, r.Title, r.FirstBlock)
			}

			count := map[string]int{}
			for _, r := range got {
				count[r.Id]++
			}

			assert.GreaterOrEqual(t, count[nTitle], 1)
			assert.GreaterOrEqual(t, count[nTitle], 2)

			assert.GreaterOrEqual(t, count[nContent], 1)

			assert.Equal(t, 0, count[nNoAccess])
		}

		{
			ch := a.Search(context.Background(), uReader, "testcase")
			got := collect(ch)
			count := map[string]int{}
			for _, r := range got {
				count[r.Id]++
			}

			assert.GreaterOrEqual(t, count[nAccess], 1)
			assert.Equal(t, 0, count[nTitle])
			assert.Equal(t, 0, count[nContent])
			assert.Equal(t, 0, count[nNoAccess])
		}

		{
			ch := a.Search(context.Background(), uEditor, "testcase")
			got := collect(ch)
			count := map[string]int{}
			for _, r := range got {
				count[r.Id]++
			}

			assert.GreaterOrEqual(t, count[nAccess], 1)
			assert.Equal(t, 0, count[nTitle])
			assert.Equal(t, 0, count[nContent])
			assert.Equal(t, 0, count[nNoAccess])
		}

		{
			ch := a.Search(context.Background(), uStranger, "testcase")
			got := collect(ch)
			assert.Len(t, got, 0)
		}

		{
			ch := a.Search(context.Background(), uAuthor, "tr")
			got := collect(ch)
			count := map[string]int{}
			for _, r := range got {
				count[r.Id]++
			}

			assert.GreaterOrEqual(t, count[nTitle], 1)
			assert.Equal(t, 0, count[nContent])

			select {
			case <-time.After(10 * time.Millisecond):
			default:
			}
		}
	})
}
