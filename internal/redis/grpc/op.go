package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "redis.grpc.Healthz"
	log.Info(op, "")

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
