package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GetNotesByUser(context.Context, *brzrpc.UserId) (*brzrpc.Notes, error) {
	const op = "redis.grpc.GetNotesByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetNoteListByUser(context.Context, *brzrpc.UserId) (*brzrpc.NoteParts, error) {
	const op = "redis.grpc.GetNoteListByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) GetNotesFromTrashByUser(context.Context, *brzrpc.UserId) (*brzrpc.Notes, error) {
	const op = "redis.grpc.GetNotesFromTrashByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) SetNotesByUser(context.Context, *brzrpc.NotesByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNotesByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) SetNotesFromTrashByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNotesFromTrashByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) SetNoteListByUser(context.Context, *brzrpc.NoteListByUser) (*emptypb.Empty, error) {
	const op = "redis.grpc.SetNoteListByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) RmNotesByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNotesByUser"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) RmNotesFromTrashByUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "redis.grpc.RmNotesFromTrashByUser"
	log.Info(op, "")

	return nil, nil
}
