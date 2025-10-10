package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
