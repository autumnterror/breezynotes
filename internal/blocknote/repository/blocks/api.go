package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository"
)

type API struct {
	db repository.NoSqlRepo
}

func NewApi(db repository.NoSqlRepo) *API {
	return &API{db: db}
}

type Repo interface {
	UpdateData(ctx context.Context, id string, data map[string]any) error
	UpdateType(ctx context.Context, id string, _type string) error
	UpdateUsed(ctx context.Context, id string, isUsedNew bool) error
	CreateBlock(ctx context.Context, b *domain2.Block) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
	Get(ctx context.Context, id string) (*domain2.Block, error)
	GetMany(ctx context.Context, ids []string) (*domain2.Blocks, error)
	GetAsFirst(ctx context.Context, id string) (string, error)
}
