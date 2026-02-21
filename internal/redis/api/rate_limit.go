package api

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
)

func (s *ServerAPI) RateLimit(ctx context.Context, req *brzrpc.RateLimitRequest) (*brzrpc.RateLimitResponse, error) {
	const op = "redis.grpc.GetNoteByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		count, ttl, err := s.rds.RateLimit(ctx, req.GetKey(), req.GetWindowMilliseconds())
		if err != nil {
			return nil, err
		}

		return &brzrpc.RateLimitResponse{
			Count: count,
			Ttl:   ttl,
		}, nil
	})

	if err != nil {
		return nil, format.Error(op, err)
	}

	return res.(*brzrpc.RateLimitResponse), nil
}
