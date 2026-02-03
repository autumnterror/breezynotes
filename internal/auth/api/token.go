package api

import (
	"context"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
)

func (s *ServerAPI) GetIdFromToken(ctx context.Context, t *brzrpc.Token) (*brzrpc.Id, error) {
	const op = "grpc.GetIdFromToken"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.API.GetIdFromToken(ctx, t.GetValue())
	})
	if err != nil {
		return nil, err
	}

	return &brzrpc.Id{Id: res.(string)}, nil
}

//func (s *ServerAPI) GenerateAccessToken(ctx context.Context, r *brzrpc.UserId) (*brzrpc.Token, error) {
//	const op = "grpc.GenerateAccessToken"
//
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		return s.API.GenerateAccessToken(ctx, r.GetUserId())
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return domain.TokenToRPC(res.(*domain.Token)), nil
//}
//
//func (s *ServerAPI) GenerateRefreshToken(ctx context.Context, r *brzrpc.UserId) (*brzrpc.Token, error) {
//	const op = "grpc.GenerateRefreshToken"
//
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		return s.API.GenerateRefreshToken(ctx, r.GetUserId())
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return domain.TokenToRPC(res.(*domain.Token)), nil
//}
//
//func (s *ServerAPI) GenerateTokens(ctx context.Context, r *brzrpc.UserId) (*brzrpc.Tokens, error) {
//	const op = "grpc.GenerateTokens"
//
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		return s.API.GenerateTokens(ctx, r.GetUserId())
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return domain.TokensToRPC(res.(*domain.Tokens)), nil
//}
//
//func (s *ServerAPI) Refresh(ctx context.Context, r *brzrpc.Token) (*brzrpc.Token, error) {
//	const op = "grpc.Refresh"
//
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		return s.API.Refresh(ctx, r.GetValue())
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return domain.TokenToRPC(res.(*domain.Token)), nil
//}
//
//func (s *ServerAPI) CheckToken(ctx context.Context, r *brzrpc.Token) (*emptypb.Empty, error) {
//	const op = "grpc.CheckToken"
//
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		return nil, s.API.CheckToken(ctx, r.GetValue())
//	})
//	if err != nil {
//		return nil, err
//	}
//	return &emptypb.Empty{}, nil
//}
