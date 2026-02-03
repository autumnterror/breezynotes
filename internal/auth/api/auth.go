package api

import (
	"context"
	"time"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
)

func (s *ServerAPI) Auth(ctx context.Context, r *brzrpc.AuthRequest) (*brzrpc.Tokens, error) {
	const op = "grpc.Auth"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		at, rt, err := s.API.Auth(ctx, r.GetEmail(), r.GetLogin(), r.GetPassword())
		if err != nil {
			return nil, err
		}
		return &brzrpc.Tokens{
			AccessToken:  at,
			RefreshToken: rt,
			ExpAccess:    time.Now().UTC().Add(s.Cfg.AccessTokenLifeTime).Unix(),
			ExpRefresh:   time.Now().UTC().Add(s.Cfg.RefreshTokenLifeTime).Unix(),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Tokens), nil
}

func (s *ServerAPI) Reg(ctx context.Context, r *brzrpc.AuthRequest) (*brzrpc.Tokens, error) {
	const op = "grpc.Reg"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		at, rt, err := s.API.Reg(ctx, r.GetEmail(), r.GetLogin(), r.GetPassword())
		if err != nil {
			return nil, err
		}
		return &brzrpc.Tokens{
			AccessToken:  at,
			RefreshToken: rt,
			ExpAccess:    time.Now().UTC().Add(s.Cfg.AccessTokenLifeTime).Unix(),
			ExpRefresh:   time.Now().UTC().Add(s.Cfg.RefreshTokenLifeTime).Unix(),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Tokens), nil
}

func (s *ServerAPI) ValidateTokens(ctx context.Context, r *brzrpc.Tokens) (*brzrpc.Token, error) {
	const op = "grpc.Reg"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.API.ValidateTokens(ctx, r.GetAccessToken(), r.GetRefreshToken())
	})

	if err != nil {
		return nil, err
	}
	if res.(string) != "" {
		return &brzrpc.Token{
			Value: res.(string),
			Exp:   time.Now().UTC().Add(s.Cfg.AccessTokenLifeTime).Unix(),
		}, nil
	}
	return nil, nil
}
