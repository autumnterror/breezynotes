package grpc

import (
	"context"
	"errors"
	brzrpc2 "github.com/autumnterror/breezynotes/api/proto/gen"

	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/blocks"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UNIFIED

func (s *ServerAPI) DeleteBlock(ctx context.Context, req *brzrpc2.NoteBlockId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.DeleteBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.blocksAPI.Delete(ctx, req.GetBlockId()); err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, mongo.ErrNotFound):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.NotFound, err.Error()),
				}
			default:
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}
		if err := s.noteAPI.DeleteBlock(ctx, req.GetNoteId(), req.GetBlockId()); err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, mongo.ErrNotFound):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.NotFound, err.Error()),
				}
			default:
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
func (s *ServerAPI) GetBlock(ctx context.Context, req *brzrpc2.BlockId) (*brzrpc2.Block, error) {
	const op = "block.note.grpc.GetBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		b, err := s.blocksAPI.Get(ctx, req.GetBlockId())
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
			Res: b,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc2.Block), nil
}

//TYPES

func (s *ServerAPI) CreateBlock(ctx context.Context, req *brzrpc2.CreateBlockRequest) (*brzrpc2.Id, error) {
	const op = "block.note.grpc.createBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		id, err := s.blocksAPI.Create(ctx, req.GetType(), req.GetNoteId(), req.GetData().AsMap())
		if err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, blocks.ErrTypeNotDefined):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Unknown, err.Error()),
				}
			default:
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}
		err = s.noteAPI.InsertBlock(ctx, req.GetNoteId(), id, int(req.GetPos()))
		if err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, mongo.ErrNotFound):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.NotFound, err.Error()),
				}
			default:
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
			return
		}
		res <- views.ResRPC{
			Res: id,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc2.Id{Id: res.(string)}, nil
}
func (s *ServerAPI) OpBlock(ctx context.Context, req *brzrpc2.OpBlockRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.OpBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.blocksAPI.OpBlock(ctx, req.GetId(), req.GetOp(), req.GetData().AsMap()); err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, blocks.ErrTypeNotDefined):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Unknown, err.Error()),
				}
			default:
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}

			return
		}

		if b, err := s.blocksAPI.Get(ctx, req.GetId()); err == nil {
			if err := s.noteAPI.UpdateUpdatedAt(ctx, b.GetNoteId()); err != nil {
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}
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

func (s *ServerAPI) GetBlockAsFirst(ctx context.Context, req *brzrpc2.BlockId) (*brzrpc2.StringResponse, error) {
	const op = "block.note.grpc.GetBlockAsFirst"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		bs, err := s.blocksAPI.GetAsFirst(ctx, req.GetBlockId())
		if err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, blocks.ErrTypeNotDefined):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Unknown, err.Error()),
				}
			default:
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Internal, err.Error()),
				}
			}

			return
		}
		res <- views.ResRPC{
			Res: bs,
			Err: nil,
		}
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc2.StringResponse{Value: res.(string)}, nil
}

func (s *ServerAPI) ChangeTypeBlock(ctx context.Context, req *brzrpc2.ChangeTypeBlockRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTypeBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.blocksAPI.ChangeType(ctx, req.GetId(), req.GetNewType()); err != nil {
			log.Warn(op, "", err)
			switch {
			case errors.Is(err, blocks.ErrTypeNotDefined):
				res <- views.ResRPC{
					Res: nil,
					Err: status.Error(codes.Unknown, err.Error()),
				}
			default:
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
