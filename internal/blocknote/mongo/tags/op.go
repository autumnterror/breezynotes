package tags

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repo interface {
	GetAllById(ctx context.Context, id string) (*brzrpc.Tags, error)
	Create(ctx context.Context, t *brzrpc.Tag) error
	Delete(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id, nTitle string) error
	UpdateColor(ctx context.Context, id, nColor string) error
	UpdateEmoji(ctx context.Context, id, nEmoji string) error
}

func (a *API) GetAllById(ctx context.Context, id string) (*brzrpc.Tags, error) {
	const op = "tags.GetAllById"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	cur, err := a.Tags().Find(ctx, bson.M{"userId": id})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	tags := &brzrpc.Tags{
		Items: []*brzrpc.Tag{},
	}

	for cur.Next(ctx) {
		var t mongo.TagDb
		if err = cur.Decode(&t); err != nil {
			return tags, format.Error(op, err)
		}
		tags.Items = append(tags.Items, mongo.FromTagDb(&t))
	}

	return tags, nil
}

// Create tag. Don't create id
func (a *API) Create(ctx context.Context, t *brzrpc.Tag) error {
	const op = "tags.Create"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	if _, err := a.Tags().InsertOne(ctx, mongo.ToTagDb(t)); err != nil {
		return format.Error(op, err)
	}

	return nil
}

func (a *API) Delete(ctx context.Context, id string) error {
	const op = "tags.Delete"

	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
	defer done()

	if _, err := a.Tags().DeleteOne(ctx, bson.D{{"_id", id}}); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// UpdateTitle return mongo.ErrNotFiend
func (a *API) UpdateTitle(ctx context.Context, id, nTitle string) error {
	const op = "tags.UpdateTitle"

	res, err := a.
		Tags().
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

// UpdateColor return mongo.ErrNotFiend
func (a *API) UpdateColor(ctx context.Context, id, nColor string) error {
	const op = "tags.UpdateTitle"

	res, err := a.
		Tags().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"color": nColor,
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

// UpdateEmoji return mongo.ErrNotFiend
func (a *API) UpdateEmoji(ctx context.Context, id, nEmoji string) error {
	const op = "tags.UpdateTitle"

	res, err := a.
		Tags().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"emoji": nEmoji,
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
