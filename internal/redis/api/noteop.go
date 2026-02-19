package api

import (
	"context"
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/redis/domain"
	"github.com/autumnterror/breezynotes/internal/redis/repository"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"slices"
)

func (s *ServerAPI) ensureCreate(ctx context.Context, idUser string) bool {
	if err := s.rds.CheckSession(ctx, idUser); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			if err := s.rds.CreateSession(ctx, idUser); err != nil {
				return false
			}
			return true
		}
		return false
	}
	return true
}

func (s *ServerAPI) GetNoteByUser(ctx context.Context, req *brzrpc.UserNoteId) (*brzrpc.NoteWithBlocks, error) {
	const op = "redis.grpc.GetNoteByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		nts, err := s.rds.GetSessionNotes(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}

		for _, n := range nts {
			if n.GetId() == req.GetNoteId() {
				return n, err
			}
		}
		return nil, domain.ErrNotFound
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.NoteWithBlocks), nil
}
func (s *ServerAPI) GetNoteListByUser(ctx context.Context, req *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "redis.grpc.GetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.rds.GetSessionNoteList(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.NoteParts{Items: res.([]*brzrpc.NotePart)}, nil
}
func (s *ServerAPI) GetNotesFromTrashByUser(ctx context.Context, req *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "redis.grpc.GetNotesFromTrashByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.rds.GetSessionNoteTrash(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.NoteParts{Items: res.([]*brzrpc.NotePart)}, nil
}

func (s *ServerAPI) SetNoteByUser(ctx context.Context, req *brzrpc.NoteByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNotesByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		nts, err := s.rds.GetSessionNotes(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}
		nts = append(nts, req.GetNote())

		return nil, s.rds.SetSessionNotes(ctx, req.GetUserId(), nts)
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) SetNotesFromTrashByUser(ctx context.Context, req *brzrpc.NoteListByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNotesFromTrashByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.rds.SetSessionNoteTrash(ctx, req.GetUserId(), req.GetItems())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) SetNoteListByUser(ctx context.Context, req *brzrpc.NoteListByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.rds.SetSessionNoteList(ctx, req.GetUserId(), req.GetItems())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) RmNoteByUser(ctx context.Context, req *brzrpc.UserNoteId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNotesByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		nts, err := s.rds.GetSessionNotes(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}

		idx := slices.IndexFunc(nts, func(b *brzrpc.NoteWithBlocks) bool {
			return b.Id == req.GetNoteId()
		})
		if idx == -1 {
			return nil, domain.ErrNotFound
		}

		nts = append(nts[:idx], nts[idx+1:]...)

		return nil, s.rds.SetSessionNotes(ctx, req.GetUserId(), nts)
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) RmNotesFromTrashByUser(ctx context.Context, req *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNotesFromTrashByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.rds.SetSessionNoteTrash(ctx, req.GetUserId(), nil)
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) RmNoteListByUser(ctx context.Context, req *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNoteListByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.rds.SetSessionNoteList(ctx, req.GetUserId(), nil)
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

//func (s *ServerAPI) CleanNoteById(ctx context.Context, req *brzrpc.NoteId) error {
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		return nil, s.rds.SetSessionNoteList(ctx, req.GetUserId(), nil)
//	})
//
//	if err != nil {
//		return nil, format.Error(op, err)
//	}
//
//	return nil, nil
//}
