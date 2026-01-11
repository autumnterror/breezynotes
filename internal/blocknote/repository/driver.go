package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Driver struct {
	Driver NoSqlRepo
}

func NewDriver(driver NoSqlRepo) *Driver {
	return &Driver{
		Driver: driver,
	}
}

type NoSqlRepo interface {
	InsertOne(ctx context.Context, document any,
		opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	InsertMany(
		ctx context.Context,
		documents any,
		opts ...options.Lister[options.InsertManyOptions],
	) (*mongo.InsertManyResult, error)
	DeleteOne(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.DeleteOneOptions],
	) (*mongo.DeleteResult, error)
	DeleteMany(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.DeleteManyOptions],
	) (*mongo.DeleteResult, error)
	UpdateOne(
		ctx context.Context,
		filter any,
		update any,
		opts ...options.Lister[options.UpdateOneOptions],
	) (*mongo.UpdateResult, error)
	UpdateMany(
		ctx context.Context,
		filter any,
		update any,
		opts ...options.Lister[options.UpdateManyOptions],
	) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter any,
		opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter any,
		opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	FindOneAndDelete(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.FindOneAndDeleteOptions]) *mongo.SingleResult
	FindOneAndReplace(
		ctx context.Context,
		filter any,
		replacement any,
		opts ...options.Lister[options.FindOneAndReplaceOptions],
	) *mongo.SingleResult
	FindOneAndUpdate(
		ctx context.Context,
		filter any,
		update any,
		opts ...options.Lister[options.FindOneAndUpdateOptions]) *mongo.SingleResult
}
