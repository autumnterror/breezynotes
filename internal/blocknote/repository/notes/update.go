package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"

	"github.com/autumnterror/utils_go/pkg/utils/format"

	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// UpdateUpdatedAt update updated time to time.Now().UTC().Unix() can return mongo.ErrNotFound
func (a *API) UpdateUpdatedAt(ctx context.Context, id string) error {
	const op = "notes.UpdateUpdatedAt"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.
		noteAPI.
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, domain2.ErrNotFound)
	}

	return nil
}

// UpdateTitle can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateTitle(ctx context.Context, id string, nTitle string) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.
		noteAPI.
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
		return format.Error(op, domain2.ErrNotFound)
	}

	return nil
}

// UpdateBlocks can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) updateBlocks(ctx context.Context, id string, blocks []string) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.
		noteAPI.
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
		return format.Error(op, domain2.ErrNotFound)
	}

	return nil
}

// UpdateBlog can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdateBlog(ctx context.Context, id string, isBlog bool) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.
		noteAPI.
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"is_blog":    isBlog,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, domain2.ErrNotFound)
	}

	return nil
}

// UpdatePublic can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) UpdatePublic(ctx context.Context, id string, isPublic bool) error {
	const op = "notes.UpdateTitle"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.
		noteAPI.
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"is_public":  isPublic,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, domain2.ErrNotFound)
	}

	return nil
}
