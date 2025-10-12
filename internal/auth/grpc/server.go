package grpc

import (
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/psql"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	brzrpc.UnimplementedAuthServiceServer
	UserAPI psql.AuthRepo
	JwtAPI  jwt.WithConfigRepo
}

func Register(server *grpc.Server, API psql.AuthRepo, JwtAPI jwt.WithConfigRepo) {
	brzrpc.RegisterAuthServiceServer(server, &ServerAPI{UserAPI: API, JwtAPI: JwtAPI})
}
