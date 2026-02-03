package notes

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
)

type API struct {
	noteAPI  repository.NoSqlRepo
	trashAPI repository.NoSqlRepo
	blockAPI blocks.Repo
}

func NewApi(noteAPI repository.NoSqlRepo, trashAPI repository.NoSqlRepo, blockAPI blocks.Repo) *API {
	return &API{noteAPI: noteAPI, trashAPI: trashAPI, blockAPI: blockAPI}
}

type Repo interface {
	CleanTrash(ctx context.Context, uid string) error
	GetNotesFromTrash(ctx context.Context, uid string) (*domain.NoteParts, error)
	ToTrash(ctx context.Context, id string) error
	FromTrash(ctx context.Context, id string) error
	FindOnTrash(ctx context.Context, id string) (*domain.Note, error)

	Create(ctx context.Context, n *domain.Note) error
	Get(ctx context.Context, id string) (*domain.Note, error)
	GetNoteListByUser(ctx context.Context, id string) (*domain.NoteParts, error)
	GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain.NoteParts, error)

	AddTagToNote(ctx context.Context, id string, tag *domain.Tag) error
	RemoveTagFromNote(ctx context.Context, idNote string) error

	InsertBlock(ctx context.Context, id, blockId string, pos int) error
	DeleteBlock(ctx context.Context, id, blockId string) error
	ChangeBlockOrder(ctx context.Context, noteID string, oldOrder, newOrder int) error

	UpdateUpdatedAt(ctx context.Context, id string) error
	UpdateTitle(ctx context.Context, id string, nTitle string) error

	ShareNote(ctx context.Context, noteId, userId, role string) error
	ChangeUserRole(ctx context.Context, noteId, userId, newRole string) error

	Search(ctx context.Context, idUser, prompt string) (<-chan *domain.NotePart, error)
}
