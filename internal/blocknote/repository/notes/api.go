package notes

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/repository"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
)

type API struct {
	db          repository.NoSqlRepo
	trashDriver repository.NoSqlRepo
	blockAPI    blocks.Repo
}

func NewApi(c repository.NoSqlRepo, trashDriver repository.NoSqlRepo, blockAPI blocks.Repo) *API {
	return &API{db: c, trashDriver: trashDriver, blockAPI: blockAPI}
}
