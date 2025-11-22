package grpc

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"time"

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
			case errors.Is(err, mongo.ErrNotFound):
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
		return nil, format.Error(op, err)
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
			switch {
			case errors.Is(err, mongo.ErrNotFound):
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
			Res: n,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.Note), nil
}

func (s *ServerAPI) GetAllNotes(ctx context.Context, req *brzrpc.GetNotesRequestPagination) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetAllNotes"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if req.GetStart() < 0 {
		return nil, status.Error(codes.InvalidArgument, "start<0")
	}

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		n, err := s.noteAPI.GetNoteListByUser(ctx, req.GetIdUser(), int(req.GetStart()), int(req.GetEnd()))
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

	return res.(*brzrpc.NoteParts), nil
}

func (s *ServerAPI) GetNotesByTag(ctx context.Context, req *brzrpc.GetNotesByTagRequestPagination) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetNotesByTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if req.GetStart() < 0 {
		return nil, status.Error(codes.InvalidArgument, "start<0")
	}

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		n, err := s.noteAPI.GetNoteListByTag(ctx, req.GetIdTag(), req.GetIdUser(), int(req.GetStart()), int(req.GetEnd()))
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

	return res.(*brzrpc.NoteParts), nil
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
		return nil, format.Error(op, err)
	}

	return nil, nil
}

// TODO check benchmark

func (s *ServerAPI) GetAllBlocksInNote(ctx context.Context, req *brzrpc.Strings) (*brzrpc.Blocks, error) {
	const op = "block.note.grpc.GetAllBlocksInNote"
	log.Info(op, "")
	ctx, done := context.WithTimeout(ctx, waitTime+5*time.Second)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		var bs []*brzrpc.Block

		for _, id := range req.Values {
			b, err := s.blocksAPI.Get(ctx, id)
			if err != nil {
				log.Warn(op, "BAD BLOCK", err)
				continue
			}
			bs = append(bs, b)
		}

		res <- views.ResRPC{
			Res: bs,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.Blocks{Items: res.([]*brzrpc.Block)}, nil
}

func (s *ServerAPI) AddTagToNote(ctx context.Context, req *brzrpc.TagToNoteRequest) (*emptypb.Empty, error) {
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

func (s *ServerAPI) RemoveTagFromNote(ctx context.Context, req *brzrpc.TagToNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.RemoveTagFromNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.noteAPI.RemoveTagFromNote(ctx, req.GetNoteId(), req.GetTagId()); err != nil {
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

func (s *ServerAPI) ChangeBlockOrder(ctx context.Context, req *brzrpc.ChangeBlockOrderRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeBlockOrder"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		log.Blue(req)
		if err := s.noteAPI.ChangeBlockOrder(ctx, req.GetId(), int(req.GetOldOrder()), int(req.GetNewOrder())); err != nil {
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
