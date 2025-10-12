package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CleanTrash(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "block.note.grpc.CleanTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) ChangeTitleNote(context.Context, *brzrpc.ChangeTitleNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeTitleNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) NoteToTrash(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.NoteToTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) NoteFromTrash(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.NoteFromTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) FindNoteInTrash(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.FindNoteInTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetNote(context.Context, *brzrpc.Id) (*brzrpc.Note, error) {
	const op = "block.note.grpc.GetNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) CreateNote(context.Context, *brzrpc.Note) (*emptypb.Empty, error) {
	const op = "block.note.grpc.CreateNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetAllBlocksInNote(context.Context, *brzrpc.Id) (*brzrpc.Blocks, error) {
	const op = "block.note.grpc.GetAllBlocksInNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetAllNotes(context.Context, *brzrpc.GetAllNotesRequest) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetAllNotes"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetNotesByTag(context.Context, *brzrpc.GetNotesByTagRequest) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetNotesByTag"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetNotesFromTrash(context.Context, *emptypb.Empty) (*brzrpc.Notes, error) {
	const op = "block.note.grpc.GetNotesFromTrash"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) AddTagToNote(context.Context, *brzrpc.AddTagToNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.AddTagToNote"
	log.Info(op, "")

	return nil, nil
}
