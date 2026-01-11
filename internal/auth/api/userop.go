package api

import (
	"context"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) DeleteUser(ctx context.Context, r *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "grpc.DeleteUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.Delete(ctx, r.GetUserId())
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) UpdateAbout(ctx context.Context, r *brzrpc.UpdateAboutRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdateAbout"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.UpdateAbout(ctx, r.GetId(), r.GetNewAbout())
	})
	if err != nil {
		return nil, format.Error(op, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) UpdateEmail(ctx context.Context, r *brzrpc.UpdateEmailRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdateEmail"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.UpdateEmail(ctx, r.GetId(), r.GetNewEmail())
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) UpdatePhoto(ctx context.Context, r *brzrpc.UpdatePhotoRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdatePhoto"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.UpdatePhoto(ctx, r.GetId(), r.GetNewPhoto())
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) ChangePasswd(ctx context.Context, r *brzrpc.ChangePasswordRequest) (*emptypb.Empty, error) {
	const op = "grpc.ChangePasswd"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.UpdatePassword(ctx, r.GetId(), r.GetNewPassword())
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) CreateUser(ctx context.Context, u *brzrpc.User) (*emptypb.Empty, error) {
	const op = "grpc.CreateUser"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.API.Create(ctx, domain.UserFromRpc(u))
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetUserDataFromToken(ctx context.Context, t *brzrpc.Token) (*brzrpc.User, error) {
	const op = "grpc.GetUserDataFromToken"
	log.Info(op, "")

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.API.GetUserDataFromToken(ctx, t.GetValue())
	})
	if err != nil {
		return nil, format.Error(op, err)
	}

	return domain.UserToRpc(res.(*domain.User)), nil
}
