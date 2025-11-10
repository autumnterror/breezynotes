package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CleanTrash rm all in a.Trash() by uid and field author
func (a *API) CleanTrash(ctx context.Context, uid string) error {
	const op = "notes.CleanTrash"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
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

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
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
		var n mongo.NoteDb
		if err = cur.Decode(&n); err != nil {
			return pts, format.Error(op, err)
		}
		np := brzrpc.NotePart{
			Id:    n.Id,
			Title: n.Title,
			Tag:   mongo.FromTagDb(n.Tag),
			//TODO first block from blocks
			FirstBlock: "",
			UpdatedAt:  n.UpdatedAt,
		}
		pts.Items = append(pts.Items, &np)
	}

	return pts, nil
}

// ToTrash get note from a.Notes() Insert on a.Trash() and remove from a.Notes()
func (a *API) ToTrash(ctx context.Context, id string) error {
	const op = "notes.ToTrash"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	n, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if _, err := a.Trash().InsertOne(ctx, mongo.ToNoteDb(n)); err != nil {
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

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	res := a.Trash().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return format.Error(op, res.Err())
	}

	var n mongo.NoteDb
	if err := res.Decode(&n); err != nil {
		return format.Error(op, err)
	}

	if err := a.Insert(ctx, mongo.FromNoteDb(&n)); err != nil {
		return format.Error(op, err)
	}

	_, err := a.Trash().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}
