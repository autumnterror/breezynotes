package net

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
				return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad access_token"})
			}

			c.Set(domain.IdFromContext, u.GetId())

			return next(c)
		}
	}
}
