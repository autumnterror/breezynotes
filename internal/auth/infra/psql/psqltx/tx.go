package psqltx

import (
	"context"
	"database/sql"
	"github.com/autumnterror/utils_go/pkg/log"
)

type TxRunner struct {
	db *sql.DB
}

func NewTxRunner(db *sql.DB) *TxRunner {
	return &TxRunner{db: db}
}

type txKey struct{}

func (r *TxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	const op = "psql.RunInTx"

	if tx, ok := TxFromContext(ctx); ok {
		if err := fn(ctx); err != nil {
			errR := tx.Rollback()
			if errR != nil {
				log.Error(op, "rollback", err)
			}
			return err
		}
		return tx.Commit()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		errR := tx.Rollback()
		if errR != nil {
			log.Error(op, "rollback", err)
		}
		return err
	}
	return tx.Commit()
}

func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
