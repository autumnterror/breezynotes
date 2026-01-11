package mongo

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (c *Client) Tags() *mongo.Collection {
	return c.C.Database(domain.Db).Collection(domain.TagColl)
}

func (c *Client) Notes() *mongo.Collection {
	return c.C.Database(domain.Db).Collection(domain.NoteColl)
}

func (c *Client) Blocks() *mongo.Collection {
	return c.C.Database(domain.Db).Collection(domain.BlockColl)
}

func (c *Client) Trash() *mongo.Collection {
	return c.C.Database(domain.Db).Collection(domain.TrashColl)
}
