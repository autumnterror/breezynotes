package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository"
)

type API struct {
	db repository.NoSqlRepo
}

func NewApi(db repository.NoSqlRepo) *API {
	return &API{db: db}
}

type Repo interface {
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*domain.Block, error)
	GetMany(ctx context.Context, ids []string) (*domain.Blocks, error)
	Create(ctx context.Context, _type, noteId string, data map[string]any) (string, error)
	OpBlock(ctx context.Context, id, opName string, data map[string]any) error
	GetAsFirst(ctx context.Context, id string) (string, error)
	ChangeType(ctx context.Context, id, newType string) error
}
