package blocks

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
)

type API struct {
	*mongo.Client
}

func NewApi(c *mongo.Client) *API {
	return &API{c}
}

var (
	ErrTypeNotDefined = errors.New("need to register type")
	ErrAlreadyUsed    = errors.New("block already in use")
)

type Repo interface {
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*brzrpc.brzrpc, error)

	OpBlock(ctx context.Context, id, opName string, data map[string]any) error
	GetAsFirst(ctx context.Context, id string) (string, error)
	ChangeType(ctx context.Context, id, newType string) error
	Create(ctx context.Context, _type, idNote string, data map[string]any) (string, error)
	UpdateUsed(ctx context.Context, id string, isUsedNew bool) error
	//Render(ctx context.Context, id, _type string) (*brzrpc.Block, error)
}
