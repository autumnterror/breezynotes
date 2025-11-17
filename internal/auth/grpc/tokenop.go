package grpc

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (s *ServerAPI) GenerateAccessToken(ctx context.Context, r *brzrpc.UserId) (*brzrpc.Token, error) {
	const op = "auth.grpc.GenerateAccessToken"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		at, err := s.JwtAPI.GenerateToken(r.GetId(), jwt.TokenTypeAccess)
		if err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: at, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return &brzrpc.Token{Value: res.(string), Exp: time.Now().UTC().Add(s.cfg.AccessTokenLifeTime).Unix()}, nil
}
func (s *ServerAPI) GenerateRefreshToken(ctx context.Context, r *brzrpc.UserId) (*brzrpc.Token, error) {
	const op = "auth.grpc.GenerateRefreshToken"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		rt, err := s.JwtAPI.GenerateToken(r.GetId(), jwt.TokenTypeRefresh)
		if err != nil {
			log.Error(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			return
		}

		res <- views.ResRPC{Res: rt, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return &brzrpc.Token{Value: res.(string), Exp: time.Now().UTC().Add(s.cfg.RefreshTokenLifeTime).Unix()}, nil
}
func (s *ServerAPI) GenerateTokens(ctx context.Context, r *brzrpc.UserId) (*brzrpc.Tokens, error) {
	const op = "auth.grpc.GenerateTokens"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()
	var ac, rt *brzrpc.Token
	var err error
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		ac, err = s.GenerateAccessToken(ctx, r)
		return err
	})
	g.Go(func() error {
		rt, err = s.GenerateRefreshToken(ctx, r)
		return err
	})
	if err := g.Wait(); err != nil {
		log.Error(op, "", err)
		return nil, status.Error(codes.Internal, "check logs")
	}

	return &brzrpc.Tokens{
		AccessToken:  ac.Value,
		RefreshToken: rt.Value,
		ExpAccess:    time.Now().UTC().Add(s.cfg.AccessTokenLifeTime).Unix(),
		ExpRefresh:   time.Now().UTC().Add(s.cfg.RefreshTokenLifeTime).Unix(),
	}, nil
}
func (s *ServerAPI) Refresh(ctx context.Context, r *brzrpc.Token) (*brzrpc.Token, error) {
	const op = "auth.grpc.Refresh"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		at, err := s.JwtAPI.Refresh(r.GetValue())
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, "refresh token expired")}
			case errors.Is(err, jwt.ErrWrongType):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.InvalidArgument, "refresh token needed")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			}
			return
		}

		res <- views.ResRPC{Res: at, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return &brzrpc.Token{Value: res.(string)}, nil
}
func (s *ServerAPI) CheckToken(ctx context.Context, r *brzrpc.Token) (*emptypb.Empty, error) {
	const op = "auth.grpc.CheckToken"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		_, err := s.JwtAPI.VerifyToken(r.GetValue())
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, "token expired")}
			default:
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.InvalidArgument, "token bad")}
			}
			return
		}

		res <- views.ResRPC{Res: nil, Err: nil}
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ServerAPI) GetIdFromToken(ctx context.Context, t *brzrpc.Token) (*brzrpc.Id, error) {
	const op = "auth.grpc.GetUserDataFromToken"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := opWithContext(ctx, func(res chan views.ResRPC) {
		token, err := s.JwtAPI.VerifyToken(t.GetValue())
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, "token expired")}
			default:
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.InvalidArgument, err.Error())}
			}
			return
		}

		tp, err := s.JwtAPI.GetTypeFromToken(token)
		if err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, err.Error())}
			return
		}
		if tp != jwt.TokenTypeAccess {
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.InvalidArgument, "access token needed")}
			return
		}

		id, err := s.JwtAPI.GetIdFromToken(token)
		if err != nil {
			log.Warn(op, "", err)
			res <- views.ResRPC{Res: nil, Err: status.Error(codes.Unauthenticated, err.Error())}
			return
		}

		res <- views.ResRPC{Res: id, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return &brzrpc.Id{Id: res.(string)}, nil
}
