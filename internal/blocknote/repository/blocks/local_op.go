package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// UpdateData can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) updateData(ctx context.Context, id string, data map[string]any) error {
	const op = "blocks.updateData"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res, err := a.
		db.
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
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

// updateUsed can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) updateUsed(ctx context.Context, id string, isUsedNew bool) error {
	const op = "blocks.updateData"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
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
		db.
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
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

// createBlock with CreatedAt and UpdatedAt time.Now().UTC().Unix(). Don't create id
func (a *API) createBlock(ctx context.Context, b *domain.Block) error {
	const op = "notes.Create"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	b.CreatedAt, b.UpdatedAt = time.Now().UTC().Unix(), time.Now().UTC().Unix()

	if _, err := a.db.InsertOne(ctx, b); err != nil {
		return format.Error(op, err)
	}
	return nil
}
