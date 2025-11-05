package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ChangeTitleNote(ctx context.Context, req *brzrpc.ChangeTitleNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTitleNote"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) GetNote(ctx context.Context, req *brzrpc.Id) (*brzrpc.Note, error) {
	const op = "block.note.grpc.Get"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) CreateNote(ctx context.Context, req *brzrpc.Note) (*emptypb.Empty, error) {
	const op = "block.note.grpc.Create"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetAllBlocksInNote(ctx context.Context, req *brzrpc.Id) (*brzrpc.Blocks, error) {
	const op = "block.note.grpc.GetAllBlocksInNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetAllNotes(ctx context.Context, req *brzrpc.GetAllNotesRequest) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetAllNotes"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetNotesByTag(ctx context.Context, req *brzrpc.GetNotesByTagRequest) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetNotesByTag"
	log.Info(op, "")

	return nil, nil
}

func (s *ServerAPI) AddTagToNote(ctx context.Context, req *brzrpc.AddTagToNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.AddTagToNote"
	log.Info(op, "")

	return nil, nil
}
