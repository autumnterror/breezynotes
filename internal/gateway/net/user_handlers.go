package net

import (
	"context"
	"net/http"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/labstack/echo/v4"
)

// GetUserData godoc
// @Summary get user data from access token
// @Description Checks access token from cookie, and return user data. If 401 call ValidateToken
// @Tags user
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/user/data [get]
func (e *Echo) GetUserData(c echo.Context) error {
	const op = "gateway.net.GetUserData"

	at, err := c.Cookie("access_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "access_token cookie missing"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	u, err := auth.GetUserDataFromToken(ctx, &brzrpc.Token{Value: at.Value})

	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	u.Password = ""

	return c.JSON(http.StatusOK, domain.UserFromRpc(u))
}

// DeleteUser godoc
// @Summary delete user account
// @Description Deletes user account by ID. Requires authentication.
// @Tags user
// @Produce json
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/user [delete]
func (e *Echo) DeleteUser(c echo.Context) error {
	const op = "gateway.net.DeleteUser"

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	_, err := e.bnAPI.API.NotesToTrash(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}
	_, err = e.bnAPI.API.CleanTrash(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes = bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	_, err = e.bnAPI.API.DeleteTags(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes = bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	_, err = auth.DeleteUser(ctx, &brzrpc.UserId{UserId: idUser})

	code, errRes = authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateAbout godoc
// @Summary update user about
// @Description Updates user "about" field. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body domain.UpdateAboutRequest true "new about text"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/user/about [patch]
func (e *Echo) UpdateAbout(c echo.Context) error {
	const op = "gateway.net.UpdateAbout"

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var req domain.UpdateAboutRequest
	if err := c.Bind(&req); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	_, err := auth.UpdateAbout(ctx, &brzrpc.UpdateAboutRequest{
		Id:       idUser,
		NewAbout: req.NewAbout,
	})

	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateEmail godoc
// @Summary update user email
// @Description Updates user email. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body domain.UpdateEmailRequest true "new email"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/user/email [patch]
func (e *Echo) UpdateEmail(c echo.Context) error {
	const op = "gateway.net.UpdateEmail"

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var req domain.UpdateEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	_, err := auth.UpdateEmail(ctx, &brzrpc.UpdateEmailRequest{
		Id:       idUser,
		NewEmail: req.NewEmail,
	})

	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdatePhoto godoc
// @Summary update user photo
// @Description Updates user photo. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body domain.UpdatePhotoRequest true "new photo"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/user/photo [patch]
func (e *Echo) UpdatePhoto(c echo.Context) error {
	const op = "gateway.net.UpdatePhoto"

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var req domain.UpdatePhotoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	_, err := auth.UpdatePhoto(ctx, &brzrpc.UpdatePhotoRequest{
		Id:       idUser,
		NewPhoto: req.NewPhoto,
	})
	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangePassword godoc
// @Summary change user password
// @Description Changes user password. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body domain.ChangePasswordRequest true "new password"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/user/pw [patch]
func (e *Echo) ChangePassword(c echo.Context) error {
	const op = "gateway.net.ChangePassword"

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var req domain.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	if req.NewPassword != req.NewPassword2 {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "password not same"})
	}

	api := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	_, err := api.ChangePasswd(ctx, &brzrpc.ChangePasswordRequest{
		Id:          idUser,
		NewPassword: req.NewPassword,
		OldPassword: req.OldPassword,
	})

	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.NoContent(http.StatusNoContent)
}
