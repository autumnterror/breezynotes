package service

import (
	"context"
	"errors"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
)

func (s *BN) CreateNote(ctx context.Context, n *domain.Note) error {
	const op = "service.CreateNote"
	if err := noteValidation(n); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, s.nts.Create(ctx, n)
	})

	return err
}

func (s *BN) GetNote(ctx context.Context, idNote, idUser string) (*domain.NoteWithBlocks, error) {
	const op = "service.Get"
	if err := idValidation(idNote); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if idUser != "" {
		if err := idValidation(idUser); err != nil {
			return nil, wrapServiceCheck(op, err)
		}
	}

	n, err := s.nts.Get(ctx, idNote, idUser)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) && !n.IsBlog && !n.IsPublic {
		return nil, domain.ErrUnauthorized
	}

	blks, err := s.blk.GetMany(ctx, n.Blocks)
	if err != nil {
		return nil, err
	}

	return &domain.NoteWithBlocks{
		Id:        n.Id,
		Title:     n.Title,
		Blocks:    blks.Blks,
		Author:    n.Author,
		Readers:   n.Readers,
		Editors:   n.Editors,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       n.Tag,
		IsBlog:    n.IsBlog,
		IsPublic:  n.IsPublic,
	}, nil
}

func (s *BN) GetNoteListByUser(ctx context.Context, idUser string) (*domain.NoteParts, error) {
	const op = "service.GetNoteListByUser"
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNoteListByUser(ctx, idUser)
}

func (s *BN) GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain.NoteParts, error) {
	const op = "service.GetNoteListByTag"
	if err := idValidation(idTag); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return s.nts.GetNoteListByTag(ctx, idTag, idUser)
}

func (s *BN) AddTagToNote(ctx context.Context, idNote, tagId, idUser string) error {
	const op = "service.AddTagToNote"
	if err := idValidation(tagId); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote, idUser)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		tag, err := s.tgs.Get(ctx, tagId)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		err = s.nts.AddTagToNote(ctx, idNote, tag)
		if err != nil {
			return nil, err
		}
		n, err = s.nts.Get(ctx, idNote, idUser)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		return nil, nil
	})

	return err
}

func (s *BN) RemoveTagFromNote(ctx context.Context, idNote string, idUser string) error {
	const op = "service.RemoveTagFromNote"

	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote, idUser)
		if err != nil {
			return nil, domain.ErrUnauthorized
		}
		if n.Author != idUser && !alg.IsIn(idUser, n.Editors) && !alg.IsIn(idUser, n.Readers) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.RemoveTagFromNote(ctx, idNote, idUser)
	})

	return err
}

