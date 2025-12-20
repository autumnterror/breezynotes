package grpc

import (
	"context"
	"errors"
	brzrpc2 "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/redis/redis"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"slices"
)

func (s *ServerAPI) ensureCreate(ctx context.Context, idUser string) bool {
	if err := s.rds.CheckSession(ctx, idUser); err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			if err := s.rds.CreateSession(ctx, idUser); err != nil {
				return false
			}
			return true
		}
		return false
	}
	return true
}

func (s *ServerAPI) GetNoteByUser(ctx context.Context, req *brzrpc2.UserNoteId) (*brzrpc2.NoteWithBlocks, error) {
	const op = "redis.grpc.GetNoteByUser"
	log.Info(op, "")
	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		nts, err := s.rds.GetSessionNotes(ctx, req.GetUserId())
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

		for _, n := range nts {
			if n.GetId() == req.GetNoteId() {
				res <- views.ResRPC{
					Res: n,
					Err: nil,
				}
				return
			}
		}
		res <- views.ResRPC{
			Res: nil,
			Err: status.Error(codes.NotFound, "not found in cache"),
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc2.NoteWithBlocks), nil
}
func (s *ServerAPI) GetNoteListByUser(ctx context.Context, req *brzrpc2.UserId) (*brzrpc2.NoteParts, error) {
	const op = "redis.grpc.GetNoteListByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		nts, err := s.rds.GetSessionNoteList(ctx, req.GetUserId())
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
			Res: nts,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc2.NoteParts{Items: res.([]*brzrpc2.NotePart)}, nil
}
func (s *ServerAPI) GetNotesFromTrashByUser(ctx context.Context, req *brzrpc2.UserId) (*brzrpc2.NoteParts, error) {
	const op = "redis.grpc.GetNotesFromTrashByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		nts, err := s.rds.GetSessionNoteTrash(ctx, req.GetUserId())
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
			Res: nts,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc2.NoteParts{Items: res.([]*brzrpc2.NotePart)}, nil
}

func (s *ServerAPI) SetNoteByUser(ctx context.Context, req *brzrpc2.NoteByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNotesByUser"
	log.Info(op, "")
	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		nts, err := s.rds.GetSessionNotes(ctx, req.GetUserId())
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
		nts = append(nts, req.GetNote())

		err = s.rds.SetSessionNotes(ctx, req.GetUserId(), nts)
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
func (s *ServerAPI) SetNotesFromTrashByUser(ctx context.Context, req *brzrpc2.NoteListByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNotesFromTrashByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		err := s.rds.SetSessionNoteTrash(ctx, req.GetUserId(), req.GetItems())
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
func (s *ServerAPI) SetNoteListByUser(ctx context.Context, req *brzrpc2.NoteListByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNoteListByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		err := s.rds.SetSessionNoteList(ctx, req.GetUserId(), req.GetItems())
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
func (s *ServerAPI) RmNoteByUser(ctx context.Context, req *brzrpc2.UserNoteId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNotesByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		nts, err := s.rds.GetSessionNotes(ctx, req.GetUserId())
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
		idx := slices.IndexFunc(nts, func(b *brzrpc2.NoteWithBlocks) bool {
			return b.Id == req.GetNoteId()
		})

		if idx == -1 {
			res <- views.ResRPC{
				Res: nil,
				Err: nil,
			}
			return
		}

		nts = append(nts[:idx], nts[idx+1:]...)

		err = s.rds.SetSessionNotes(ctx, req.GetUserId(), nts)
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
func (s *ServerAPI) RmNotesFromTrashByUser(ctx context.Context, req *brzrpc2.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNotesFromTrashByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		err := s.rds.SetSessionNoteTrash(ctx, req.GetUserId(), nil)
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

func (s *ServerAPI) RmNoteListByUser(ctx context.Context, req *brzrpc2.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNoteListByUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if !s.ensureCreate(ctx, req.GetUserId()) {
		return nil, status.Error(codes.Internal, "can't create session")
	}

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		err := s.rds.SetSessionNoteList(ctx, req.GetUserId(), nil)
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
