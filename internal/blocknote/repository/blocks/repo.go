package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

type Repo interface {
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*domain.Block, error)
	Create(ctx context.Context, _type, noteId string, data map[string]any) (string, error)
	OpBlock(ctx context.Context, id, opName string, data map[string]any) error
	GetAsFirst(ctx context.Context, id string) (string, error)
	ChangeType(ctx context.Context, id, newType string) error
}
