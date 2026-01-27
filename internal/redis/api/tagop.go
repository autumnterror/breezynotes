package api

import (
	"context"
	brzrpc2 "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetTagsByUser(ctx context.Context, req *brzrpc2.UserId) (*brzrpc2.Tags, error) {
	const op = "redis.grpc.GetTagsByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.rds.GetSessionTags(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc2.Tags{Items: res.([]*brzrpc2.Tag)}, nil
}

func (s *ServerAPI) SetTagsByUser(ctx context.Context, req *brzrpc2.TagsByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetTagsByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.rds.SetSessionTags(ctx, req.GetUserId(), req.GetItems())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) RmTagsByUser(ctx context.Context, req *brzrpc2.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmTagsByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.rds.SetSessionTags(ctx, req.GetUserId(), nil)
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
