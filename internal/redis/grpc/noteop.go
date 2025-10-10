package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetNotesByUser(context.Context, *brzrpc.UserId) (*brzrpc.Notes, error) {
	return nil, nil
}
func (s *ServerAPI) GetNoteListByUser(context.Context, *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	return nil, nil
}
func (s *ServerAPI) GetNotesFromTrashByUser(context.Context, *brzrpc.UserId) (*brzrpc.Notes, error) {
	return nil, nil
}
func (s *ServerAPI) SetNotesByUser(context.Context, *brzrpc.NotesByUser) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) SetNotesFromTrashByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) SetNoteListByUser(context.Context, *brzrpc.NoteListByUser) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) RmNotesByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) RmNotesFromTrashByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	return nil, nil
}
