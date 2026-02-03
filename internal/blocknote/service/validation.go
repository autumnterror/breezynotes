package service

import (
	"errors"
	"fmt"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
)

var (
	ErrBadServiceCheck = errors.New("bad service check")
)

func wrapServiceCheck(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", ErrBadServiceCheck, err)
}

func idValidation(id string) error {
	if uid.Validate(id) {
		return nil
	}
	return errors.New("id not in uuid")
}

func stringEmpty(s string) bool {
	if s == "" {
		return true
	}
	return false
}

func digEmpty[T int64 | int | float32 | float64](s T) bool {
	if s == 0 {
		return true
	}
	return false
}

func tagValidation(t *domain.Tag) error {
	if err := idValidation(t.Id); err != nil {
		return err
	}
	if err := idValidation(t.UserId); err != nil {
		return err
	}
	if stringEmpty(t.Title) {
		return errors.New("title is empty")
	}
	if stringEmpty(t.Color) {
		return errors.New("email is empty")
	}
	if stringEmpty(t.Emoji) {
		return errors.New("photo is empty")
	}

	return nil
}

func noteValidation(n *domain.Note) error {
	var errs []error

	validateIDs := func(ids []string, fieldName string) {
		for _, id := range ids {
			if err := idValidation(id); err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", fieldName, err))
			}
		}
	}

	validateIDs(n.Blocks, "blocks")
	validateIDs(n.Readers, "readers")
	validateIDs(n.Editors, "editors")

	if err := idValidation(n.Id); err != nil {
		errs = append(errs, fmt.Errorf("id: %w", err))
	}
	if stringEmpty(n.Title) {
		errs = append(errs, errors.New("title is empty"))
	}
	if stringEmpty(n.Author) {
		errs = append(errs, errors.New("author is empty"))
	}
	if n.Tag != nil {
		if err := tagValidation(n.Tag); err != nil {
			errs = append(errs, err)
		}
	}
	if digEmpty(n.UpdatedAt) {
		errs = append(errs, errors.New("updated_at is empty"))
	}
	if digEmpty(n.CreatedAt) {
		errs = append(errs, errors.New("created_at is empty"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
