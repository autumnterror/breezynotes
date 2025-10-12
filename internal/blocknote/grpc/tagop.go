package grpc

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateTag(context.Context, *brzrpc.Tag) (*emptypb.Empty, error) {
	const op = "block.note.grpc.CreateTag"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) UpdateTagTitle(context.Context, *brzrpc.UpdateTagTitleRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.UpdateTagTitle"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) UpdateTagColor(context.Context, *brzrpc.UpdateTagColorRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.UpdateTagColor"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) UpdateTagEmoji(context.Context, *brzrpc.UpdateTagEmojiRequest) (*emptypb.Empty, error) {
	const op = "block.note.grpc.UpdateTagEmoji"
	log.Info(op, "")

	return nil, nil
}
func (s *ServerAPI) DeleteTag(context.Context, *brzrpc.Id) (*emptypb.Empty, error) {
	const op = "block.note.grpc.DeleteTag"
	log.Info(op, "")

	return nil, nil
}
