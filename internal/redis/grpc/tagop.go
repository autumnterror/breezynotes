package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetTagsByUser(context.Context, *brzrpc.UserId) (*brzrpc.Tags, error) {
	return nil, nil
}

func (s *ServerAPI) SetTagsByUser(context.Context, *brzrpc.TagsByUser) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ServerAPI) RmTagsByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	return nil, nil
}
