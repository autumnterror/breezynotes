package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetTagsByUser(ctx context.Context, req *brzrpc.UserId) (*brzrpc.Tags, error) {
	const op = "redis.grpc.GetTagsByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		tgs, err := s.rds.GetSessionTags(ctx, req.GetId())
		if err != nil {
			switch {
			default:
				log.Warn(op, "", err)
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}
		res <- views.ResRPC{
			Res: tgs,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.Tags{Items: res.([]*brzrpc.Tag)}, nil
}

func (s *ServerAPI) SetTagsByUser(ctx context.Context, req *brzrpc.TagsByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetTagsByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		err := s.rds.SetSessionTags(ctx, req.GetUserId(), req.GetItems())
		if err != nil {
			switch {
			default:
				log.Warn(op, "", err)
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}

		res <- views.ResRPC{
			Res: nil,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) RmTagsByUser(ctx context.Context, req *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmTagsByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		err := s.rds.SetSessionTags(ctx, req.GetId(), nil)
		if err != nil {
			switch {
			default:
				log.Warn(op, "", err)
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}

		res <- views.ResRPC{
			Res: nil,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
