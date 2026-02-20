package notes

import (
	"context"
	"errors"
	"time"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (a *API) ShareNote(ctx context.Context, noteId, userId, role string) error {
	const op = "notes.ShareNote"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	var update bson.M
	switch role {
	case domain.ReaderRole:
		update = bson.M{
			"$addToSet": bson.M{"readers": userId},
			"$pull":     bson.M{"editors": userId},
			"$set":      bson.M{"updated_at": time.Now().UTC().Unix()},
		}
	case domain.EditorRole:
		update = bson.M{
			"$addToSet": bson.M{"editors": userId},
			"$pull":     bson.M{"readers": userId},
			"$set":      bson.M{"updated_at": time.Now().UTC().Unix()},
		}
	default:
		return format.Error(op, errors.New("invalid role specified"))
	}

	res, err := a.noteAPI.UpdateOne(
		ctx,
		bson.M{"_id": noteId},
		update,
	)

	if err != nil {
		return format.Error(op, err)
	}
	if res.MatchedCount == 0 {
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

func (a *API) DeleteRole(ctx context.Context, noteId, userId string) error {
	const op = "notes.DeleteRole"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res, err := a.noteAPI.UpdateOne(
		ctx,
		bson.M{"_id": noteId},
		bson.M{
			"$pull": bson.M{"editors": userId, "readers": userId},
			"$set":  bson.M{"updated_at": time.Now().UTC().Unix()},
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

//func (a *API) ChangeUserRole(ctx context.Context, noteId, userId, newRole string) error {
//	const op = "notes.ChangeUserRole"
//	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
//	defer done()
//
//	var update bson.M
//	switch newRole {
//	case domain.ReaderRole:
//		update = bson.M{
//			"$pull":     bson.M{"editors": userId},
//			"$addToSet": bson.M{"readers": userId},
//			"$set":      bson.M{"updated_at": time.Now().UTC().Unix()},
//		}
//	case domain.EditorRole:
//		update = bson.M{
//			"$pull":     bson.M{"readers": userId},
//			"$addToSet": bson.M{"editors": userId},
//			"$set":      bson.M{"updated_at": time.Now().UTC().Unix()},
//		}
//	default:
//		return format.Error(op, errors.New("invalid new role specified"))
//	}
//
//	res, err := a.noteAPI.UpdateOne(
//		ctx,
//		bson.M{"_id": noteId},
//		update,
//	)
//
//	if err != nil {
//		return format.Error(op, err)
//	}
//	if res.MatchedCount == 0 {
//		return format.Error(op, domain.ErrNotFound)
//	}
//
//	return nil
//}
