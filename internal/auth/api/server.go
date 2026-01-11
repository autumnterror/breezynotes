package api

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/auth/service"
	"google.golang.org/grpc"
	"time"
)

type ServerAPI struct {
	brzrpc.UnimplementedAuthServiceServer
	API *service.AuthService
}

func Register(server *grpc.Server, s *service.AuthService) {
	brzrpc.RegisterAuthServiceServer(server, &ServerAPI{API: s})
}

const (
	waitTime = 3 * time.Second
)
