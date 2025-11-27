package notes

import (
	"context"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// InsertBlock can return mongo.ErrNotFound. Set updated_at to time.Now().UTC().Unix()
func (a *API) InsertBlock(ctx context.Context, id, blockId string, pos int) error {
	const op = "notes.InsertBlock"

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
				"$push": bson.M{
					"blocks": bson.M{
						"$each":     []string{blockId},
						"$position": pos,
					},
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

// ChangeBlockOrder вставляет блок в срез на новое место
func (a *API) ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error {
	const op = "blocks.ChangeBlockOrder"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	if oldOrder == newOrder {
		return nil
	}

	n, err := a.Get(ctx, noteID)
	if err != nil {
		return fmt.Errorf("%s: get note failed: %w", op, err)
	}
	if n == nil {
		return fmt.Errorf("%s: note not found", op)
	}

	blocks := n.Blocks
	l := len(blocks)
	if l == 0 {
		return nil
	}

	if oldOrder < 0 || oldOrder >= l {
		return format.Error(op, ErrBadRequest)
	}

	val := blocks[oldOrder]

	arrWithout := make([]string, 0, l-1)
	arrWithout = append(arrWithout, blocks[:oldOrder]...)
	arrWithout = append(arrWithout, blocks[oldOrder+1:]...)

	newBlocks := make([]string, 0, l)

	if newOrder < 0 {
		newOrder = 0
	} else if newOrder >= l {
		newOrder = l - 1
		newBlocks = append(newBlocks, arrWithout[:newOrder]...)
		newBlocks = append(newBlocks, val)
	} else {
		newBlocks = append(newBlocks, arrWithout[:newOrder]...)
		newBlocks = append(newBlocks, val)
		newBlocks = append(newBlocks, arrWithout[newOrder:]...)

		// Если ничего не изменилось — выходим
		changed := false
		if len(newBlocks) == len(blocks) {
			for i := range blocks {
				if blocks[i] != newBlocks[i] {
					changed = true
					break
				}
			}
		}
		if !changed {
			return nil
		}
	}
	if err := a.UpdateBlocks(ctx, noteID, newBlocks); err != nil {
		return format.Error(op, err)
	}
	return nil
}
