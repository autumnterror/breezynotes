package notes

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Get return note by id can return mongo.ErrNotFound
func (a *API) Get(ctx context.Context, id string) (*brzrpc.Note, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Notes().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, format.Error(op, mongo.ErrNotFound)
		}
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

	cur, err := a.Notes().Find(ctx, bson.M{"author": id}, options.Find().SetSort(bson.M{"updated_at": -1}))

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

	cur, err := a.Notes().Find(ctx, bson.M{"tag._id": id, "author": idUser}, options.Find().SetSort(bson.M{"updated_at": -1}))
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &brzrpc.NoteParts{
		Items: []*brzrpc.NotePart{},
	}

	idx := 0
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
		idx++
	}

	return nts, nil
}

// GetAllByUser return note by id author
func (a *API) getAllByUser(ctx context.Context, id string) (*brzrpc.Notes, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	cur, err := a.Notes().Find(ctx, bson.M{"author": id}, options.Find().SetSort(bson.M{"updated_at": -1}))
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
