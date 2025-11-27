package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// AddTagToNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix(). If tag exist func rewrite it
func (a *API) AddTagToNote(ctx context.Context, id string, tag *brzrpc.Tag) error {
	const op = "notes.AddTagToNote"

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
					"tag":        views.ToTagDb(tag),
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil || res.MatchedCount == 0 {
		if res != nil && res.MatchedCount == 0 {
			return format.Error(op, mongo.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}

// DeleteBlock can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) DeleteBlock(ctx context.Context, id, blockId string) error {
	const op = "notes.DeleteBlock"

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
				"$pull": bson.M{
					"blocks": blockId,
				},
				"$set": bson.M{
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

// RemoveTagFromNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix().
func (a *API) RemoveTagFromNote(ctx context.Context, id string, tagID string) error {
	const op = "notes.RemoveTagFromNote"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Notes().
		UpdateOne(
			ctx,
			bson.M{
				"_id":     id,
				"tag._id": tagID,
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
			return format.Error(op, mongo.ErrNotFound)
		}
		return format.Error(op, err)
	}

	return nil
}
