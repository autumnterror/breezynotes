package net

import (
	"context"
	"net/http"
	"time"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/labstack/echo/v4"
)

func (e *Echo) ValidateTokenMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const op = "gateway.net.TokenValidate"

			at, err := c.Cookie("access_token")
			if err != nil {
				at = &http.Cookie{Value: "BAD"}
			}
			rt, err := c.Cookie("refresh_token")
			if err != nil {
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
					return next(c)
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
				return next(c)
			}
			return next(c)
		}
	}
}
func ValidateID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if idReq := c.QueryParam("id"); idReq != "" {
				if !uid.Validate(idReq) {
					return c.JSON(http.StatusUnauthorized, domain.Error{Error: "id is not in format"})
				}
			}
			return next(c)
		}
	}
}

func (e *Echo) GetUserId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
			defer done()
			at, err := c.Cookie("access_token")
			if err != nil {
				return next(c)
			}
			u, err := e.authAPI.API.GetIdFromToken(ctx, &brzrpc.Token{Value: at.Value})
			if err != nil || u.GetId() == "" {
				return next(c)
				//return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad access_token"})
			}

			if !uid.Validate(u.GetId()) {
				return c.JSON(http.StatusUnauthorized, domain.Error{Error: "id is not in format"})
			}

			c.Set(domain.IdFromContext, u.GetId())

			return next(c)
		}
	}
}
