package api

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/redis/repository"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	brzrpc.UnimplementedRedisServiceServer
	rds repository.Repo
}

func Register(server *grpc.Server, rds repository.Repo) {
	brzrpc.RegisterRedisServiceServer(server, &ServerAPI{rds: rds})
}
