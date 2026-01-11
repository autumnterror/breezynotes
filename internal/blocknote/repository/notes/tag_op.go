package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"

	"github.com/autumnterror/breezynotes/pkg/utils/format"

	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// AddTagToNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix(). If tag exist func rewrite it
func (a *API) AddTagToNote(ctx context.Context, id string, tag *domain.Tag) error {
	const op = "notes.AddTagToNote"

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
					"tag":        tag,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil || res.MatchedCount == 0 {
		if res != nil && res.MatchedCount == 0 {
			return format.Error(op, domain.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}

// RemoveTagFromNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix().
func (a *API) RemoveTagFromNote(ctx context.Context, id string, tagId string) error {
	const op = "notes.RemoveTagFromNote"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res, err := a.
		db.
		UpdateOne(
			ctx,
			bson.M{
				"_id":     id,
				"tag._id": tagId,
			},
			bson.M{
				"$set": bson.M{
					"tag":        nil,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)

	if err != nil || res.MatchedCount == 0 {
		if res != nil && res.MatchedCount == 0 {
			return format.Error(op, domain.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}
