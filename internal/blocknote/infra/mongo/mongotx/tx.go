package mongotx

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

type TxRunner struct {
	db *mongo.Client
}

func NewTxRunner(db *mongo.Client) *TxRunner {
	return &TxRunner{db: db}
}

func (r *TxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	const op = "mongo.RunInTx"

	session, err := r.db.StartSession()
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer session.EndSession(ctx)

	txnOptions := options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.Majority())

	result, err := session.WithTransaction(ctx, fn, txnOptions)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return result, nil
}

func (r *TxRunner) Healthz(ctx context.Context) error {
	var result map[string]any

	err := r.db.Database(domain2.Db).RunCommand(ctx, map[string]interface{}{"ping": 1}).Decode(&result)
	if err != nil {
		return err
	}
	return nil
}
