package service

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
)

//func (s *BN) GetTag(ctx context.Context, id string) (*domain.Tag, error) {
//	const op = "service.GetTag"
//	if err := idValidation(id); err != nil {
//		return nil, wrapServiceCheck(op, err)
//	}
//
//	return s.tgs.Get(ctx, id)
//}

func (s *BN) GetAllByIdTag(ctx context.Context, id string) (*domain.Tags, error) {
	const op = "service.GetAllByIdTag"
	if err := idValidation(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return s.tgs.GetAllById(ctx, id)
}

func (s *BN) CreateTag(ctx context.Context, t *domain.Tag) error {
	const op = "service.CreateTag"
	if err := tagValidation(t); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.tgs.Create(ctx, t)
	})

	return err
}

func (s *BN) DeleteTag(ctx context.Context, idTag, idUser string) error {
	const op = "service.DeleteTag"
	if err := idValidation(idTag); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		tag, err := s.tgs.Get(ctx, idTag)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if tag.UserId != idUser {
			return nil, domain.ErrUnauthorized
		}
		return nil, s.tgs.Delete(ctx, idTag)
	})

	return err
}
func (s *BN) UpdateTitleTag(ctx context.Context, idTag, idUser, nTitle string) error {
	const op = "service.UpdateTitleTag"
	if err := idValidation(idTag); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(nTitle) {
		return wrapServiceCheck(op, errors.New("title is empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		tag, err := s.tgs.Get(ctx, idTag)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if tag.UserId != idUser {
			return nil, domain.ErrUnauthorized
		}
		return nil, s.tgs.UpdateTitle(ctx, idTag, nTitle)
	})

	return err
}

func (s *BN) UpdateColorTag(ctx context.Context, idTag, idUser, nColor string) error {
	const op = "service.UpdateColorTag"
	if err := idValidation(idTag); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(nColor) {
		return wrapServiceCheck(op, errors.New("color is empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		tag, err := s.tgs.Get(ctx, idTag)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if tag.UserId != idUser {
			return nil, domain.ErrUnauthorized
		}
		return nil, s.tgs.UpdateColor(ctx, idTag, nColor)
	})

	return err
}

func (s *BN) UpdateEmojiTag(ctx context.Context, idTag, idUser, nEmoji string) error {
	const op = "service.UpdateEmojiTag"
	if err := idValidation(idTag); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}
	if stringEmpty(nEmoji) {
		return wrapServiceCheck(op, errors.New("emoji is empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		tag, err := s.tgs.Get(ctx, idTag)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if tag.UserId != idUser {
			return nil, domain.ErrUnauthorized
		}
		return nil, s.tgs.UpdateEmoji(ctx, idTag, nEmoji)
	})

	return err
}
