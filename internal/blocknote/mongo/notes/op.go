package notes

import (
	"context"
	"fmt"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Get return note by id
func (a *API) Get(ctx context.Context, id string) (*brzrpc.Note, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Notes().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, format.Error(op, res.Err())
	}

	var n views.NoteDb
	err := res.Decode(&n)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return views.FromNoteDb(&n), nil
}

// GetNoteListByUser use func a.blockAPI.GetAsFirst
func (a *API) GetNoteListByUser(ctx context.Context, id string) (*brzrpc.NoteParts, error) {
	const op = "notes.GetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	cur, err := a.Notes().Find(ctx, bson.M{"author": id})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &brzrpc.NoteParts{
		Items: []*brzrpc.NotePart{},
	}

	for cur.Next(ctx) {
		var n views.NoteDb
		if err = cur.Decode(&n); err != nil {
			return nts, format.Error(op, err)
		}
		fb := ""
		if len(n.Blocks) > 0 {
			nfb, err := a.blockAPI.GetAsFirst(ctx, n.Blocks[0])
			if err == nil {
				fb = nfb
			} else {
				log.Warn(op, "get as first", err)
			}
		}
		nts.Items = append(nts.Items, &brzrpc.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        views.FromTagDb(n.Tag),
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
		})
	}

	return nts, nil
}

// GetNoteListByTag use func a.blockAPI.GetAsFirst
func (a *API) GetNoteListByTag(ctx context.Context, id, idUser string) (*brzrpc.NoteParts, error) {
	const op = "notes.GetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	cur, err := a.Notes().Find(ctx, bson.M{"tag._id": id, "author": idUser})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &brzrpc.NoteParts{
		Items: []*brzrpc.NotePart{},
	}

	for cur.Next(ctx) {
		var n views.NoteDb
		if err = cur.Decode(&n); err != nil {
			return nts, format.Error(op, err)
		}
		fb := ""
		if len(n.Blocks) > 0 {
			nfb, err := a.blockAPI.GetAsFirst(ctx, n.Blocks[0])
			if err == nil {
				fb = nfb
			}
		}
		nts.Items = append(nts.Items, &brzrpc.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        views.FromTagDb(n.Tag),
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
		})
	}

	return nts, nil
}

// GetAllByUser return note by id author
func (a *API) GetAllByUser(ctx context.Context, id string) (*brzrpc.Notes, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	cur, err := a.Notes().Find(ctx, bson.M{"author": id}, options.Find().SetSort(bson.M{"updatedAt": -1}))
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &brzrpc.Notes{
		Items: []*brzrpc.Note{},
	}

	for cur.Next(ctx) {
		var n views.NoteDb
		if err = cur.Decode(&n); err != nil {
			return nts, format.Error(op, err)
		}
		nts.Items = append(nts.Items, views.FromNoteDb(&n))
	}

	return nts, nil
}

// Create note with CreatedAt and UpdatedAt time.Now().UTC().Unix(). Don't create id
func (a *API) Create(ctx context.Context, n *brzrpc.Note) error {
	const op = "notes.Create"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	n.CreatedAt, n.UpdatedAt = time.Now().UTC().Unix(), time.Now().UTC().Unix()

	if _, err := a.Notes().InsertOne(ctx, views.ToNoteDb(n)); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// Insert created note in a.Notes() without any change
func (a *API) Insert(ctx context.Context, n *brzrpc.Note) error {
	const op = "notes.Insert"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if _, err := a.Notes().InsertOne(ctx, views.ToNoteDb(n)); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// Delete note WARNING DO NOT USE FROM gRPC NEED PASS TRASH CYCLE
func (a *API) Delete(ctx context.Context, id string) error {
	const op = "notes.Delete"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.Notes().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil || res.DeletedCount == 0 {
		if res.DeletedCount == 0 {
			return format.Error(op, mongo.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}

// UpdateUpdatedAt update updated time to time.Now().UTC().Unix() can return mongo.ErrNotFound
func (a *API) UpdateUpdatedAt(ctx context.Context, id string) error {
	const op = "notes.UpdateUpdatedAt"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Notes().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"updatedAt": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}

	return nil
}

// UpdateTitle can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateTitle(ctx context.Context, id string, nTitle string) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Notes().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"title":      nTitle,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}

	return nil
}

// UpdateBlocks can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateBlocks(ctx context.Context, id string, blocks []string) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Notes().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"blocks":     blocks,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}

	return nil
}

// InsertBlock can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) InsertBlock(ctx context.Context, id, blockId string, pos int) error {
	const op = "notes.InsertBlock"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Notes().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$push": bson.M{
					"blocks": bson.M{
						"$each":     []string{blockId},
						"$position": pos,
					},
				},
				"$set": bson.M{
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}

	return nil
}

// AddTagToNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix(). If tag exist func rewrite it
func (a *API) AddTagToNote(ctx context.Context, id string, tag *brzrpc.Tag) error {
	const op = "notes.AddTagToNote"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Notes().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
				//"tags._id": bson.M{"$ne": tag.Id},
			},
			bson.M{
				"$set": bson.M{
					"tag":        views.ToTagDb(tag),
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil || res.MatchedCount == 0 {
		if res != nil && res.MatchedCount == 0 {
			return format.Error(op, mongo.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}

// ChangeBlockOrder вставляет блок в срез на новое место
func (a *API) ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error {
	const op = "blocks.ChangeBlockOrder"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if oldOrder == newOrder {
		return nil
	}

	n, err := a.Get(ctx, noteID)
	if err != nil {
		return fmt.Errorf("%s: get note failed: %w", op, err)
	}
	if n == nil {
		return fmt.Errorf("%s: note not found", op)
	}

	blocks := n.Blocks
	l := len(blocks)
	if l == 0 {
		return nil
	}

	if oldOrder < 0 || oldOrder >= l {
		return format.Error(op, ErrBadRequest)
	}

	val := blocks[oldOrder]

	arrWithout := make([]string, 0, l-1)
	arrWithout = append(arrWithout, blocks[:oldOrder]...)
	arrWithout = append(arrWithout, blocks[oldOrder+1:]...)

	newBlocks := make([]string, 0, l)

	if newOrder < 0 {
		newOrder = 0
	} else if newOrder >= l {
		newOrder = l - 1
		newBlocks = append(newBlocks, arrWithout[:newOrder]...)
		newBlocks = append(newBlocks, val)
	} else {
		newBlocks = append(newBlocks, arrWithout[:newOrder]...)
		newBlocks = append(newBlocks, val)
		newBlocks = append(newBlocks, arrWithout[newOrder:]...)

		// Если ничего не изменилось — выходим
		changed := false
		if len(newBlocks) == len(blocks) {
			for i := range blocks {
				if blocks[i] != newBlocks[i] {
					changed = true
					break
				}
			}
		}
		if !changed {
			return nil
		}
	}
	if err := a.UpdateBlocks(ctx, noteID, newBlocks); err != nil {
		return format.Error(op, err)
	}
	return nil
	// Корректируем целевой индекс с учётом удаления:
	// если мы перемещаем элемент вправо (old < new) — после удаления
	// все индексы справа сдвинулись на -1, значит целевая позиция уменьшается на 1.
	// to := newOrder
	// if oldOrder < newOrder {
	// 	to = newOrder - 1
	// }

	// Зажимаем to в допустимый промежуток [0..len(arrWithout)]
	// (равенство правой границе означает «вставить в конец»)
	// if to < 0 {
	// 	to = 0
	// } else if to > len(arrWithout) {
	// 	to = len(arrWithout)
	// }

	// Формируем итоговый массив

}
