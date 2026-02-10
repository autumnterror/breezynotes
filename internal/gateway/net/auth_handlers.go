package net

import (
	"context"
	"net/http"
	"time"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/labstack/echo/v4"
)

// Auth godoc
// @Summary Authorize user
// @Description Authenticates user and returns access/refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param User body domain.AuthRequest true "Login or Email and Password"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/auth [post]
func (e *Echo) Auth(c echo.Context) error {
	const op = "gateway.net.Auth"

	api := e.authAPI.API

	var r domain.AuthRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	res, err := api.Auth(ctx, &brzrpc.AuthRequest{
		Email:    r.Email,
		Login:    r.Login,
		Password: r.Password,
	})
	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    res.GetAccessToken(),
		Path:     "/",
		HttpOnly: true,
		//SameSite: http.SameSiteLaxMode,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(res.ExpAccess, 0).UTC(),
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    res.GetRefreshToken(),
		Path:     "/",
		HttpOnly: true,
		//SameSite: http.SameSiteLaxMode,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(res.ExpRefresh, 0).UTC(),
	})

	return c.JSON(http.StatusOK, domain.AuthResponse{
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
		ExpAccess:    res.ExpAccess,
		ExpRefresh:   res.ExpRefresh,
		Metadata:     domain.UserFromRpc(res.GetMetadata()),
	})
}

// Reg godoc
// @Summary Register new user
// @Description Validates registration data, creates user and returns tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param User body domain.UserRegister true "Reg data"
// @Success 200 {object} domain.Tokens
// @Failure 400 {object} domain.Error
// @Failure 302 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/auth/reg [post]
func (e *Echo) Reg(c echo.Context) error {
	const op = "gateway.net.Reg"

	auth := e.authAPI.API

	var u domain.UserRegister
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}
	if u.Pw1 != u.Pw2 {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "password not same"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	tokens, err := auth.Reg(ctx, &brzrpc.AuthRequest{
		Email:    u.Email,
		Login:    u.Login,
		Password: u.Pw1,
	})
	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    tokens.GetAccessToken(),
		Path:     "/",
		HttpOnly: true,
		//SameSite: http.SameSiteLaxMode,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(tokens.ExpAccess, 0).UTC(),
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.GetRefreshToken(),
		Path:     "/",
		HttpOnly: true,
		//SameSite: http.SameSiteLaxMode,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(tokens.ExpRefresh, 0).UTC(),
	})

	return c.JSON(http.StatusOK, &domain.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpAccess:    tokens.ExpAccess,
		ExpRefresh:   tokens.ExpRefresh,
	})
}

// ValidateToken godoc
// @Summary Validate token (uses cookies)
// @Description Checks access token from cookie, tries to refresh if expired. If 410 (GONE) need to re-auth
// @Tags auth
// @Produce json
// @Success 200 {object} domain.Message
// @Success 201 {object} domain.Token
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/auth/token [get]
func (e *Echo) ValidateToken(c echo.Context) error {
	const op = "gateway.net.TokenValidate"

	at, err := c.Cookie("access_token")
	if err != nil {
		at = &http.Cookie{Value: "BAD"}
	}
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		log.Warn(op, "refresh_token cookie missing", err)
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "refresh_token cookie missing"})
	}

	auth := e.authAPI.API

	ctx, cancel := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer cancel()

	token, err := auth.ValidateTokens(ctx, &brzrpc.Tokens{
		AccessToken:  at.Value,
		RefreshToken: rt.Value,
	})
	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if token != nil {
		if token.Value == "" {
			return c.JSON(http.StatusOK, domain.Message{Message: "tokens valid"})
		}
		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			Value:    token.GetValue(),
			Path:     "/",
			HttpOnly: true,
			//SameSite: http.SameSiteLaxMode,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Unix(token.GetExp(), 0).UTC(),
		})
		return c.JSON(http.StatusCreated, domain.Token{
			Value: token.GetValue(),
			Exp:   token.GetExp(),
		})
	}

	return c.JSON(http.StatusOK, domain.Message{Message: "tokens valid"})
}
