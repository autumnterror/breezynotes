package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"

	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CleanTrash(ctx context.Context, req *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.CleanTrash"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.noteAPI.CleanTrash(ctx, req.GetUserId()); err != nil {
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
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) NoteToTrash(ctx context.Context, req *brzrpc.NoteId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ToTrash"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.noteAPI.ToTrash(ctx, req.GetNoteId()); err != nil {
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
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) NoteFromTrash(ctx context.Context, req *brzrpc.NoteId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.FromTrash"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.noteAPI.FromTrash(ctx, req.GetNoteId()); err != nil {
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
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) FindNoteInTrash(ctx context.Context, req *brzrpc.NoteId) (*brzrpc.Note, error) {
	const op = "block.note.grpc.FindNoteInTrash"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		n, err := s.noteAPI.FindFromTrash(ctx, req.GetNoteId())
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
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.Note), nil
}

func (s *ServerAPI) GetNotesFromTrash(ctx context.Context, req *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetNotesFromTrash"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		nts, err := s.noteAPI.GetNotesFromTrash(ctx, req.GetUserId())
		if err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{
				Res: nil,
				Err: status.Error(codes.Internal, err.Error()),
			}
			return
		}
		res <- views.ResRPC{
			Res: nts,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.NoteParts), nil
}
