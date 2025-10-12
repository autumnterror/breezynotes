package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) DeleteBlock(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.DeleteBlock"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) CreateBlock(context.Context, *brzrpc.Block) (*emptypb.Empty, error) {
	const op = "block.note.grpc.CreateBlock"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) OpBlock(context.Context, *brzrpc.OpBlockRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.OpBlock"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetBlock(context.Context, *brzrpc.Id) (*brzrpc.Block, error) {
	const op = "block.note.grpc.GetBlock"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetBlockAsFirst(context.Context, *brzrpc.Id) (*brzrpc.StringResponse, error) {
	const op = "block.note.grpc.GetBlockAsFirst"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) ChangeBlockOrder(context.Context, *brzrpc.ChangeBlockOrderRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeBlockOrder"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) ChangeTypeBlock(context.Context, *brzrpc.ChangeTypeBlockRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTypeBlock"
	log.Info(op, "")

	return nil, nil
}
