package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ShareNote(context.Context, *brzrpc.ShareNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ShareNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) ChangeUserRole(context.Context, *brzrpc.ChangeUserRoleRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeUserRole"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "block.note.grpc.healthz"
	log.Info(op, "")

	return nil, nil
}
