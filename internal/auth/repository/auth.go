package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	Authentication(ctx context.Context, email, login, pw string) (string, error)
}

// Authentication search user login and password in database and compare
func (d Driver) Authentication(ctx context.Context, email, login, pw string) (string, error) {
	const op = "psql.Authentication"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	var (
		query string
		arg   string
	)

	switch {
	case login != "":
		query = `SELECT id, password FROM users WHERE login = $1`
		arg = login
	case email != "":
		query = `SELECT id, password FROM users WHERE email = $1`
		arg = email
	default:
		return "", format.Error(op, domain.ErrWrongInput)
	}

	var hashed string
	var id string
	if err := d.Driver.QueryRowContext(ctx, query, arg).Scan(&id, &hashed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrNotFound
		}
		return "", format.Error(op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.ErrPasswordIncorrect
		}
		return "", format.Error(op, err)
	}

	return id, nil
}
