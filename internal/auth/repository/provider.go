package repository

import "context"

type Provider interface {
	Auth(ctx context.Context) AuthRepo
	User(ctx context.Context) UserRepo
	Health(ctx context.Context) HealthRepo
}
