package psqltx

import (
	"context"
	"database/sql"
	"github.com/autumnterror/breezynotes/internal/auth/repository"
)

type RepoProvider struct {
	db *sql.DB
}

func NewRepoProvider(db *sql.DB) *RepoProvider {
	return &RepoProvider{db: db}
}

func (p *RepoProvider) Auth(ctx context.Context) repository.AuthRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}

func (p *RepoProvider) User(ctx context.Context) repository.UserRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}
func (p *RepoProvider) Health(ctx context.Context) repository.HealthRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}
