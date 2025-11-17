package net

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
// @Success 200 {object} emptypb.Empty
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/user/{id} [delete]
func (e *Echo) DeleteUser(c echo.Context) error {
	const op = "gateway.net.DeleteUser"
	log.Info(op, "")

	idInt := c.Get("id")
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := auth.DeleteUser(ctx, &brzrpc.UserId{Id: idUser})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "delete failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete failed"})
		}
		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "user does not exist"})
		case codes.Unauthenticated:
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "authentication required"})
		case codes.Internal:
			return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "internal server error"})
		default:
			log.Error(op, "delete failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete failed"})
		}
	}

	return c.JSON(http.StatusOK, &emptypb.Empty{})
}
