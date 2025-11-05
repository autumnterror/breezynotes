package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
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
		if err := s.noteAPI.CleanTrash(ctx, req.GetId()); err != nil {
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

func (s *ServerAPI) NoteToTrash(ctx context.Context, req *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ToTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) NoteFromTrash(ctx context.Context, req *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.FromTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) FindNoteInTrash(ctx context.Context, req *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.FindNoteInTrash"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) GetNotesFromTrash(ctx context.Context, req *brzrpc.Id) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetNotesFromTrash"
	log.Info(op, "")

	return nil, nil
}
