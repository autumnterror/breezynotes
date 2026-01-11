package service

import "context"

type TxRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
	Healthz(ctx context.Context) error
}
