package net

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func ValidateID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if idReq := c.QueryParam("id"); idReq != "" {
				if !uid.Validate(idReq) {
					return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "id is not in format"})
				}
			}
			return next(c)
		}
	}
}

func (e *Echo) GetUserId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer done()
			at, err := c.Cookie("access_token")
			if err != nil {
				return next(c)
			}
			u, err := e.authAPI.API.GetIdFromToken(ctx, &brzrpc.Token{Value: at.Value})
			if err != nil || u.GetId() == "" {
				return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad access_token"})
			}

			c.Set("id", u.GetId())

			return next(c)
		}
	}
}
