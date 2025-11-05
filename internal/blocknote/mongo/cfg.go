package mongo

import (
	"errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	Db     = "blocknotedb"
	Tags   = "tags"
	Notes  = "notes"
	Blocks = "blocks"
)

var (
	ErrNotFiend = errors.New("not fiend")
)

func (c *Client) Tags() *mongo.Collection {
	return c.C.Database(Db).Collection(Tags)
}

func (c *Client) Notes() *mongo.Collection {
	return c.C.Database(Db).Collection(Notes)
}

func (c *Client) Blocks() *mongo.Collection {
	return c.C.Database(Db).Collection(Blocks)
}

func (c *Client) Trash() *mongo.Collection {
	return c.C.Database(Db).Collection(Blocks)
}
