package grpc

import (
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	brzrpc.UnimplementedBlockNoteServiceServer
}

func Register(server *grpc.Server) {
	brzrpc.RegisterBlockNoteServiceServer(server, &ServerAPI{})
}
