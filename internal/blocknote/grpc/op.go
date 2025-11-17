package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) ShareNote(ctx context.Context, req *brzrpc.ShareNoteRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ShareNote"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) ChangeUserRole(context.Context, *brzrpc.ChangeUserRoleRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.ChangeUserRole"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) Healthz(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "block.note.grpc.healthz"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	if err := s.noteAPI.Healthz(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}

func opWithContext(ctx context.Context, op func(chan views.ResRPC)) (interface{}, error) {
	res := make(chan views.ResRPC, 1)

	go op(res)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "context deadline")
	case r := <-res:
		return r.Res, r.Err
	}
}
