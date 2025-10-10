package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) DeleteUser(context.Context, *brzrpc.UserId) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) UpdateAbout(context.Context, *brzrpc.UpdateAboutRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) UpdateEmail(context.Context, *brzrpc.UpdateEmailRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) UpdatePhoto(context.Context, *brzrpc.UpdatePhotoRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) ChangePasswd(context.Context, *brzrpc.ChangePasswordRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) CreateUser(context.Context, *brzrpc.User) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) GetUserDataFromToken(context.Context, *brzrpc.Token) (*brzrpc.User, error) {
	return nil, nil
}
