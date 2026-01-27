package api

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/redis/domain"

	"github.com/autumnterror/utils_go/pkg/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "redis.grpc.Healthz"

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
		return nil, status.Error(codes.DeadlineExceeded, "Context dead")
	case r := <-res:
		if r.err != nil {
			switch {
			case errors.Is(r.err, domain.ErrNotFound):
				return nil, status.Error(codes.NotFound, r.err.Error())
			default:
				log.Error(op, "", r.err)
				return nil, status.Error(codes.Internal, "check logs")
			}
		}
		if r.res != nil {
			return r.res, nil
		}
		return nil, nil
	}
}
