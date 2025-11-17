package notes

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/blocks"
	"github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
)

type API struct {
	*mongo.Client
	blockAPI blocks.Repo
}

func NewApi(c *mongo.Client, blockAPI blocks.Repo) *API {
	return &API{c, blockAPI}
}

type Repo interface {
	GetNotesFromTrash(ctx context.Context, uid string) (*brzrpc.NoteParts, error)
	CleanTrash(ctx context.Context, uid string) error
	ToTrash(ctx context.Context, id string) error
	FromTrash(ctx context.Context, id string) error

	Get(ctx context.Context, id string) (*brzrpc.Note, error)
	GetAllByTag(ctx context.Context, id string) (*brzrpc.Notes, error)
	GetAllByUser(ctx context.Context, id string) (*brzrpc.Notes, error)
	Create(ctx context.Context, n *brzrpc.Note) error
	Insert(ctx context.Context, n *brzrpc.Note) error
	Delete(ctx context.Context, id string) error
	UpdateBlocks(ctx context.Context, id string, blocks []string) error
	UpdateUpdatedAt(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id string, nTitle string) error
	AddTagToNote(ctx context.Context, id string, tag *brzrpc.Tag) error

	ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error
}

var (
	ErrBadRequest = errors.New("bad fields")
)
