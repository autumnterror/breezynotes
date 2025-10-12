package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetTagsByUser(context.Context, *brzrpc.UserId) (*brzrpc.Tags, error) {
	const op = "redis.grpc.GetTagsByUser"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) SetTagsByUser(context.Context, *brzrpc.TagsByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetTagsByUser"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) RmTagsByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmTagsByUser"
	log.Info(op, "")

	return nil, nil
}
