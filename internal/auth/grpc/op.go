package grpc

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/psql"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Auth(ctx context.Context, r *brzrpc.AuthRequest) (*brzrpc.UserId, error) {
	const op = "auth.grpc.Auth"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		id, err := s.UserAPI.Authentication(ctx, r)
		if err != nil {
			switch {
			case errors.Is(err, psql.ErrNoUser):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, err.Error())}
			case errors.Is(err, psql.ErrPasswordIncorrect):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, err.Error())}
			case errors.Is(err, psql.ErrWrongInput):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.InvalidArgument, err.Error())}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			}
			return
		}

		res <- views.ResRPC{Res: id, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return &brzrpc.UserId{Id: res.(string)}, nil
}

func (s *ServerAPI) Healthz(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "auth.grpc.Healthz"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.Healthz(ctx); err != nil {
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, err.Error())}
			return
		}
		res <- views.ResRPC{Res: nil, Err: nil}
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opWithContext(ctx context.Context, op func(chan views.ResRPC)) (interface{}, error) {
	res := make(chan views.ResRPC, 1)

	go op(res)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "context deadline")
	case r := <-res:
		return r.Res, r.Err
	}
}
