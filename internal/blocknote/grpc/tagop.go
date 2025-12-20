package grpc

import (
	"context"
	"errors"
	brzrpc2 "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateTag(ctx context.Context, t *brzrpc2.Tag) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.CreateTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.Create(ctx, t); err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagTitle(ctx context.Context, req *brzrpc2.UpdateTagTitleRequest) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.UpdateTagTitle"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.UpdateTitle(ctx, req.GetId(), req.GetTitle()); err != nil {
			if errors.Is(err, mongo.ErrNotFound) {
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "not found")}
				return
			} else {
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
				return
			}
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagColor(ctx context.Context, req *brzrpc2.UpdateTagColorRequest) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.UpdateTagColor"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.UpdateColor(ctx, req.GetId(), req.GetColor()); err != nil {
			if errors.Is(err, mongo.ErrNotFound) {
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "not found")}
				return
			} else {
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
				return
			}
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagEmoji(ctx context.Context, req *brzrpc2.UpdateTagEmojiRequest) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.UpdateTagEmoji"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.UpdateEmoji(ctx, req.GetId(), req.GetEmoji()); err != nil {
			if errors.Is(err, mongo.ErrNotFound) {
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "not found")}
				return
			} else {
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
				return
			}
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) DeleteTag(ctx context.Context, req *brzrpc2.TagId) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.DeleteTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.Delete(ctx, req.GetTagId()); err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) GetTagsByUser(ctx context.Context, req *brzrpc2.UserId) (*brzrpc2.Tags, error) {
	const op = "blocknote.grpc.GetTagsByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		tgs, err := s.tagAPI.GetAllById(ctx, req.GetUserId())
		if err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: tgs, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc2.Tags), nil
}

func (s *ServerAPI) GetTag(ctx context.Context, req *brzrpc2.TagId) (*brzrpc2.Tag, error) {
	const op = "blocknote.grpc.GetTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		tg, err := s.tagAPI.Get(ctx, req.GetTagId())
		if err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: tg, Err: nil}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc2.Tag), nil
}
