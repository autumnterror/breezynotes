package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// createBlock with CreatedAt and UpdatedAt time.Now().UTC().Unix(). Don't create id
func (a *API) createBlock(ctx context.Context, b *views.BlockDb) error {
	const op = "notes.Create"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	b.CreatedAt, b.UpdatedAt = time.Now().UTC().Unix(), time.Now().UTC().Unix()

	if _, err := a.Blocks().InsertOne(ctx, b); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// Delete can return mongo.ErrNotFound
func (a *API) Delete(ctx context.Context, id string) error {
	const op = "blocks.Delete"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.Blocks().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil || res.DeletedCount == 0 {
		if res.DeletedCount == 0 {
			return format.Error(op, mongo.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}

func (a *API) Get(ctx context.Context, id string) (*brzrpc.Block, error) {
	const op = "blocks.Get"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Blocks().FindOne(ctx, bson.D{{"_id", id}})
	if res.Err() != nil {
		return nil, format.Error(op, res.Err())
	}

	var b views.BlockDb
	if err := res.Decode(&b); err != nil {
		return nil, format.Error(op, err)
	}

	return views.FromBlockDb(&b), nil
}

// UpdateData can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) updateData(ctx context.Context, id string, data map[string]any) error {
	const op = "blocks.updateData"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Blocks().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"data":       data,
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

// UpdateUsed can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateUsed(ctx context.Context, id string, isUsedNew bool) error {
	const op = "blocks.updateData"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	var filter bson.M
	if isUsedNew {
		filter = bson.M{
			"_id":     id,
			"is_used": false,
		}
	} else {
		filter = bson.M{
			"_id":     id,
			"is_used": true,
		}
	}

	res, err := a.
		Blocks().
		UpdateOne(
			ctx,
			filter,
			bson.M{
				"$set": bson.M{
					"is_used":    isUsedNew,
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
