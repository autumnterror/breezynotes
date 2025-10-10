package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Auth(ctx context.Context, request *brzrpc.AuthRequest) (*brzrpc.Tokens, error) {
	return nil, nil
}

func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
