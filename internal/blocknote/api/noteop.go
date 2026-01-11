package api

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"time"

	"github.com/autumnterror/breezynotes/pkg/utils/format"

	"github.com/autumnterror/breezynotes/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ChangeTitleNote(ctx context.Context, req *brzrpc.ChangeTitleNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTitleNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.noteAPI.UpdateTitle(ctx, req.GetId(), req.GetTitle())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) GetNote(ctx context.Context, req *brzrpc.NoteId) (*brzrpc.Note, error) {
	const op = "block.note.grpc.GetNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.noteAPI.GetNote(ctx, req.GetNoteId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.Note), nil
}

func (s *ServerAPI) GetAllNotes(ctx context.Context, req *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetAllNotes"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.noteAPI.GetNoteListByUser(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.NoteParts), nil
}

func (s *ServerAPI) GetNotesByTag(ctx context.Context, req *brzrpc.UserTagId) (*brzrpc.NoteParts, error) {
	const op = "block.note.grpc.GetNotesByTag"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.noteAPI.GetNoteListByTag(ctx, req.GetTagId(), req.GetUserId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.NoteParts), nil
}

// TODO GET ALL BLOCKS IN ONE REQUEST

func (s *ServerAPI) GetAllBlocksInNote(ctx context.Context, req *brzrpc.Strings) (*brzrpc.Blocks, error) {
	const op = "block.note.grpc.GetAllBlocksInNote"
	log.Info(op, "")
	ctx, done := context.WithTimeout(ctx, waitTime+5*time.Second)
	defer done()

	//res, err := handleCRUDResponse(ctx, op, func() (any, error) {
	//	return s.noteAPI.GetNoteListByTag(ctx, req.GetTagId(), req.GetUserId())
	//})

	//res, err := opWithContext(ctx, func(res chan views.ResRPC) {
	//	var bs []*brzrpc.Block
	//
	//	for _, id := range req.Values {
	//		b, err := s.blocksAPI.Get(ctx, id)
	//		if err != nil {
	//			log.Warn(op, "BAD BLOCK", err)
	//			continue
	//		}
	//		bs = append(bs, b)
	//	}
	//
	//	res <- views.ResRPC{
	//		Res: bs,
	//		Err: nil,
	//	}
	//})
	//
	//if err != nil {
	//	return nil, format.Error(op, err)
	//}
	//
	//return &brzrpc.Blocks{Items: res.([]*brzrpc.Block)}, nil

	return nil, nil
}

func (s *ServerAPI) CreateNote(ctx context.Context, req *brzrpc.Note) (*emptypb.Empty, error) {
	const op = "block.note.grpc.Create"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.noteAPI.Create(ctx, domain.ToNoteDb(req))
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) AddTagToNote(ctx context.Context, req *brzrpc.NoteTagId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.AddTagToNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.noteAPI.AddTagToNote(ctx, req.GetNoteId(), req.GetTagId())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return nil, nil
}

func (s *ServerAPI) RemoveTagFromNote(ctx context.Context, req *brzrpc.NoteTagId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.RemoveTagFromNote"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.noteAPI.RemoveTagFromNote(ctx, req.GetNoteId(), req.GetTagId())
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

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.noteAPI.ChangeBlockOrder(ctx, req.GetId(), int(req.GetOldOrder()), int(req.GetNewOrder()))
	})

	if err != nil {
		return nil, format.Error(op, err)
	}
	return nil, nil
}
