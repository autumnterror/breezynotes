package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// Get return note by id
func (a *API) Get(ctx context.Context, id string) (*brzrpc.Note, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	res := a.Notes().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, format.Error(op, res.Err())
	}

	var n mongo.NoteDb
	err := res.Decode(&n)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return mongo.FromNoteDb(&n), nil
}

// Create note with CreatedAt and UpdatedAt time.Now().UTC().Unix(). Don't create id
func (a *API) Create(ctx context.Context, n *brzrpc.Note) error {
	const op = "notes.Create"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	n.CreatedAt = time.Now().UTC().Unix()
	n.UpdatedAt = time.Now().UTC().Unix()

	if _, err := a.Notes().InsertOne(ctx, mongo.ToNoteDb(n)); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// insert created note in a.Notes() without any change
func (a *API) insert(ctx context.Context, n *brzrpc.Note) error {
	const op = "notes.insert"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	if _, err := a.Notes().InsertOne(ctx, mongo.ToNoteDb(n)); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// Delete note WARNING DO NOT USE FROM gRPC NEED PASS TRASH CYCLE
func (a *API) Delete(ctx context.Context, id string) error {
	const op = "notes.Delete"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	res, err := a.Notes().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil || res.DeletedCount == 0 {
		if res.DeletedCount == 0 {
			return format.Error(op, mongo.ErrNotFiend)
		}
		return format.Error(op, err)
	}

	return nil
}

// UpdateUpdatedAt update updated time to time.Now().UTC().Unix() can return mongo.ErrNotFiend
func (a *API) UpdateUpdatedAt(ctx context.Context, id string) error {
	const op = "notes.UpdateUpdatedAt"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
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
	if err != nil || res.MatchedCount == 0 {
		if res.MatchedCount == 0 {
			return format.Error(op, mongo.ErrNotFiend)
		}
		return format.Error(op, err)
	}

	return nil
}

// UpdateTitle can return mongo.ErrNotFiend
func (a *API) UpdateTitle(ctx context.Context, id string, nTitle string) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
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
					"title": nTitle,
				},
			},
		)
	if err != nil || res.MatchedCount == 0 {
		if res.MatchedCount == 0 {
			return format.Error(op, mongo.ErrNotFiend)
		}
		return format.Error(op, err)
	}

	return nil
}

//TODO GetNoteList WITH BLOCKZ
