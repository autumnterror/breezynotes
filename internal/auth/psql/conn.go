package psql

import (
	"database/sql"
	"errors"
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	_ "github.com/lib/pq"
)

type PostgresDb struct {
	Driver *sql.DB
}

type SqlRepo interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Driver struct {
	driver SqlRepo
}

func NewDriver(driver SqlRepo) *Driver {
	return &Driver{
		driver: driver,
	}
}

func MustConnect(cfg *config.Config) *PostgresDb {
	db, err := NewConnect(cfg)
	if err != nil {
		log.Panic(err)
	}
	return db
}

// NewConnect is constructor of PostgresDb. Construct with connection
func NewConnect(cfg *config.Config) (*PostgresDb, error) {
	const op = "psql.NewConnect"

	db, err := sql.Open("postgres", cfg.Uri)
	if err != nil {
		return nil, format.Error(op, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, format.Error(op, err)
	}

	log.Success(op, "Connection to postgresSQL is established")
	return &PostgresDb{Driver: db}, nil
}

func (d *PostgresDb) Disconnect() error {
	const op = "psql.PostgresDb.Disconnect"
	err := d.Driver.Close()
	if err != nil {
		return format.Error(op, err)
	}
	err = d.Driver.Ping()
	if err == nil {
		return format.Error(op, errors.New("failed to disconnect"))
	}

	log.Success(op, "Connection to postgresSQL terminated")
	return nil
}

type AuthRepo interface {
	Authentication(u *brzrpc.AuthRequest) (string, error)
	GetAll() ([]*brzrpc.User, error)
	Create(u *brzrpc.User) error
	UpdatePhoto(id, np string) error
	UpdatePassword(id, newPassword string) error
	UpdateEmail(id, email string) error
	UpdateAbout(id, about string) error
	Delete(id string) error
	GetInfo(id string) (*brzrpc.User, error)
}
