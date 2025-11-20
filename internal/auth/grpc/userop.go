package grpc

import (
	"context"
	"database/sql"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/psql"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) DeleteUser(ctx context.Context, r *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "auth.grpc.DeleteUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.Delete(ctx, r.GetId()); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "user does not exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
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

func (s *ServerAPI) UpdateAbout(ctx context.Context, r *brzrpc.UpdateAboutRequest) (*emptypb.Empty, error) {
	const op = "auth.grpc.UpdateAbout"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.UpdateAbout(ctx, r.GetId(), r.GetNewAbout()); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "user does not exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
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
func (s *ServerAPI) UpdateEmail(ctx context.Context, r *brzrpc.UpdateEmailRequest) (*emptypb.Empty, error) {
	const op = "auth.grpc.UpdateEmail"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.UpdateEmail(ctx, r.GetId(), r.GetNewEmail()); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "user does not exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
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
func (s *ServerAPI) UpdatePhoto(ctx context.Context, r *brzrpc.UpdatePhotoRequest) (*emptypb.Empty, error) {
	const op = "auth.grpc.UpdatePhoto"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.UpdatePhoto(ctx, r.GetId(), r.GetNewPhoto()); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "user does not exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
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
func (s *ServerAPI) ChangePasswd(ctx context.Context, r *brzrpc.ChangePasswordRequest) (*emptypb.Empty, error) {
	const op = "auth.grpc.ChangePasswd"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.UpdatePassword(ctx, r.GetId(), r.GetNewPassword()); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "user does not exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
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
func (s *ServerAPI) CreateUser(ctx context.Context, u *brzrpc.User) (*emptypb.Empty, error) {
	const op = "auth.grpc.CreateUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := opWithContext(ctx, func(res chan views.ResRPC) {
		if err := s.UserAPI.Create(ctx, u); err != nil {
			switch {
			case errors.Is(err, psql.ErrAlreadyExist):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.AlreadyExists, "user already exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
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
func (s *ServerAPI) GetUserDataFromToken(ctx context.Context, t *brzrpc.Token) (*brzrpc.User, error) {
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

		info, err := s.UserAPI.GetInfo(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				log.Warn(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.NotFound, "user does not exist")}
			default:
				log.Error(op, "", err)
				res <- views.ResRPC{Res: nil, Err: status.Error(codes.Internal, "check logs")}
			}
			return
		}

		res <- views.ResRPC{Res: info, Err: nil}
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.User), nil
}
