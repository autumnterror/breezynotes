package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetBlocksOnNote(context.Context, *brzrpc.Id) (*brzrpc.Blocks, error) {
	const op = "redis.grpc.GetBlocksOnNote"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) SetBlocksOnNote(context.Context, *brzrpc.BlocksOnNote) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetBlocksOnNote"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) RmBlocksOnNote(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmBlocksOnNote"
	log.Info(op, "")

	return nil, nil
}
