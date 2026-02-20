package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// UpdateData can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateData(ctx context.Context, id string, data map[string]any) error {
	const op = "blocks.UpdateData"
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

// UpdateType can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateType(ctx context.Context, id string, _type string) error {
	const op = "blocks.UpdateData"

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
					"type":       _type,
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

// UpdateUsed can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateUsed(ctx context.Context, id string, isUsedNew bool) error {
	const op = "blocks.UpdateData"

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

// CreateBlock with CreatedAt and UpdatedAt time.Now().UTC().Unix(). Don't create id
func (a *API) CreateBlock(ctx context.Context, b *domain.Block) error {
	const op = "blocks.Create"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	b.CreatedAt, b.UpdatedAt = time.Now().UTC().Unix(), time.Now().UTC().Unix()

	if _, err := a.db.InsertOne(ctx, b); err != nil {
		return format.Error(op, err)
	}
	return nil
}
