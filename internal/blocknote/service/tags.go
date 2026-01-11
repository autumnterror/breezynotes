package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

func (s *TagsService) Get(ctx context.Context, id string) (*domain.Tag, error) {
	const op = "service.GetTag"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.tgs.Get(ctx, id)
}

func (s *TagsService) GetAllById(ctx context.Context, id string) (*domain.Tags, error) {
	const op = "service.GetAllByIdTag"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return s.tgs.GetAllById(ctx, id)
}

func (s *TagsService) Create(ctx context.Context, t *domain.Tag) error {
	const op = "service.CreateTag"
	if err := tagValidation(t); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.tgs.Create(ctx, t)
	})

	return err
}

func (s *TagsService) Delete(ctx context.Context, id string) error {
	const op = "service.DeleteTag"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.tgs.Delete(ctx, id)
	})

	return err
}
func (s *TagsService) UpdateTitle(ctx context.Context, id, nTitle string) error {
	const op = "service.UpdateTitleTag"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(nTitle) {
		return wrapServiceCheck(op, errors.New("nTitle is empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.tgs.UpdateTitle(ctx, id, nTitle)
	})

	return err
}

func (s *TagsService) UpdateColor(ctx context.Context, id, nColor string) error {
	const op = "service.UpdateColorTag"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(nColor) {
		return wrapServiceCheck(op, errors.New("nColor is empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.tgs.UpdateTitle(ctx, id, nColor)
	})

	return err
}

func (s *TagsService) UpdateEmoji(ctx context.Context, id, nEmoji string) error {
	const op = "service.UpdateEmojiTag"
	if err := idValidation(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(nEmoji) {
		return wrapServiceCheck(op, errors.New("nEmoji is empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.tgs.UpdateTitle(ctx, id, nEmoji)
	})

	return err
}
