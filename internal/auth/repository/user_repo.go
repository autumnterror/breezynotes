package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/autumnterror/breezynotes/internal/auth/domain"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	CreateAdmin(ctx context.Context) (string, error)
	UpdatePhoto(ctx context.Context, id, np string) error
	UpdatePassword(ctx context.Context, id, newPassword string) error
	UpdateEmail(ctx context.Context, id, email string) error
	UpdateAbout(ctx context.Context, id, about string) error
	Delete(ctx context.Context, id string) error
	GetInfo(ctx context.Context, id string) (*domain.User, error)
	GetIdFromLogin(ctx context.Context, login string) (string, error)
}

func (d Driver) CreateAdmin(ctx context.Context) (string, error) {
	const op = "users.CreateAdmin"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	query := `
				INSERT INTO users (id, login, email, about, password, photo)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
	id := uid.New()
	u := &domain.User{
		Id:       id,
		Login:    "admin",
		Email:    "admin@breezy.su",
		About:    "im admin",
		Photo:    "https://i.ibb.co/0RvkgTL8/0cceeba7-ad31-42dd-b91a-eb8e9d358524.png",
		Password: "admin",
	}

	if idA, err := d.GetIdFromLogin(ctx, "admin"); err == nil {
		if idA != "" {
			return idA, nil
		}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", format.Error(op, err)
	}

	_, err = d.Driver.ExecContext(ctx, query, u.Id, u.Login, u.Email, u.About, hashedPass, u.Photo)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				return "", format.Error(op, domain.ErrAlreadyExists)
			case "23503": // foreign_key_violation
				return "", format.Error(op, domain.ErrForeignKey)
			}
		}
		return "", format.Error(op, err)
	}

	return id, nil
}

// getAll only for test
func (d Driver) getAll(ctx context.Context) ([]*domain.User, error) {
	const op = "users.getAll"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	var ls []*domain.User
	rows, err := d.Driver.QueryContext(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var us domain.User
		if err := rows.Scan(&us.Id, &us.Login, &us.Email, &us.About, &us.Password, &us.Photo); err != nil {
			log.Error(op, "rows scan error", err)
			continue
		}
		ls = append(ls, &us)
	}

	return ls, nil
}

// GetInfo get info about user by login. May send sql.ErrNoRows
func (d Driver) GetInfo(ctx context.Context, id string) (*domain.User, error) {
	const op = "users.GetInfo"
	query := `
		SELECT login,email,about, photo FROM users
		WHERE id = $1
	`
	var u domain.User
	if err := d.Driver.QueryRowContext(ctx, query, id).Scan(&u.Login, &u.Email, &u.About, &u.Photo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, format.Error(op, domain.ErrNotFound)
		}
		return nil, format.Error(op, err)
	}
	u.Id = id
	u.Password = ""

	return &u, nil
}

// GetIdFromLogin get info about user by login. May send sql.ErrNoRows
func (d Driver) GetIdFromLogin(ctx context.Context, login string) (string, error) {
	const op = "users.GetInfo"
	query := `
		SELECT id FROM users
		WHERE login = $1
	`

	var id string
	if err := d.Driver.QueryRowContext(ctx, query, login).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", format.Error(op, domain.ErrNotFound)
		}
		return "", format.Error(op, err)
	}

	return id, nil
}

// Create new user
func (d Driver) Create(ctx context.Context, u *domain.User) error {
	const op = "users.Create"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	query := `
				INSERT INTO users (id, login, email, about, password, photo)
				VALUES ($1, $2, $3, $4, $5, $6)
			`

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return format.Error(op, err)
	}

	_, err = d.Driver.ExecContext(ctx, query, u.Id, u.Login, u.Email, u.About, hashedPass, u.Photo)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				return format.Error(op, domain.ErrAlreadyExists)
			case "23503": // foreign_key_violation
				return format.Error(op, domain.ErrForeignKey)
			}
		}
		return format.Error(op, err)
	}

	return nil
}

// UpdatePassword updates user's password by user ID.
func (d Driver) UpdatePassword(ctx context.Context, id, newPassword string) error {
	const op = "users.UpdatePassword"
	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return format.Error(op, err)
	}

	res, err := d.Driver.ExecContext(ctx, `UPDATE users SET password = $1 WHERE id = $2`, hashedPass, id)
	if err != nil {
		return format.Error(op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return format.Error(op, err)
	}
	if rowsAffected == 0 {
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

// UpdatePhoto updates user's photo by user ID.
// Returns sql.ErrNoRows if user not found.
func (d Driver) UpdatePhoto(ctx context.Context, id, np string) error {
	const op = "users.UpdatePhoto"

	ctx, done := context.WithTimeout(ctx, waitTime)
	defer done()

	res, err := d.Driver.ExecContext(ctx, `UPDATE users SET photo = $1 WHERE id = $2`, np, id)
	if err != nil {
		return format.Error(op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return format.Error(op, err)
	}
	if rowsAffected == 0 {
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

// UpdateEmail updates user's email by user ID.
// Returns sql.ErrNoRows if user not found.
func (d Driver) UpdateEmail(ctx context.Context, id, email string) error {
	const op = "users.UpdateEmail"

	res, err := d.Driver.ExecContext(ctx, `UPDATE users SET email = $1 WHERE id = $2`, email, id)
	if err != nil {
		return format.Error(op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return format.Error(op, err)
	}
	if rowsAffected == 0 {
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

// UpdateAbout updates user's about section by user ID.
// Returns sql.ErrNoRows if user not found.
func (d Driver) UpdateAbout(ctx context.Context, id, about string) error {
	const op = "users.UpdateAbout"

	res, err := d.Driver.ExecContext(ctx, `UPDATE users SET about = $1 WHERE id = $2`, about, id)
	if err != nil {
		return format.Error(op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return format.Error(op, err)
	}
	if rowsAffected == 0 {
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

// Delete user. May send sql.ErrNoRows
func (d Driver) Delete(ctx context.Context, id string) error {
	const op = "users.delete"

	query := `
				DELETE FROM users
				WHERE id = $1
			`
	res, err := d.Driver.ExecContext(ctx, query, id)
	if err != nil {
		return format.Error(op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return format.Error(op, err)
	}
	if rowsAffected == 0 {
		return format.Error(op, domain.ErrNotFound)
	}
	return nil
}
