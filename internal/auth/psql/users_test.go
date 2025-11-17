package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersOperations(t *testing.T) {
	t.Parallel()

	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	userID := uid.New()
	user := &brzrpc.User{
		Id:       userID,
		Login:    "testlogin",
		Email:    "testemail@example.com",
		About:    "initial about",
		Password: "testpassword",
		Photo:    "testphoto",
	}

	print := func() {
		t.Run("getAll", func(t *testing.T) {
			users, err := repo.GetAll(context.TODO())
			assert.NoError(t, err)
			fmt.Println("🔍 Current Users in DB:")
			for _, u := range users {
				log.Println(format.Struct(u))
			}
		})
	}

	t.Run("create", func(t *testing.T) {
		err := repo.Create(context.TODO(), user)
		assert.NoError(t, err)
	})
	print()

	t.Run("get info", func(t *testing.T) {
		info, err := repo.GetInfo(context.TODO(), user.Id)
		assert.NoError(t, err)
		fmt.Printf("👤 GetInfo: %s", format.Struct(info))
	})

	t.Run("update email", func(t *testing.T) {
		newEmail := "newemail@example.com"
		err := repo.UpdateEmail(context.TODO(), user.Id, newEmail)
		assert.NoError(t, err)
		user.Email = newEmail
		info, err := repo.GetInfo(context.TODO(), user.Id)
		assert.NoError(t, err)
		fmt.Printf("👤 GetInfo after update email: %s", format.Struct(info))
	})

	t.Run("update about", func(t *testing.T) {
		newAbout := "updated about"
		err := repo.UpdateAbout(context.TODO(), user.Id, newAbout)
		assert.NoError(t, err)
		user.About = newAbout
		info, err := repo.GetInfo(context.TODO(), user.Id)
		assert.NoError(t, err)
		fmt.Printf("👤 GetInfo after update about: %s", format.Struct(info))
	})

	t.Run("update password", func(t *testing.T) {
		err := repo.UpdatePassword(context.TODO(), user.Id, "newSecurePassword")
		assert.NoError(t, err)
		log.Println("after update password")
		print()
	})

	t.Run("delete", func(t *testing.T) {
		err := repo.Delete(context.TODO(), user.Id)
		assert.NoError(t, err)
		log.Println("after delete")
		print()
	})
}

func setupTestTx(t *testing.T) (*Driver, *sql.Tx, func()) {
	pdb := MustConnect(config.Test())

	tx, err := pdb.Driver.Begin()
	assert.NoError(t, err)

	return NewDriver(tx), tx, func() {
		assert.NoError(t, tx.Rollback())
		assert.NoError(t, pdb.Disconnect())
	}
}

func TestCreateDuplicateUser(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	uid := uid.New()
	user := &brzrpc.User{
		Id:       uid,
		Login:    "duplicate",
		Email:    "dup@example.com",
		About:    "test",
		Password: "password",
	}

	assert.NoError(t, repo.Create(context.TODO(), user))
	user.Id = uid.New()
	err := repo.Create(context.TODO(), user)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAlreadyExist))
}

func TestUpdateNonExistentUser(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	err := repo.UpdateAbout(context.TODO(), uid.New(), "123")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestDeleteNonExistentUser(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	err := repo.Delete(context.TODO(), uid.New())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestGetInfo_InvalidID(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	_, err := repo.GetInfo(context.TODO(), uid.New())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestAuthLogin(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	uid := uid.New()
	user := &brzrpc.User{
		Id:       uid,
		Login:    "login",
		Email:    "login@example.com",
		About:    "test",
		Password: "password",
	}

	assert.NoError(t, repo.Create(context.TODO(), user))

	_, err := repo.Authentication(context.TODO(), &brzrpc.AuthRequest{
		Email:    "",
		Login:    user.GetLogin(),
		Password: user.GetPassword(),
	})

	assert.NoError(t, err)
}

func TestAuthEmail(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	uid := uid.New()
	user := &brzrpc.User{
		Id:       uid,
		Login:    "login",
		Email:    "login@example.com",
		About:    "test",
		Password: "password",
	}

	assert.NoError(t, repo.Create(context.TODO(), user))

	_, err := repo.Authentication(context.TODO(), &brzrpc.AuthRequest{
		Email:    user.GetEmail(),
		Login:    "",
		Password: user.GetPassword(),
	})
	assert.NoError(t, err)
}

func TestAuthWrongInput1(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	uid := uid.New()
	user := &brzrpc.User{
		Id:       uid,
		Login:    "login",
		Email:    "login@example.com",
		About:    "test",
		Password: "password",
	}

	assert.NoError(t, repo.Create(context.TODO(), user))

	_, err := repo.Authentication(context.TODO(), &brzrpc.AuthRequest{
		Email:    user.GetEmail(),
		Login:    "",
		Password: "",
	})
	assert.True(t, errors.Is(err, ErrWrongInput))
}

func TestAuthWrongInput2(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	uid := uid.New()
	user := &brzrpc.User{
		Id:       uid,
		Login:    "login",
		Email:    "login@example.com",
		About:    "test",
		Password: "password",
	}

	assert.NoError(t, repo.Create(context.TODO(), user))

	_, err := repo.Authentication(context.TODO(), &brzrpc.AuthRequest{
		Email:    "",
		Login:    "",
		Password: "123",
	})
	assert.True(t, errors.Is(err, ErrWrongInput))
}

func TestAuthPwIncorrect(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupTestTx(t)
	defer cleanup()

	uid := uid.New()
	user := &brzrpc.User{
		Id:       uid,
		Login:    "login",
		Email:    "login@example.com",
		About:    "test",
		Password: "password",
	}

	assert.NoError(t, repo.Create(context.TODO(), user))

	_, err := repo.Authentication(context.TODO(), &brzrpc.AuthRequest{
		Email:    "",
		Login:    user.GetLogin(),
		Password: "123",
	})
	assert.True(t, errors.Is(err, ErrPasswordIncorrect))
}
