package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetBlocksOnNote(context.Context, *brzrpc.Id) (*brzrpc.Blocks, error) {
	return nil, nil
}

func (s *ServerAPI) SetBlocksOnNote(context.Context, *brzrpc.BlocksOnNote) (*emptypb.Empty, error) {

	return nil, nil
}

func (s *ServerAPI) RmBlocksOnNote(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	return nil, nil
}
