package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) DeleteBlock(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) CreateBlock(context.Context, *brzrpc.Block) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) OpBlock(context.Context, *brzrpc.OpBlockRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) GetBlock(context.Context, *brzrpc.Id) (*brzrpc.Block, error) {
	return nil, nil
}
func (s *ServerAPI) GetBlockAsFirst(context.Context, *brzrpc.Id) (*brzrpc.StringResponse, error) {
	return nil, nil
}
func (s *ServerAPI) ChangeBlockOrder(context.Context, *brzrpc.ChangeBlockOrderRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) ChangeTypeBlock(context.Context, *brzrpc.ChangeTypeBlockRequest) (*emptypb.Empty, error) {
	return nil, nil
}
