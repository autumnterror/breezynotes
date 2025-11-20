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
	FindFromTrash(ctx context.Context, id string) (*brzrpc.Note, error)

	Get(ctx context.Context, id string) (*brzrpc.Note, error)
	//GetAllByTag(ctx context.Context, id string) (*brzrpc.Notes, error)
	GetAllByUser(ctx context.Context, id string) (*brzrpc.Notes, error)
	GetNoteListByUser(ctx context.Context, id string) (*brzrpc.NoteParts, error)
	GetNoteListByTag(ctx context.Context, id, idUser string) (*brzrpc.NoteParts, error)
	Create(ctx context.Context, n *brzrpc.Note) error
	Insert(ctx context.Context, n *brzrpc.Note) error
	Delete(ctx context.Context, id string) error
	UpdateBlocks(ctx context.Context, id string, blocks []string) error
	UpdateUpdatedAt(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id string, nTitle string) error
	AddTagToNote(ctx context.Context, id string, tag *brzrpc.Tag) error
	InsertBlock(ctx context.Context, id, block string, pos int) error

	ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error

	Healthz(ctx context.Context) error
}

var (
	ErrBadRequest = errors.New("bad fields")
)

func (a *API) Healthz(ctx context.Context) error {
	var result map[string]interface{}

	err := a.Blocks().Database().RunCommand(ctx, map[string]interface{}{"ping": 1}).Decode(&result)
	if err != nil {
		return err
	}
	return nil
}
