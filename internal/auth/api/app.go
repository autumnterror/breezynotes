package api

import (
	"context"
	"fmt"
	"net"

	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/internal/auth/service"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type App struct {
	gRPCServer *grpc.Server
	cfg        *config.Config
}

func New(cfg *config.Config, as *service.AuthService) *App {
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 0,
		}),
	)
	Register(s, as, cfg)

	return &App{
		gRPCServer: s,
		cfg:        cfg,
	}
}

// MustRun running gRPC server and panic if error
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc.auth.App"

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
	const op = "grpc.auth.Stop"
	a.gRPCServer.GracefulStop()
	log.Success(op, "grpc server is stop "+fmt.Sprint(a.cfg.Port))
}

func (s *ServerAPI) Healthz(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "grpc.Healthz"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.Health(ctx)
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &emptypb.Empty{}, nil
}