func (s *BN) UpdateTitleNote(ctx context.Context, idNote, idUser, nTitle string) error {
	const op = "service.UpdateTitleNote"

	if err := idValidation(idNote); err != nil {
		return wrapServiceCheck(op, err)
	}
	if err := idValidation(idUser); err != nil {
		return wrapServiceCheck(op, err)
	}

	if stringEmpty(nTitle) {
		return wrapServiceCheck(op, errors.New("title empty"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		n, err := s.nts.Get(ctx, idNote, idUser)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if n.Author != idUser && !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.UpdateTitle(ctx, idNote, nTitle)
	})

	return err
}

func (s *BN) Search(ctx context.Context, idUser, prompt string) (<-chan *domain.NotePart, error) {
	const op = "service.Search"
	if err := idValidation(idUser); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if stringEmpty(prompt) {
		return nil, nil
	}
	return s.nts.Search(ctx, idUser, prompt)
}

func (s *BN) ShareNote(ctx context.Context, idNote, idUser, idUserToShare, role string) error {
	const op = "service.ShareNote"

	if stringEmpty(role) {
		return wrapServiceCheck(op, errors.New("role is empty"))
	}
	if idValidation(idNote) != nil {
		return wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(idUser) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}
	if idValidation(idUserToShare) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}
	if idUser == idUserToShare {
		return wrapServiceCheck(op, errors.New("can't share for yourself"))
	}

	switch role {
	case domain.EditorRole:
	case domain.ReaderRole:
	default:
		return wrapServiceCheck(op, errors.New("role undefined"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if n, err := s.nts.Get(ctx, idNote, idUser); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		} else {
			if n.Author == idUserToShare {
				return nil, wrapServiceCheck(op, errors.New("can't share for author"))
			}
		}

		return nil, s.nts.ShareNote(ctx, idNote, idUserToShare, role)
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *BN) AddPublicNote(ctx context.Context, idNote, idUser string) error {
	const op = "service.PublicNote"

	if idValidation(idNote) != nil {
		return wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(idUser) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {

		if n, err := s.nts.Get(ctx, idNote, idUser); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author == idUser || alg.IsIn(idUser, n.Editors) || alg.IsIn(idUser, n.Readers) {
			return nil, nil
		} else if !n.IsPublic || !n.IsBlog {
			return nil, domain.ErrUnauthorized
		}

		return nil, s.nts.ShareNote(ctx, idNote, idUser, domain.ReaderRole)
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *BN) PublicNote(ctx context.Context, idNote, idUser string) error {
	const op = "service.PublicNote"

	if idValidation(idNote) != nil {
		return wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(idUser) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		isPublic := false
		if n, err := s.nts.Get(ctx, idNote, idUser); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		} else {
			if !n.IsPublic {
				isPublic = true
			}
		}
		return nil, s.nts.UpdatePublic(ctx, idNote, isPublic)
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *BN) BlogNote(ctx context.Context, idNote, idUser string) error {
	const op = "service.BlogNote"

	if idValidation(idNote) != nil {
		return wrapServiceCheck(op, errors.New("bad note id"))
	}
	if idValidation(idUser) != nil {
		return wrapServiceCheck(op, errors.New("bad user id"))
	}

	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
		isBlog := false
		if n, err := s.nts.Get(ctx, idNote, idUser); err != nil {
			return nil, domain.ErrNotFound
		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) {
			return nil, domain.ErrUnauthorized
		} else {
			if !n.IsBlog {
				isBlog = true
			}
		}
		return nil, s.nts.UpdateBlog(ctx, idNote, isBlog)
	})

	if err != nil {
		return err
	}
	return nil
}

// func (s *BN) ChangeUserRole(ctx context.Context, idNote, idUser, idUserToChange, newRole string) error {
// 	const op = "service.ChangeUserRole"
// 	if stringEmpty(newRole) {
// 		return wrapServiceCheck(op, errors.New("newRole is empty"))
// 	}
// 	if idValidation(idNote) != nil {
// 		return wrapServiceCheck(op, errors.New("bad note id"))
// 	}
// 	if idValidation(idUser) != nil {
// 		return wrapServiceCheck(op, errors.New("bad user id"))
// 	}
// 	if idValidation(idUserToChange) != nil {
// 		return wrapServiceCheck(op, errors.New("bad user id"))
// 	}
// 	if idUser == idUserToChange {
// 		return wrapServiceCheck(op, errors.New("can't change for yourself"))
// 	}

// 	switch newRole {
// 	case domain.EditorRole:
// 	case domain.ReaderRole:
// 	default:
// 		return wrapServiceCheck(op, errors.New("role undefined"))
// 	}
// 	_, err := s.tx.RunInTx(ctx, func(ctx context.Context) (interface{}, error) {
// 		if n, err := s.nts.Get(ctx, idNote, idUser); err != nil {
// 			return nil, domain.ErrNotFound
// 		} else if n.Author != idUser && !alg.IsIn(idUser, n.Editors) {
// 			return nil, domain.ErrUnauthorized
// 		} else {
// 			if n.Author == idUserToChange {
// 				return nil, wrapServiceCheck(op, errors.New("can't share for author"))
// 			}
// 		}

// 		return nil, s.nts.ChangeUserRole(ctx, idNote, idUserToChange, newRole)
// 	})

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
