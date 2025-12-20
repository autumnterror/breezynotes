package tags

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repo interface {
	Get(ctx context.Context, id string) (*brzrpc.Tag, error)
	GetAllById(ctx context.Context, id string) (*brzrpc.Tags, error)
	Create(ctx context.Context, t *brzrpc.Tag) error
	Delete(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id, nTitle string) error
	UpdateColor(ctx context.Context, id, nColor string) error
	UpdateEmoji(ctx context.Context, id, nEmoji string) error
}

func (a *API) Get(ctx context.Context, id string) (*brzrpc.Tag, error) {
	const op = "tags.GetNote"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Tags().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, format.Error(op, res.Err())
	}
	var t views.TagDb
	if err := res.Decode(&t); err != nil {
		return nil, format.Error(op, err)
	}

	return views.FromTagDb(&t), nil
}

func (a *API) GetAllById(ctx context.Context, id string) (*brzrpc.Tags, error) {
	const op = "tags.GetAllById"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	cur, err := a.Tags().Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	tags := &brzrpc.Tags{
		Items: []*brzrpc.Tag{},
	}

	for cur.Next(ctx) {
		var t views.TagDb
		if err = cur.Decode(&t); err != nil {
			return tags, format.Error(op, err)
		}
		tags.Items = append(tags.Items, views.FromTagDb(&t))
	}

	return tags, nil
}

// Create tag. Don't create id
func (a *API) Create(ctx context.Context, t *brzrpc.Tag) error {
	const op = "tags.Create"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if _, err := a.Tags().InsertOne(ctx, views.ToTagDb(t)); err != nil {
		return format.Error(op, err)
	}

	return nil
}

func (a *API) Delete(ctx context.Context, id string) error {
	const op = "tags.Delete"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if _, err := a.Tags().DeleteOne(ctx, bson.D{{"_id", id}}); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// UpdateTitle return mongo.ErrNotFound
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

	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}
	return nil
}

// UpdateColor return mongo.ErrNotFound
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

	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}
	return nil
}

// UpdateEmoji return mongo.ErrNotFound
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

	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, mongo.ErrNotFound)
	}
	return nil
}
