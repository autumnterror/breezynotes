package api

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/format"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ChangeTitleNote(ctx context.Context, req *brzrpc.ChangeTitleNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTitleNote"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.UpdateTitleNote(ctx, req.GetIdNote(), req.GetIdUser(), req.GetTitle())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) GetNote(ctx context.Context, req *brzrpc.UserNoteId) (*brzrpc.NoteWithBlocks, error) {
	const op = "block.note.grpc.Get"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.GetNote(ctx, req.GetNoteId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.FromNoteWithBlocksDb(res.(*domain.NoteWithBlocks)), nil
}

func (s *ServerAPI) GetAllNotes(ctx context.Context, req *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetAllNotes"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.GetNoteListByUser(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.FromNotePartsDb(res.(*domain.NoteParts)), nil
}

func (s *ServerAPI) GetNotesByTag(ctx context.Context, req *brzrpc.UserTagId) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetNotesByTag"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.GetNoteListByTag(ctx, req.GetTagId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.FromNotePartsDb(res.(*domain.NoteParts)), nil
}

func (s *ServerAPI) CreateNote(ctx context.Context, req *brzrpc.Note) (*emptypb.Empty, error) {
	const op = "block.note.grpc.Create"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.CreateNote(ctx, domain.ToNoteDb(req))
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) AddTagToNote(ctx context.Context, req *brzrpc.NoteTagUserId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.AddTagToNote"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.AddTagToNote(ctx, req.GetNoteId(), req.GetTagId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) RemoveTagFromNote(ctx context.Context, req *brzrpc.NoteTagUserId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.RemoveTagFromNote"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.RemoveTagFromNote(ctx, req.GetNoteId(), req.GetTagId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
