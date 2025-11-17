package grpc

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateTag(ctx context.Context, t *brzrpc.Tag) (*emptypb.Empty, error) {
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
		return nil, err
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagTitle(ctx context.Context, req *brzrpc.UpdateTagTitleRequest) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.UpdateTagTitle"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.UpdateTitle(ctx, req.GetId(), req.GetTitle()); err != nil {
			if errors.Is(err, mongo.ErrNotFiend) {
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
		return nil, err
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagColor(ctx context.Context, req *brzrpc.UpdateTagColorRequest) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.UpdateTagColor"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.UpdateColor(ctx, req.GetId(), req.GetColor()); err != nil {
			if errors.Is(err, mongo.ErrNotFiend) {
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
		return nil, err
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagEmoji(ctx context.Context, req *brzrpc.UpdateTagEmojiRequest) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.UpdateTagEmoji"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.UpdateEmoji(ctx, req.GetId(), req.GetEmoji()); err != nil {
			if errors.Is(err, mongo.ErrNotFiend) {
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
		return nil, err
	}

	return nil, nil
}
func (s *ServerAPI) DeleteTag(ctx context.Context, req *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "blocknote.grpc.DeleteTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.tagAPI.Delete(ctx, req.GetId()); err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetTagsByUser(ctx context.Context, req *brzrpc.Id) (*brzrpc.Tags, error) {
	const op = "blocknote.grpc.GetTagsByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		tgs, err := s.tagAPI.GetAllById(ctx, req.GetId())
		if err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: tgs, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Tags), nil
}

func (s *ServerAPI) GetTag(ctx context.Context, req *brzrpc.Id) (*brzrpc.Tag, error) {
	const op = "blocknote.grpc.GetTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		tg, err := s.tagAPI.Get(ctx, req.GetId())
		if err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: tg, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Tag), nil
}
