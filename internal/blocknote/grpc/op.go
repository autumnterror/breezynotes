package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ShareNote(context.Context, *brzrpc.ShareNoteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) ChangeUserRole(context.Context, *brzrpc.ChangeUserRoleRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
