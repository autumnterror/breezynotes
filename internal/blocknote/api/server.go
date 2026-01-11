package api

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/service"

	"google.golang.org/grpc"
	"time"
)

type ServerAPI struct {
	brzrpc.UnimplementedBlockNoteServiceServer
	tagAPI    *service.TagsService
	noteAPI   *service.NotesService
	blocksAPI *service.BlocksService
}

func Register(
	server *grpc.Server,
	tagAPI *service.TagsService,
	noteAPI *service.NotesService,
	blocksAPI *service.BlocksService,
) {
	brzrpc.RegisterBlockNoteServiceServer(server, &ServerAPI{
		tagAPI:    tagAPI,
		noteAPI:   noteAPI,
		blocksAPI: blocksAPI,
	})
}

const (
	waitTime = 3 * time.Second
)
