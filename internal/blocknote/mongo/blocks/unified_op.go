package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// createBlock with CreatedAt and UpdatedAt time.Now().UTC().Unix(). Don't create id
func (a *API) createBlock(ctx context.Context, b *views.BlockDb) error {
	const op = "notes.Create"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	b.CreatedAt, b.UpdatedAt = time.Now().UTC().Unix(), time.Now().UTC().Unix()

	if _, err := a.Blocks().InsertOne(ctx, b); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// Delete can return mongo.ErrNotFiend
func (a *API) Delete(ctx context.Context, id string) error {
	const op = "blocks.Delete"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.Blocks().DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil || res.DeletedCount == 0 {
		if res.DeletedCount == 0 {
			return format.Error(op, mongo.ErrNotFiend)
		}
		return format.Error(op, err)
	}

	return nil
}

func (a *API) Get(ctx context.Context, id string) (*brzrpc.Block, error) {
	const op = "blocks.Get"
	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res := a.Blocks().FindOne(ctx, bson.D{{"_id", id}})
	if res.Err() != nil {
		return nil, format.Error(op, res.Err())
	}

	var b views.BlockDb
	if err := res.Decode(&b); err != nil {
		return nil, format.Error(op, err)
	}

	return views.FromBlockDb(&b), nil
}

// UpdateData can return mongo.ErrNotFiend. Set updated_at to time.Now().UTC().Unix()
func (a *API) updateData(ctx context.Context, id string, data map[string]any) error {
	const op = "blocks.updateData"

	ctx, done := context.WithTimeout(ctx, views.WaitTime)
	defer done()

	res, err := a.
		Blocks().
		UpdateOne(
			ctx,
			bson.M{
				"_id": id,
			},
			bson.M{
				"$set": bson.M{
					"data":       data,
					"updated_at": time.Now().UTC().Unix(),
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

//func (a *API) ChangeBlockOrderWithPipeline(ctx context.Context, noteID string, oldOrder, newOrder int) error {
//	const op = "blocks.ChangeBlockOrder"
//
//	ctx, done := context.WithTimeout(ctx, mongo.WaitTime)
//	defer done()
//
//	// Собираем update-пайплайн (MongoDB 4.2+).
//	// Логика в терминах выражений Mongo:
//	//  from = clamp(oldOrder, 0, len(blocks)-1)
//	//  to   = clamp(newOrder, 0, len(blocks)-1)
//	//  val = blocks[from]
//	//  arrWithout = blocks[:from] + blocks[from+1:]
//	//  result = if from==to -> blocks
//	//           else        -> arrWithout[:to] + [val] + arrWithout[to:]
//	//
//	// Всё происходит на сервере, без двух этапов $pull/$push.
//	n := "$blocks"
//
//	pipeline := mongo2.Pipeline{
//		{{
//			Key: "$set",
//			Value: bson.D{
//				{
//					Key: "blocks",
//					Value: bson.D{
//						{
//							Key: "$let",
//							Value: bson.D{
//								{ // ВНЕШНИЙ $let: только базовые переменные
//									Key: "vars",
//									Value: bson.D{
//										{Key: "arr", Value: n},
//										{Key: "n", Value: bson.D{{Key: "$size", Value: n}}},
//										{Key: "from", Value: bson.D{
//											{Key: "$max", Value: bson.A{
//												0,
//												bson.D{{Key: "$min", Value: bson.A{
//													oldOrder,
//													bson.D{{Key: "$subtract", Value: bson.A{bson.D{{Key: "$size", Value: n}}, 1}}},
//												}}},
//											}},
//										}},
//										{Key: "to", Value: bson.D{
//											{Key: "$max", Value: bson.A{
//												0,
//												bson.D{{Key: "$min", Value: bson.A{
//													newOrder,
//													bson.D{{Key: "$subtract", Value: bson.A{bson.D{{Key: "$size", Value: n}}, 1}}},
//												}}},
//											}},
//										}},
//									},
//								},
//								{ // ВНУТРЕННИЙ $let: вычисляем val и arrWithout, тут уже можно $$from/$$arr
//									Key: "in",
//									Value: bson.D{
//										{
//											Key: "$let",
//											Value: bson.D{
//												{
//													Key: "vars",
//													Value: bson.D{
//														{Key: "val", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$$arr", "$$from"}}}},
//
//														// длина правой части при построении arrWithout: max(0, n - (from+1))
//														{Key: "rightLen", Value: bson.D{
//															{Key: "$max", Value: bson.A{
//																0,
//																bson.D{{Key: "$subtract", Value: bson.A{"$$n", bson.D{{Key: "$add", Value: bson.A{"$$from", 1}}}}}},
//															}},
//														}},
//
//														// arrWithout = left ++ (right если rightLen>0)
//														{Key: "arrWithout", Value: bson.D{
//															{Key: "$concatArrays", Value: bson.A{
//																// left = arr[:from]
//																bson.D{{Key: "$slice", Value: bson.A{"$$arr", 0, "$$from"}}},
//
//																// right = cond(rightLen>0 ? arr[from+1 : from+1+rightLen] : [])
//																bson.D{{Key: "$cond", Value: bson.A{
//																	bson.D{{Key: "$gt", Value: bson.A{"$$rightLen", 0}}},
//																	bson.D{{Key: "$slice", Value: bson.A{
//																		"$$arr",
//																		bson.D{{Key: "$add", Value: bson.A{"$$from", 1}}},
//																		"$$rightLen",
//																	}}},
//																	bson.A{}, // []
//																}}},
//															}},
//														}},
//													},
//												},
//												// --- ВНУТРЕННИЙ $let (in) ---
//												{
//													Key: "in",
//													Value: bson.D{
//														{Key: "$cond", Value: bson.A{
//															bson.D{{Key: "$eq", Value: bson.A{"$$from", "$$to"}}},
//															"$$arr",
//															bson.D{{Key: "$concatArrays", Value: bson.A{
//																// head = arrWithout[:to]
//																bson.D{{Key: "$slice", Value: bson.A{"$$arrWithout", 0, "$$to"}}},
//
//																// вставляем val
//																bson.A{"$$val"},
//
//																// tail = cond(tailLen>0 ? arrWithout[to : to+tailLen] : [])
//																bson.D{{Key: "$cond", Value: bson.A{
//																	bson.D{{Key: "$gt", Value: bson.A{
//																		bson.D{{Key: "$subtract", Value: bson.A{bson.D{{Key: "$size", Value: "$$arrWithout"}}, "$$to"}}},
//																		0,
//																	}}},
//																	bson.D{{Key: "$slice", Value: bson.A{
//																		"$$arrWithout",
//																		"$$to",
//																		bson.D{{Key: "$subtract", Value: bson.A{bson.D{{Key: "$size", Value: "$$arrWithout"}}, "$$to"}}},
//																	}}},
//																	bson.A{}, // []
//																}}},
//															}}},
//														}},
//													},
//												},
//											},
//										},
//									},
//								},
//							},
//						},
//					},
//				},
//				{Key: "updated_at", Value: time.Now().Unix()},
//			},
//		}},
//	}
//
//	filter := bson.M{"_id": noteID}
//
//	res, err := a.Notes().UpdateOne(ctx, filter, pipeline)
//	if err != nil {
//		return format.Error(op, err)
//	}
//	if res.MatchedCount == 0 {
//		return format.Error(op, errors.New("note not found"))
//	}
//
//	return nil
//}
