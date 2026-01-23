package notes

import (
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
