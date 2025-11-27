package notes

import (
	"context"

	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CleanTrash rm all in a.Trash() by uid and field author
func (a *API) CleanTrash(ctx context.Context, uid string) error {
	const op = "notes.CleanTrash"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	_, err := a.Trash().DeleteMany(ctx, bson.M{"author": uid})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}

// GetNotesFromTrash by author id
func (a *API) GetNotesFromTrash(ctx context.Context, uid string) (*brzrpc.NoteParts, error) {
	const op = "notes.CleanTrash"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	cur, err := a.Trash().Find(ctx, bson.M{"author": uid})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	pts := &brzrpc.NoteParts{
		Items: []*brzrpc.NotePart{},
	}

	for cur.Next(ctx) {
		var n views.NoteDb
		if err = cur.Decode(&n); err != nil {
			return pts, format.Error(op, err)
		}
		fb := ""
		if len(n.Blocks) > 0 {
			nfb, err := a.blockAPI.GetAsFirst(ctx, n.Blocks[0])
			if err == nil {
				fb = nfb
			}
		}
		np := brzrpc.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        views.FromTagDb(n.Tag),
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
		}
		pts.Items = append(pts.Items, &np)
	}

	return pts, nil
}

// ToTrash get note from a.Notes() Insert on a.Trash() and remove from a.Notes()
func (a *API) ToTrash(ctx context.Context, id string) error {
	const op = "notes.ToTrash"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	n, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if _, err := a.Trash().InsertOne(ctx, views.ToNoteDb(n)); err != nil {
		return format.Error(op, err)
	}

	if err := a.Delete(ctx, n.Id); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// FromTrash remove note in a.Trash() and Insert in a.Notes()
func (a *API) FromTrash(ctx context.Context, id string) error {
	const op = "notes.FromTrash"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Trash().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return format.Error(op, res.Err())
	}

	var n views.NoteDb
	if err := res.Decode(&n); err != nil {
		return format.Error(op, err)
	}

	if err := a.Insert(ctx, views.FromNoteDb(&n)); err != nil {
		return format.Error(op, err)
	}

	_, err := a.Trash().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}

// FindFromTrash return note by id from Trash
func (a *API) FindFromTrash(ctx context.Context, id string) (*brzrpc.Note, error) {
	const op = "notes.FindFromTrash"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Trash().FindOne(ctx, bson.M{"_id": id})
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
