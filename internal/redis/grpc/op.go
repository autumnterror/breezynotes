package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "redis.grpc.Healthz"
	log.Info(op, "")

	return nil, nil
}
