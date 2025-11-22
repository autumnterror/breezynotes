package grpc

import (
	"github.com/autumnterror/breezynotes/internal/redis/redis"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	brzrpc.UnimplementedRedisServiceServer
	rds redis.Repo
}

func Register(server *grpc.Server, rds redis.Repo) {
	brzrpc.RegisterRedisServiceServer(server, &ServerAPI{rds: rds})
}
