package net

import (
	"context"
	brzrpc2 "github.com/autumnterror/breezynotes/api/proto/gen"
	"net/http"
	"time"

	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
	"github.com/autumnterror/breezynotes/pkg/utils/validate"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth godoc
// @Summary Authorize user
// @Description Authenticates user and returns access/refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param User body brzrpc.AuthRequest true "Login or Email and Password"
// @Success 200 {object} brzrpc.Tokens
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/auth [post]
func (e *Echo) Auth(c echo.Context) error {
	const op = "gateway.net.Auth"
	log.Info(op, "")
	api := e.authAPI.API

	var r brzrpc2.AuthRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "auth bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	id, err := api.Auth(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "authentication error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "authentication error"})
		}

		switch st.Code() {
		case codes.Unauthenticated:
			log.Warn(op, "wrong login or password", err)
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "wrong login or password"})
		case codes.InvalidArgument:
			log.Warn(op, "bad argument", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad argument"})
		default:
			log.Error(op, "auth req", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "authentication error"})
		}
	}

	tokens, err := api.GenerateTokens(ctx, id)
	if err != nil {
		log.Error(op, "token generation error", err)
		return c.JSON(http.StatusBadGateway, views.SWGError{Error: "token generation error"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    tokens.GetAccessToken(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(tokens.ExpAccess, 0).UTC(),
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.GetRefreshToken(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(tokens.ExpRefresh, 0).UTC(),
	})

	log.Success(op, "")

	return c.JSON(http.StatusOK, brzrpc2.Tokens{
		AccessToken:  tokens.GetAccessToken(),
		RefreshToken: tokens.GetRefreshToken(),
	})
}

// Reg godoc
// @Summary Register new user
// @Description Validates registration data, creates user and returns tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param User body views.UserRegister true "Reg data"
// @Success 200 {object} brzrpc.Tokens
// @Failure 400 {object} views.SWGError
// @Failure 302 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/auth/reg [post]
func (e *Echo) Reg(c echo.Context) error {
	const op = "gateway.net.Reg"
	log.Info(op, "")

	auth := e.authAPI.API

	var u views.UserRegister
	if err := c.Bind(&u); err != nil {
		log.Warn(op, "bad JSON", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}
	if u.Pw1 != u.Pw2 {
		log.Warn(op, "password not same", nil)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "password not same"})
	}

	if !validate.Password(u.Pw1) {
		log.Warn(op, "password not in policy", nil)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "password not in policy"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	id := uid.New()
	_, err := auth.CreateUser(ctx, &brzrpc2.User{
		Id:       id,
		Login:    u.Login,
		Email:    u.Email,
		About:    "Write me!",
		Photo:    "images/default.png",
		Password: u.Pw1,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "user creation failed", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "user creation failed"})
		}

		switch st.Code() {
		case codes.AlreadyExists:
			log.Warn(op, "user creation failed", err)
			return c.JSON(http.StatusFound, views.SWGError{Error: "user already exist"})
		default:
			log.Error(op, "", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "user creation failed"})
		}
	}

	tokens, err := auth.GenerateTokens(ctx, &brzrpc2.UserId{UserId: id})
	if err != nil {
		log.Error(op, "token generation error", err)
		return c.JSON(http.StatusBadGateway, views.SWGError{Error: "token generation error"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    tokens.GetAccessToken(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(tokens.ExpAccess, 0).UTC(),
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.GetRefreshToken(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(tokens.ExpRefresh, 0).UTC(),
	})

	log.Success(op, "")

	return c.JSON(http.StatusOK, &brzrpc2.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// ValidateToken godoc
// @Summary Validate token (uses cookies)
// @Description Checks access token from cookie, tries to refresh if expired. If 410 (GONE) need to re-auth
// @Tags auth
// @Produce json
// @Success 200 {object} views.SWGMessage
// @Success 201 {object} brzrpc.Token
// @Failure 400 {object} views.SWGError
// @Failure 410 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/auth/token [get]
func (e *Echo) ValidateToken(c echo.Context) error {
	const op = "gateway.net.TokenValidate"
	log.Info(op, "")

	at, err := c.Cookie("access_token")
	if err != nil {
		at = &http.Cookie{Value: "BAD"}
	}
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		log.Warn(op, "refresh_token cookie missing", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "refresh_token cookie missing"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	if _, err := auth.CheckToken(ctx, &brzrpc2.Token{Value: at.Value}); err != nil {
		newAt, err := auth.Refresh(ctx, &brzrpc2.Token{Value: rt.Value})
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				log.Error(op, "refresh failed", err)
				return c.JSON(http.StatusBadGateway, views.SWGError{Error: "refresh failed"})
			}

			switch st.Code() {
			case codes.Unauthenticated:
				log.Warn(op, "refresh token expired. Terminate authorization", err)
				return c.JSON(http.StatusGone, views.SWGError{Error: "refresh token expired. Terminate authorization"})

			case codes.InvalidArgument:
				log.Warn(op, "invalid refresh token", err)
				return c.JSON(http.StatusBadRequest, views.SWGError{Error: "invalid refresh token"})

			default:
				log.Error(op, "refresh failed", err)
				return c.JSON(http.StatusBadGateway, views.SWGError{Error: "refresh failed"})
			}
		}
		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			Value:    newAt.GetValue(),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Unix(newAt.GetExp(), 0).UTC(),
		})
		return c.JSON(http.StatusCreated, newAt)
	}
	log.Success(op, "")
	return c.JSON(http.StatusOK, views.SWGMessage{Message: "tokens valid"})
}
