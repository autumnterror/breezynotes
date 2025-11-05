package tags

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
)

type API struct {
	*mongo.Client
}

func NewApi(c *mongo.Client) *API {
	return &API{c}
}
