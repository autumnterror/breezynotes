package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"

	"github.com/autumnterror/utils_go/pkg/utils/format"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Create note. Don't create id
func (a *API) Create(ctx context.Context, n *domain2.Note) error {
	const op = "notes.Create"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	if _, err := a.noteAPI.InsertOne(ctx, n); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// Insert created note in a.repo() without any change
func (a *API) insert(ctx context.Context, n *domain2.Note) error {
	const op = "notes.insert"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	if _, err := a.noteAPI.InsertOne(ctx, n); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// delete note WARNING DO NOT USE FROM gRPC NEED PASS TRASH CYCLE
func (a *API) delete(ctx context.Context, id string) error {
	const op = "notes.delete"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.noteAPI.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil || res.DeletedCount == 0 {
		if res.DeletedCount == 0 {
			return format.Error(op, domain2.ErrNotFound)
		}
		return format.Error(op, err)
	}

	res, err = a.noteTagsAPI.DeleteOne(ctx, bson.D{{"note_id", id}})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}

// deleteMany note WARNING DO NOT USE FROM gRPC NEED PASS TRASH CYCLE
func (a *API) deleteMany(ctx context.Context, ids []string) error {
	const op = "notes.delete"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()
	if len(ids) == 0 {
		return nil
	}
	filter := bson.D{{"_id", bson.D{{"$in", ids}}}}
	_, err := a.noteAPI.DeleteMany(ctx, filter)

	if err != nil {
		return format.Error(op, err)
	}

	return nil
}
