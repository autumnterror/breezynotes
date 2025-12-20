package grpc

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/tags"
	"google.golang.org/grpc"
	"time"
)

type ServerAPI struct {
	brzrpc.UnimplementedBlockNoteServiceServer
	tagAPI    tags.Repo
	noteAPI   notes.Repo
	blocksAPI blocks.Repo
}

func Register(server *grpc.Server, tagAPI tags.Repo, noteAPI notes.Repo, blocksAPI blocks.Repo) {
	brzrpc.RegisterBlockNoteServiceServer(server, &ServerAPI{
		tagAPI:    tagAPI,
		noteAPI:   noteAPI,
		blocksAPI: blocksAPI,
	})
}

const (
	waitTime = 3 * time.Second
)
