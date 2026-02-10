package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"

	"github.com/autumnterror/utils_go/pkg/utils/format"

	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// AddTagToNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix(). If tag exist func rewrite it
func (a *API) AddTagToNote(ctx context.Context, id string, tag *domain2.Tag) error {
	const op = "notes.AddTagToNote"

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
					//"tag":        tag,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)
	if err != nil || res.MatchedCount == 0 {
		if res != nil && res.MatchedCount == 0 {
			return format.Error(op, domain2.ErrNotFound)
		}
		return format.Error(op, err)
	}

	_, err = a.
		noteTagsAPI.
		InsertOne(
			ctx,
			bson.M{
				"note_id": id,
				"tag":     tag,
			},
		)
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}

// RemoveTagFromNote can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix().
func (a *API) RemoveTagFromNote(ctx context.Context, idNote string, idUser string) error {
	const op = "notes.RemoveTagFromNote"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res, err := a.
		noteAPI.
		UpdateOne(
			ctx,
			bson.M{
				"_id": idNote,
			},
			bson.M{
				"$set": bson.M{
					//"tag":        nil,
					"updated_at": time.Now().UTC().Unix(),
				},
			},
		)

	if err != nil || res.MatchedCount == 0 {
		if res != nil && res.MatchedCount == 0 {
			return format.Error(op, domain2.ErrNotFound)
		}
		return format.Error(op, err)
	}

	resDelete, err := a.
		noteTagsAPI.
		DeleteOne(
			ctx,
			bson.M{
				"note_id":     idNote,
				"tag.user_id": idUser,
			},
		)
	if err != nil {
		return format.Error(op, err)
	}
	if resDelete.DeletedCount == 0 {
		return format.Error(op, domain2.ErrNotFound)
	}

	return nil
}
