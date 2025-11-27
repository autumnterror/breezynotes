package net

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/validate"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

// GetUserData godoc
// @Summary get user data from access token
// @Description Checks access token from cookie, and return user data. If 401 call ValidateToken
// @Tags user
// @Produce json
// @Success 200 {object} brzrpc.User
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user/data [get]
func (e *Echo) GetUserData(c echo.Context) error {
	const op = "gateway.net.GetUserData"
	log.Info(op, "")

	at, err := c.Cookie("access_token")
	if err != nil {
		log.Warn(op, "access_token cookie missing", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "access_token cookie missing"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	u, err := auth.GetUserDataFromToken(ctx, &brzrpc.Token{Value: at.Value})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get failed"})
		}
		switch st.Code() {
		case codes.Unauthenticated:
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "token expired"})
		case codes.InvalidArgument:
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "token is bad, or not access"})
		case codes.NotFound:
			return c.JSON(http.StatusGone, views.SWGError{Error: "user does not exist"})
		default:
			log.Error(op, "get failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get failed"})
		}
	}
	u.Password = ""
	return c.JSON(http.StatusOK, u)
}

// DeleteUser godoc
// @Summary delete user account
// @Description Deletes user account by ID. Requires authentication.
// @Tags user
// @Produce json
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user [delete]
func (e *Echo) DeleteUser(c echo.Context) error {
	const op = "gateway.net.DeleteUser"
	log.Info(op, "")

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	_, err := auth.DeleteUser(ctx, &brzrpc.UserId{UserId: idUser})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "delete failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete failed"})
		}
		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "user does not exist"})
		case codes.Internal:
			return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "internal server error"})
		default:
			log.Error(op, "delete failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete failed"})
		}
	}

	return c.NoContent(http.StatusOK)
}

// UpdateAbout godoc
// @Summary update user about
// @Description Updates user "about" field. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body views.UpdateAboutRequest true "new about text"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 500 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user/about [patch]
func (e *Echo) UpdateAbout(c echo.Context) error {
	const op = "gateway.net.UpdateAbout"
	log.Info(op, "")

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var req views.UpdateAboutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "invalid request body"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	_, err := auth.UpdateAbout(ctx, &brzrpc.UpdateAboutRequest{
		Id:       idUser,
		NewAbout: req.NewAbout,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "update about failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update about failed"})
		}

		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "user does not exist"})
		case codes.Internal:
			return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "internal server error"})
		default:
			log.Error(op, "update about failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update about failed"})
		}
	}

	return c.NoContent(http.StatusOK)
}

// UpdateEmail godoc
// @Summary update user email
// @Description Updates user email. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body views.UpdateEmailRequest true "new email"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 500 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user/email [patch]
func (e *Echo) UpdateEmail(c echo.Context) error {
	const op = "gateway.net.UpdateEmail"
	log.Info(op, "")

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var req brzrpc.UpdateEmailRequest
	if err := c.Bind(&req); err != nil || req.NewEmail == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "invalid request body"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	_, err := auth.UpdateEmail(ctx, &brzrpc.UpdateEmailRequest{
		Id:       idUser,
		NewEmail: req.NewEmail,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "update email failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update email failed"})
		}

		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "user does not exist"})
		case codes.Internal:
			return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "internal server error"})
		default:
			log.Error(op, "update email failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update email failed"})
		}
	}

	return c.NoContent(http.StatusOK)
}

// UpdatePhoto godoc
// @Summary update user photo
// @Description Updates user photo. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body views.UpdatePhotoRequest true "new photo"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 500 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user/photo [patch]
func (e *Echo) UpdatePhoto(c echo.Context) error {
	const op = "gateway.net.UpdatePhoto"
	log.Info(op, "")

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var req views.UpdatePhotoRequest
	if err := c.Bind(&req); err != nil || req.NewPhoto == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "invalid request body"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	_, err := auth.UpdatePhoto(ctx, &brzrpc.UpdatePhotoRequest{
		Id:       idUser,
		NewPhoto: req.NewPhoto,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "update photo failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update photo failed"})
		}

		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "user does not exist"})

		case codes.Internal:
			return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "internal server error"})
		default:
			log.Error(op, "update photo failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update photo failed"})
		}
	}

	return c.NoContent(http.StatusOK)
}

// ChangePassword godoc
// @Summary change user password
// @Description Changes user password. Requires authentication.
// @Tags user
// @Accept json
// @Produce json
// @Param request body views.ChangePasswordRequest true "new password"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 500 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user/pw [patch]
func (e *Echo) ChangePassword(c echo.Context) error {
	const op = "gateway.net.ChangePassword"
	log.Info(op, "")

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var req views.ChangePasswordRequest
	if err := c.Bind(&req); err != nil || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "invalid request body"})
	}

	if req.NewPassword != req.NewPassword2 {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "password not same"})
	}
	if !validate.Password(req.NewPassword) {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "password not in policy"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	_, err := auth.ChangePasswd(ctx, &brzrpc.ChangePasswordRequest{
		Id:          idUser,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "change password failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change password failed"})
		}

		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "user does not exist"})
		case codes.Internal:
			return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "internal server error"})
		default:
			log.Error(op, "change password failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change password failed"})
		}
	}

	return c.NoContent(http.StatusOK)
}
