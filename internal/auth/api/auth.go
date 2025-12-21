package grpc

import (
	"context"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
)

func (s *ServerAPI) Auth(ctx context.Context, r *brzrpc.AuthRequest) (*brzrpc.UserId, error) {
	const op = "grpc.Auth"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.API.Auth(ctx, r.GetEmail(), r.GetLogin(), r.GetPassword())
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return &brzrpc.UserId{UserId: res.(string)}, nil
}
