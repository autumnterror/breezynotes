package api

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ShareNote(ctx context.Context, req *brzrpc.ShareNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ShareNote"
	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.ShareNote(ctx, req.GetNoteId(), req.GetUserId(), req.GetUserIdToShare(), req.GetRole())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}
	return nil, nil
}
func (s *ServerAPI) ChangeUserRole(ctx context.Context, req *brzrpc.ChangeUserRoleRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeUserRole"
	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.ChangeUserRole(ctx, req.GetNoteId(), req.GetUserId(), req.GetUserIdToChange(), req.GetNewRole())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}
	return nil, nil
}
