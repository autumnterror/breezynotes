package notes

import (
	"context"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
