package api

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/service"
	"github.com/autumnterror/breezynotes/pkg/log"
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

func handleCRUDResponse(ctx context.Context, op string, action func() (any, error)) (any, error) {
	type resV struct {
		res any
		err error
	}
	res := make(chan resV, 1)
	go func() {
		r, err := action()
		res <- resV{r, err}
	}()
	select {
	case <-ctx.Done():
		log.Error(op, "Context dead", ctx.Err())
		return nil, status.Error(codes.DeadlineExceeded, "Context dead")
	case r := <-res:
		if r.err != nil {
			log.Error(op, "", r.err)
			switch {
			case errors.Is(r.err, domain.ErrNotFound):
				return nil, status.Error(codes.NotFound, r.err.Error())
			case errors.Is(r.err, domain.ErrTypeNotDefined):
				return nil, status.Error(codes.FailedPrecondition, r.err.Error())
			case errors.Is(r.err, domain.ErrAlreadyUsed):
				return nil, status.Error(codes.PermissionDenied, r.err.Error())
			case errors.Is(r.err, domain.ErrBadRequest), errors.Is(r.err, service.ErrBadServiceCheck):
				return nil, status.Error(codes.InvalidArgument, r.err.Error())
			default:
				return nil, status.Error(codes.Internal, "check logs")
			}
		}
		log.Green(op)
		if r.res != nil {
			return r.res, nil
		}
		return nil, nil
	}
}
