package net

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// Healthz godoc
// @Summary check health of gateway
// @Description
// @Tags healthz
// @Produce json
// @Success 200 {object} views.SWGMessage
// @Router /api/health [get]
func (e *Echo) Healthz(c echo.Context) error {
	const op = "gateway.net.Healthz"
	log.Info(op, "")

	ctx, done := context.WithTimeout(c.Request().Context(), time.Second)
	defer done()
	_, err := e.bnAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.SWGMessage{Message: "bad blocknote"})
	}
	_, err = e.rdsAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.SWGMessage{Message: "bad redis"})
	}
	_, err = e.authAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.SWGMessage{Message: "bad auth"})
	}

	return c.JSON(http.StatusOK, views.SWGMessage{Message: "HEALTHZ"})
}
