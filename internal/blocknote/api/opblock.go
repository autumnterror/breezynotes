package api

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetAllBlocksInNote(ctx context.Context, req *brzrpc.Strings) (*brzrpc.Blocks, error) {
	const op = "block.note.grpc.GetAllBlocksInNote"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.GetBlocks(ctx, req.GetValues())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.FromBlocksDb(res.(*domain.Blocks)), nil
}

func (s *ServerAPI) ChangeBlockOrder(ctx context.Context, req *brzrpc.ChangeBlockOrderRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeBlockOrder"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.ChangeBlockOrder(ctx, req.GetNoteId(), req.GetUserId(), int(req.GetOldOrder()), int(req.GetNewOrder()))
	})

	if err != nil {
		return nil, format.Error(op, err)
	}
	return nil, nil
}

// UNIFIED

func (s *ServerAPI) DeleteBlock(ctx context.Context, req *brzrpc.NoteBlockUserId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.DeleteBlock"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.DeleteBlock(ctx, req.GetNoteId(), req.GetBlockId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}
func (s *ServerAPI) GetBlock(ctx context.Context, req *brzrpc.NoteBlockUserId) (*brzrpc.Block, error) {
	const op = "block.note.grpc.GetBlock"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.GetBlock(ctx, req.GetBlockId(), req.GetNoteId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.Block), nil
}

//TYPES

func (s *ServerAPI) CreateBlock(ctx context.Context, req *brzrpc.CreateBlockRequest) (*brzrpc.Id, error) {
	const op = "block.note.grpc.createBlock"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.service.CreateBlock(ctx, req.GetType(), req.GetNoteId(), req.GetData().AsMap(), int(req.GetPos()), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.Id{Id: res.(string)}, nil
}
func (s *ServerAPI) OpBlock(ctx context.Context, req *brzrpc.OpBlockRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.OpBlock"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.OpBlock(ctx, req.GetBlockId(), req.GetOp(), req.GetData().AsMap(), req.GetNoteId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) ChangeTypeBlock(ctx context.Context, req *brzrpc.ChangeTypeBlockRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTypeBlock"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.ChangeTypeBlock(ctx, req.GetBlockId(), req.GetNoteId(), req.GetUserId(), req.GetNewType())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

// TODO check is it need

func (s *ServerAPI) GetBlockAsFirst(ctx context.Context, req *brzrpc.BlockId) (*brzrpc.StringResponse, error) {
	//const op = "block.note.grpc.GetBlockAsFirst"
	//
	//ctx, done := context.WithTimeout(ctx, waitTime)
	//defer done()
	//
	//res, err := handleCRUDResponse(ctx, op, func() (any, error) {
	//	return s.service.GetAsFirst(ctx, req.GetBlockId())
	//})
	//
	//res, err := opWithContext(ctx, func(res chan domain.ResRPC) {
	//	bs, err := s.blocksAPI.GetAsFirst(ctx, req.GetBlockId())
	//	if err != nil {
	//		log.Warn(op, "", err)
	//		switch {
	//		case errors.Is(err, blocks.ErrTypeNotDefined):
	//			res <- domain.ResRPC{
	//				Res: nil,
	//				Err: status.Error(codes.Unknown, err.Error()),
	//			}
	//		default:
	//			res <- domain.ResRPC{
	//				Res: nil,
	//				Err: status.Error(codes.Internal, err.Error()),
	//			}
	//		}
	//
	//		return
	//	}
	//	res <- domain.ResRPC{
	//		Res: bs,
	//		Err: nil,
	//	}
	//})
	//
	//if err != nil {
	//	return nil, format.Error(op, err)
	//}
	//
	//return &brzrpc.StringResponse{Value: res.(string)}, nil
	return nil, nil
}
