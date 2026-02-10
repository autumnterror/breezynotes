package mongo

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (c *Client) Tags() *mongo.Collection {
	return c.C.Database(domain2.Db).Collection(domain2.TagColl)
}

func (c *Client) Notes() *mongo.Collection {
	return c.C.Database(domain2.Db).Collection(domain2.NoteColl)
}

func (c *Client) Blocks() *mongo.Collection {
	return c.C.Database(domain2.Db).Collection(domain2.BlockColl)
}

func (c *Client) Trash() *mongo.Collection {
	return c.C.Database(domain2.Db).Collection(domain2.TrashColl)
}
func (c *Client) NoteTags() *mongo.Collection {
	return c.C.Database(domain2.Db).Collection(domain2.NoteTagsColl)
}
