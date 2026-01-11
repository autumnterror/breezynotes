package tags

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo interface {
	Get(ctx context.Context, id string) (*domain.Tag, error)
	GetAllById(ctx context.Context, id string) (*domain.Tags, error)
	Create(ctx context.Context, t *domain.Tag) error
	Delete(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id, nTitle string) error
	UpdateColor(ctx context.Context, id, nColor string) error
	UpdateEmoji(ctx context.Context, id, nEmoji string) error
}

func (a *API) Get(ctx context.Context, id string) (*domain.Tag, error) {
	const op = "tags.Get"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res := a.db.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, format.Error(op, domain.ErrNotFound)
		}
		return nil, format.Error(op, err)
	}
	var t domain.Tag
	if err := res.Decode(&t); err != nil {
		return nil, format.Error(op, err)
	}

	return &t, nil
}

func (a *API) GetAllById(ctx context.Context, id string) (*domain.Tags, error) {
	const op = "tags.GetAllById"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	cur, err := a.db.Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	tags := &domain.Tags{
		Tgs: []*domain.Tag{},
	}

	for cur.Next(ctx) {
		var t domain.Tag
		if err = cur.Decode(&t); err != nil {
			return nil, format.Error(op, err)
		}
		tags.Tgs = append(tags.Tgs, &t)
	}

	return tags, nil
}

// Create tag. Don't create id
func (a *API) Create(ctx context.Context, t *domain.Tag) error {
	const op = "tags.Create"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	if _, err := a.db.InsertOne(ctx, t); err != nil {
		return format.Error(op, err)
	}

	return nil
}

func (a *API) Delete(ctx context.Context, id string) error {
	const op = "tags.delete"

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

// UpdateTitle return mongo.ErrNotFound
func (a *API) UpdateTitle(ctx context.Context, id, nTitle string) error {
	const op = "tags.UpdateTitle"

	res, err := a.
		db.
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
		return format.Error(op, domain.ErrNotFound)
	}
	return nil
}

// UpdateColor return mongo.ErrNotFound
func (a *API) UpdateColor(ctx context.Context, id, nColor string) error {
	const op = "tags.UpdateColor"

	res, err := a.
		db.
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
		return format.Error(op, domain.ErrNotFound)
	}
	return nil
}

// UpdateEmoji return mongo.ErrNotFound
func (a *API) UpdateEmoji(ctx context.Context, id, nEmoji string) error {
	const op = "tags.UpdateEmoji"

	res, err := a.
		db.
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
		return format.Error(op, domain.ErrNotFound)
	}
	return nil
}
