package repository

import (
	"context"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/redis/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Rdb *redis.Client
}

func New(cfg *config.Config) *Client {
	return &Client{Rdb: redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPasswd,
		DB:       cfg.RedisDb,
	})}
}

type Repo interface {
	GetSessionNoteList(ctx context.Context, id string) ([]*brzrpc.NotePart, error)
	SetSessionNoteList(ctx context.Context, id string, noteParts []*brzrpc.NotePart) error
	GetSessionNoteTrash(ctx context.Context, id string) ([]*brzrpc.NotePart, error)
	SetSessionNoteTrash(ctx context.Context, id string, noteTrash []*brzrpc.NotePart) error
	GetSessionNotes(ctx context.Context, id string) ([]*brzrpc.NoteWithBlocks, error)
	SetSessionNotes(ctx context.Context, id string, notes []*brzrpc.NoteWithBlocks) error
	GetSessionTags(ctx context.Context, id string) ([]*brzrpc.Tag, error)
	SetSessionTags(ctx context.Context, id string, tags []*brzrpc.Tag) error
	CreateSession(ctx context.Context, id string) error
	CheckSession(ctx context.Context, id string) error
}
