package api

import (
	"context"
	"errors"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"

	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UNIFIED

func (s *ServerAPI) DeleteBlock(ctx context.Context, req *brzrpc.NoteBlockId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.DeleteBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.noteAPI.DeleteBlock(ctx, req.GetNoteId(), req.GetBlockId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) GetBlock(ctx context.Context, req *brzrpc.BlockId) (*brzrpc.Block, error) {
	const op = "block.note.grpc.GetBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.blocksAPI.Get(ctx, req.GetBlockId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.Block), nil
}

//TYPES

func (s *ServerAPI) CreateBlock(ctx context.Context, req *brzrpc.CreateBlockRequest) (*brzrpc.Id, error) {
	const op = "block.note.grpc.createBlock"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.noteAPI.CreateBlock(ctx, req.GetType(), req.GetNoteId(), req.GetData().AsMap(), int(req.GetPos()))
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.Id{Id: res.(string)}, nil
}
func (s *ServerAPI) OpBlock(ctx context.Context, req *brzrpc.OpBlockRequest) (*emptypb.Empty, error) {
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
			if err := s.noteAPI.updateUpdatedAt(ctx, b.GetNoteId()); err != nil {
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

func (s *ServerAPI) GetBlockAsFirst(ctx context.Context, req *brzrpc.BlockId) (*brzrpc.StringResponse, error) {
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

	return &brzrpc.StringResponse{Value: res.(string)}, nil
}

func (s *ServerAPI) ChangeTypeBlock(ctx context.Context, req *brzrpc.ChangeTypeBlockRequest) (*emptypb.Empty, error) {
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
