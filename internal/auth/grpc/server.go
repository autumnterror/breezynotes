package grpc

import (
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/psql"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/grpc"
	"time"
)

type ServerAPI struct {
	brzrpc.UnimplementedAuthServiceServer
	UserAPI psql.AuthRepo
	JwtAPI  jwt.WithConfigRepo
	cfg     *config.Config
}

func Register(server *grpc.Server, API psql.AuthRepo, JwtAPI jwt.WithConfigRepo, cfg *config.Config) {
	brzrpc.RegisterAuthServiceServer(server, &ServerAPI{UserAPI: API, JwtAPI: JwtAPI, cfg: cfg})
}

const (
	waitTime = 3 * time.Second
)
