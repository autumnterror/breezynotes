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

func (s *ServerAPI) ChangeTitleNote(ctx context.Context, req *brzrpc.ChangeTitleNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTitleNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.noteAPI.UpdateTitle(ctx, req.GetId(), req.GetTitle()); err != nil {
			switch {
			case errors.Is(err, mongo.ErrNotFiend):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.NotFound, err.Error()),
				}
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
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetNote(ctx context.Context, req *brzrpc.Id) (*brzrpc.Note, error) {
	const op = "block.note.grpc.Get"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		n, err := s.noteAPI.Get(ctx, req.GetId())
		if err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{
				Res: nil,
				Err: status.Error(codes.Internal, err.Error()),
			}
			return
		}
		res <- views.ResRPC{
			Res: n,
			Err: nil,
		}
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Note), nil
}

func (s *ServerAPI) GetAllNotes(ctx context.Context, req *brzrpc.Id) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetAllNotes"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		n, err := s.noteAPI.GetAllByUser(ctx, req.GetId())
		if err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{
				Res: nil,
				Err: status.Error(codes.Internal, err.Error()),
			}
			return
		}
		res <- views.ResRPC{
			Res: n,
			Err: nil,
		}
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Notes), nil
}

func (s *ServerAPI) GetNotesByTag(ctx context.Context, req *brzrpc.Id) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetNotesByTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		n, err := s.noteAPI.GetAllByTag(ctx, req.GetId())
		if err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{
				Res: nil,
				Err: status.Error(codes.Internal, err.Error()),
			}
			return
		}
		res <- views.ResRPC{
			Res: n,
			Err: nil,
		}
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Notes), nil
}

func (s *ServerAPI) CreateNote(ctx context.Context, req *brzrpc.Note) (*emptypb.Empty, error) {
	const op = "block.note.grpc.Create"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.noteAPI.Create(ctx, req); err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{
				Res: nil,
				Err: status.Error(codes.Internal, err.Error()),
			}
			return
		}
		res <- views.ResRPC{
			Res: nil,
			Err: nil,
		}
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllBlocksInNote(ctx context.Context, req *brzrpc.Id) (*brzrpc.Blocks, error) {
	const op = "block.note.grpc.GetAllBlocksInNote"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) AddTagToNote(ctx context.Context, req *brzrpc.AddTagToNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.AddTagToNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		tag, err := s.tagAPI.Get(ctx, req.GetTagId())
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.InvalidArgument, err.Error()),
				}
			default:
				log.Warn(op, "", err)
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}

		if err := s.noteAPI.AddTagToNote(ctx, req.GetNoteId(), tag); err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{
				Res: nil,
				Err: status.Error(codes.Internal, err.Error()),
			}
			return
		}
		res <- views.ResRPC{
			Res: nil,
			Err: nil,
		}
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}
