package api

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/breezynotes/internal/auth/service"
	"github.com/autumnterror/breezynotes/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
			case errors.Is(r.err, domain.ErrAlreadyExists):
				return nil, status.Error(codes.AlreadyExists, r.err.Error())
			case errors.Is(r.err, domain.ErrForeignKey):
				return nil, status.Error(codes.FailedPrecondition, r.err.Error())
			case errors.Is(r.err, domain.ErrTokenExpired):
				return nil, status.Error(codes.Unauthenticated, r.err.Error())
			case errors.Is(r.err, domain.ErrTokenWrongType), errors.Is(r.err, domain.ErrTokenInvalid):
				return nil, status.Error(codes.InvalidArgument, r.err.Error())
			case errors.Is(r.err, service.ErrBadServiceCheck), errors.Is(r.err, domain.ErrWrongInput), errors.Is(r.err, domain.ErrPasswordIncorrect):
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
