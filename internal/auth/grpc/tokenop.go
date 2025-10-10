package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) GenerateAccessToken(context.Context, *brzrpc.UserId) (*brzrpc.Token, error) {
	return nil, nil
}
func (s *ServerAPI) GenerateRefreshToken(context.Context, *brzrpc.UserId) (*brzrpc.Token, error) {
	return nil, nil
}
func (s *ServerAPI) GenerateTokens(context.Context, *brzrpc.UserId) (*brzrpc.Tokens, error) {
	return nil, nil
}
func (s *ServerAPI) Refresh(context.Context, *brzrpc.Token) (*brzrpc.Token, error) {
	return nil, nil
}
func (s *ServerAPI) CheckToken(context.Context, *brzrpc.Token) (*emptypb.Empty, error) {
	return nil, nil
}
