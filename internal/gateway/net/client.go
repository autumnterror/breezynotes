package net

import (
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/auth"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/blocknote"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/redis"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"strings"
)

type Echo struct {
	echo    *echo.Echo
	cfg     *config.Config
	authAPI *auth.Client
	bnAPI   *blocknote.Client
	rdsAPI  *redis.Client
}

func New(
	cfg *config.Config,
	authAPI *auth.Client,
	bnAPI *blocknote.Client,
	rdsAPI *redis.Client,
) *Echo {
	e := &Echo{
		echo:    echo.New(),
		cfg:     cfg,
		authAPI: authAPI,
		bnAPI:   bnAPI,
		rdsAPI:  rdsAPI,
	}

	e.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	e.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins:     []string{"http://localhost:5500", "http://127.0.0.1:5500", "http://localhost:8080"},
		AllowOriginFunc: func(origin string) (bool, error) {
			return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1"), nil
		},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	//e.echo.Use(middleware.Logger(), middleware.Recover())

	api := e.echo.Group("/api", ValidateID(), e.GetUserId())
	{
		api.GET("/health", e.Healthz)

		auth := api.Group("/auth")
		{
			auth.GET("/token", e.ValidateToken)
			auth.POST("", e.Auth)
			auth.POST("/reg", e.Reg)
		}
		user := api.Group("/user")
		{
			user.GET("/data", e.GetUserData)
		}

		notes := api.Group("/notes")
		{
			notes.GET("", e.GetNote)
			notes.POST("", e.CreateNote)

			notes.GET("/all", e.GetAllNotes)
			notes.GET("/by-tag", e.GetNotesByTag)
			notes.GET("/blocks", e.GetAllBlocksInNote)
			notes.PATCH("/change-title", e.ChangeTitleNote)

			notes.POST("/add-tag", e.AddTagToNote)
		}

		blocks := api.Group("/blocks")
		{
			blocks.GET("", e.GetBlock)
			blocks.POST("", e.CreateBlock)
			blocks.DELETE("", e.DeleteBlock)

			blocks.GET("/as-first", e.GetBlockAsFirst)
			blocks.POST("/op", e.OpBlock)
			blocks.PATCH("/change-type", e.ChangeTypeBlock)
			blocks.PATCH("/change-order", e.ChangeBlockOrder)
		}

		trash := api.Group("/trash")
		{
			trash.DELETE("", e.CleanTrash)
			trash.PUT("/to", e.NoteToTrash)
			trash.PUT("/from", e.NoteFromTrash)
			//trash.POST("/note/find", e.FindNoteInTrash)
			trash.GET("", e.GetNotesFromTrash)
		}

		tags := api.Group("/tags")
		{
			tags.GET("/by-user", e.GetTagsByUser)

			tags.POST("", e.CreateTag)

			tags.PUT("/title", e.UpdateTagTitle)
			tags.PUT("/color", e.UpdateTagColor)
			tags.PUT("/emoji", e.UpdateTagEmoji)

			tags.DELETE("", e.DeleteTag)
		}
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
