package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
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

// ToTrash get note from a.Notes() insert on a.Trash() and remove from a.Notes()
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

// FromTrash remove note in a.Trash() and insert in a.Notes()
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

	if err := a.insert(ctx, mongo.FromNoteDb(&n)); err != nil {
		return format.Error(op, err)
	}

	_, err := a.Trash().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}
