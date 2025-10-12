package net

import (
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type Echo struct {
	echo *echo.Echo
	cfg  *config.Config
}

func New(cfg *config.Config) *Echo {
	e := &Echo{
		echo: echo.New(),
		cfg:  cfg,
	}

	e.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	e.echo.Use(middleware.Logger(), middleware.Recover())

	health := e.echo.Group("/api")
	{
		health.GET("/health", e.Healthz)
	}

	return e
}

func (e *Echo) MustRun() {
	const op = "net.Run"

	if err := e.echo.Start(fmt.Sprintf(":%d", e.cfg.Port)); err != nil && !errors.Is(http.ErrServerClosed, err) {
		e.echo.Logger.Fatal(format.Error(op, err))
	}
}

func (e *Echo) Stop() error {
	const op = "net.Stop"

	if err := e.echo.Close(); err != nil {
		return format.Error(op, err)
	}
	return nil
}
