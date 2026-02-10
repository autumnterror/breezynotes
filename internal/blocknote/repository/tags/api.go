package tags

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/repository"
)

type API struct {
	db          repository.NoSqlRepo
	noteTagsAPI repository.NoSqlRepo
}

func NewApi(db repository.NoSqlRepo, noteTagsAPI repository.NoSqlRepo) *API {
	return &API{db: db, noteTagsAPI: noteTagsAPI}
}
