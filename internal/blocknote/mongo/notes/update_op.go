package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

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
