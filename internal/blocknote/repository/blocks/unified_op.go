package blocks

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Delete can return mongo.ErrNotFound
func (a *API) Delete(ctx context.Context, id string) error {
	const op = "blocks.delete"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res, err := a.db.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return format.Error(op, err)
	}
	if res.DeletedCount == 0 {
		return format.Error(op, domain.ErrNotFound)
	}
	return nil
}

func (a *API) Get(ctx context.Context, id string) (*domain.Block, error) {
	const op = "blocks.GetNote"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res := a.db.FindOne(ctx, bson.D{{"_id", id}})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, format.Error(op, domain.ErrNotFound)
		}
		return nil, format.Error(op, res.Err())
	}

	var b domain.Block
	if err := res.Decode(&b); err != nil {
		return nil, format.Error(op, err)
	}

	return &b, nil
}
