package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type Client struct {
	C   *mongo.Client
	cfg *config.Config
}

func MustConnect(cfg *config.Config) *Client {
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	c, err := newConnect(cfg, ctx)
	if err != nil {
		log.Panic(err)
	}
	return c
}

func newConnect(
	cfg *config.Config,
	ctx context.Context,
) (*Client, error) {
	const op = "app.newConnect"

	c, err := mongo.Connect(options.Client().ApplyURI(cfg.Uri))
	if err != nil {
		return nil, format.Error(op, fmt.Errorf("error in connection to app: %w", err))
	}

	if err = c.Ping(ctx, nil); err != nil {
		return nil, format.Error(op, fmt.Errorf("PING: connection not established: %w", err))
	}

	log.Green("Connection to MongoDB is established successfully")

	return &Client{
		C:   c,
		cfg: cfg,
	}, nil
}

func (c *Client) Disconnect() (err error) {
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	
	if err = c.C.Disconnect(ctx); err != nil {
		return fmt.Errorf("error when terminating connection to app: %w", err)
	}
	if err = c.C.Ping(ctx, nil); err == nil {
		return errors.New("PING: terminating connection to MongoDB is failed")
	}
	log.Green("Connection to MongoDB is terminated successfully")

	return nil
}
