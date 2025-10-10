package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CleanTrash(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) ChangeTitleNote(context.Context, *brzrpc.ChangeTitleNoteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) NoteToTrash(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) NoteFromTrash(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) FindNoteInTrash(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) GetNote(context.Context, *brzrpc.Id) (*brzrpc.Note, error) {
	return nil, nil
}
func (s *ServerAPI) CreateNote(context.Context, *brzrpc.Note) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) GetAllBlocksInNote(context.Context, *brzrpc.Id) (*brzrpc.Blocks, error) {
	return nil, nil
}
func (s *ServerAPI) GetAllNotes(context.Context, *brzrpc.GetAllNotesRequest) (*brzrpc.Notes, error) {
	return nil, nil
}
func (s *ServerAPI) GetNotesByTag(context.Context, *brzrpc.GetNotesByTagRequest) (*brzrpc.Notes, error) {
	return nil, nil
}
func (s *ServerAPI) GetNotesFromTrash(context.Context, *emptypb.Empty) (*brzrpc.Notes, error) {
	return nil, nil
}
func (s *ServerAPI) AddTagToNote(context.Context, *brzrpc.AddTagToNoteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
