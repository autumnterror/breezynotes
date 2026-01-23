package api

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"

	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CleanTrash(ctx context.Context, req *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.CleanTrash"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.CleanTrash(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) NoteToTrash(ctx context.Context, req *brzrpc.UserNoteId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ToTrash"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.ToTrash(ctx, req.GetNoteId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) NoteFromTrash(ctx context.Context, req *brzrpc.UserNoteId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.FromTrash"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.FromTrash(ctx, req.GetNoteId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) FindNoteInTrash(ctx context.Context, req *brzrpc.UserNoteId) (*brzrpc.Note, error) {
	const op = "block.note.grpc.FindNoteInTrash"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.FindOnTrash(ctx, req.GetNoteId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.FromNoteDb(res.(*domain.Note)), nil
}

func (s *ServerAPI) GetNotesFromTrash(ctx context.Context, req *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetNotesFromTrash"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.GetNotesFromTrash(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.FromNotePartsDb(res.(*domain.NoteParts)), nil
}
