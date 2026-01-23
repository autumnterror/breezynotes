package api

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/internal/auth/service"
	"google.golang.org/grpc"
	"time"
)

type ServerAPI struct {
	brzrpc.UnimplementedAuthServiceServer
	API *service.AuthService
	Cfg *config.Config
}

func Register(server *grpc.Server, s *service.AuthService, cfg *config.Config) {
	brzrpc.RegisterAuthServiceServer(server, &ServerAPI{API: s, Cfg: cfg})
}

const (
	waitTime = 3 * time.Second
)
