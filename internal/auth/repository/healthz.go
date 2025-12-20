package repository

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
)

type HealthRepo interface {
	Healthz(ctx context.Context) error
}

func (d Driver) Healthz(ctx context.Context) error {
	const op = "psql.Healthz"
	if _, err := d.Driver.ExecContext(ctx, "SELECT 1;"); err != nil {
		return format.Error(op, err)
	}
	return nil
}
