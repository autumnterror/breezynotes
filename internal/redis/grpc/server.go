package grpc

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/redis/redis"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	brzrpc.UnimplementedRedisServiceServer
	rds redis.Repo
}

func Register(server *grpc.Server, rds redis.Repo) {
	brzrpc.RegisterRedisServiceServer(server, &ServerAPI{rds: rds})
}
