package api

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateTag(ctx context.Context, t *brzrpc.Tag) (*emptypb.Empty, error) {
	const op = "grpc.CreateTag"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.CreateTag(ctx, domain.ToTagDb(t))
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *ServerAPI) UpdateTagTitle(ctx context.Context, req *brzrpc.UpdateTagTitleRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdateTagTitle"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.UpdateTitleTag(ctx, req.GetIdTag(), req.GetIdUser(), req.GetTitle())
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (s *ServerAPI) UpdateTagColor(ctx context.Context, req *brzrpc.UpdateTagColorRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdateTagColor"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.UpdateColorTag(ctx, req.GetIdTag(), req.GetIdUser(), req.GetColor())
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *ServerAPI) UpdateTagEmoji(ctx context.Context, req *brzrpc.UpdateTagEmojiRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdateTagEmoji"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.UpdateEmojiTag(ctx, req.GetIdTag(), req.GetIdUser(), req.GetEmoji())
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ServerAPI) UpdateTagPinned(ctx context.Context, req *brzrpc.UserTagId) (*emptypb.Empty, error) {
	const op = "grpc.UpdateTagPinned"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.UpdatePinned(ctx, req.GetTagId(), req.GetUserId())
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ServerAPI) DeleteTag(ctx context.Context, req *brzrpc.UserTagId) (*emptypb.Empty, error) {
	const op = "grpc.DeleteTag"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.DeleteTag(ctx, req.GetTagId(), req.GetUserId())
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ServerAPI) DeleteTags(ctx context.Context, req *brzrpc.UserId) (*emptypb.Empty, error) {
	const op = "grpc.DeleteTags"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.service.DeleteTags(ctx, req.GetUserId())
	})

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ServerAPI) GetTagsByUser(ctx context.Context, req *brzrpc.UserId) (*brzrpc.Tags, error) {
	const op = "grpc.GetTagsByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		r, err := s.service.GetAllByIdTag(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}
		return domain.FromTagsDb(r), nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Tags), nil
}

func (s *ServerAPI) GetPinnedTagsByUser(ctx context.Context, req *brzrpc.UserId) (*brzrpc.Tags, error) {
	const op = "grpc.GetPinnedTagsByUser"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		r, err := s.service.GetAllByIdTagPinned(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}
		return domain.FromTagsDb(r), nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*brzrpc.Tags), nil
}

//func (s *ServerAPI) GetTag(ctx context.Context, req *brzrpc.TagId) (*brzrpc.Tag, error) {
//	const op = "grpc.GetTag"
//
//	ctx, done := context.WithTimeout(ctx, waitTime)
//	defer done()
//
//	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
//		r, err := s.service.GetTag(ctx, req.GetTagId())
//		if err != nil {
//			return nil, err
//		}
//		return domain.FromTagDb(r), nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return res.(*brzrpc.Tag), nil
//}
