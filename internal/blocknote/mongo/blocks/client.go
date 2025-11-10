package blocks

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
)

type API struct {
	*mongo.Client
}

func NewApi(c *mongo.Client) *API {
	return &API{c}
}

var (
	ErrTypeNotDefined = errors.New("need to register type")
)

type Repo interface {
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*brzrpc.Block, error)
	UpdateData(ctx context.Context, id string, data map[string]any) error
	OpBlock(ctx context.Context, id, opName string, data map[string]any) error
	GetAsFirst(ctx context.Context, id string) (string, error)
	ChangeType(ctx context.Context, id, newType string) error
}
