package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

type Repo interface {
	CleanTrash(ctx context.Context, uid string) error
	GetNotesFromTrash(ctx context.Context, uid string) (*domain.NoteParts, error)
	ToTrash(ctx context.Context, id string) error
	FromTrash(ctx context.Context, id string) error
	FindOnTrash(ctx context.Context, id string) (*domain.Note, error)

	Create(ctx context.Context, n *domain.Note) error
	//insert(ctx context.Context, n *domain.Note) error
	Get(ctx context.Context, id string) (*domain.Note, error)
	GetNoteListByUser(ctx context.Context, id string) (*domain.NoteParts, error)
	GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain.NoteParts, error)

	AddTagToNote(ctx context.Context, id string, tag *domain.Tag) error
	RemoveTagFromNote(ctx context.Context, id string, tagId string) error

	InsertBlock(ctx context.Context, id, blockId string, pos int) error
	DeleteBlock(ctx context.Context, id, blockId string) error
	ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error

	UpdateUpdatedAt(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id string, nTitle string) error
	//updateBlocks(ctx context.Context, id string, blocks []string) error
}
