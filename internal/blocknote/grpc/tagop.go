package grpc

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateTag(context.Context, *brzrpc.Tag) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) UpdateTagTitle(context.Context, *brzrpc.UpdateTagTitleRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) UpdateTagColor(context.Context, *brzrpc.UpdateTagColorRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) UpdateTagEmoji(context.Context, *brzrpc.UpdateTagEmojiRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ServerAPI) DeleteTag(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	return nil, nil
}
