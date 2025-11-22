package grpc

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/redis/config"
	"github.com/autumnterror/breezynotes/internal/redis/redis"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type App struct {
	gRPCServer *grpc.Server
	cfg        *config.Config
}

func New(cfg *config.Config, rds redis.Repo) *App {
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 0,
		}),
	)
	Register(s, rds)

	return &App{
		gRPCServer: s,
		cfg:        cfg,
	}
}

const (
	waitTime = time.Second
)

// MustRun running gRPC server and panic if error
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc.redis.App"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return format.Error(op, err)
	}
	log.Success(op, "grpc server is running "+fmt.Sprint(a.cfg.Port))

	if err := a.gRPCServer.Serve(l); err != nil {
		return format.Error(op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpc.redis.Stop"
	a.gRPCServer.GracefulStop()
	log.Success(op, "grpc server is stop "+fmt.Sprint(a.cfg.Port))
}
