package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
)

type API struct {
	*mongo.Client
}

func NewApi(c *mongo.Client) *API {
	return &API{c}
}

type Repo interface {
	CleanTrash(ctx context.Context, uid string) error
	ToTrash(ctx context.Context, id string) error
	FromTrash(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*brzrpc.Note, error)
	Create(ctx context.Context, n *brzrpc.Note) error
	insert(ctx context.Context, n *brzrpc.Note) error
	Delete(ctx context.Context, id string) error
	UpdateUpdatedAt(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id string, nTitle string) error
}
