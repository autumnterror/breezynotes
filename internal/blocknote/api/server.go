package api

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/service"

	"google.golang.org/grpc"
	"time"
)

type ServerAPI struct {
	brzrpc.UnimplementedBlockNoteServiceServer
	service *service.BN
}

func Register(
	server *grpc.Server,
	noteAPI *service.BN,
) {
	brzrpc.RegisterBlockNoteServiceServer(server, &ServerAPI{
		service: noteAPI,
	})
}

const (
	waitTime = 3 * time.Second
)
