package grpc

import (
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	brzrpc.UnimplementedRedisServiceServer
}

func Register(server *grpc.Server) {
	brzrpc.RegisterRedisServiceServer(server, &ServerAPI{})
}
