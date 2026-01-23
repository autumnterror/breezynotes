package net

import (
	"errors"
	"fmt"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/auth"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/blocknote"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/redis"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/utils_go/pkg/utils/format"
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
		AllowHeaders:     []string{echo.HeaderContentType},
		AllowCredentials: true,
	}))

	//e.echo.Use(middleware.Logger(), middleware.Recover())
	//e.echo.Static("/", "./example/html")

	api := e.echo.Group("/api", ValidateID(), e.GetUserId())
	{
		api.GET("/healthz", e.Healthz)

		auth := api.Group("/auth")
		{
			auth.GET("/token", e.ValidateToken)
			auth.POST("", e.Auth)
			auth.POST("/reg", e.Reg)
		}
		user := api.Group("/user")
		{
			user.GET("/data", e.GetUserData)
			user.DELETE("", e.DeleteUser)
			user.PATCH("/about", e.UpdateAbout)
			user.PATCH("/email", e.UpdateEmail)
			user.PATCH("/photo", e.UpdatePhoto)
			user.PATCH("/pw", e.ChangePassword)
		}

		notes := api.Group("/note")
		{
			notes.GET("", e.GetNote)
			notes.GET("/find", e.Search)
			notes.POST("", e.CreateNote)

			notes.GET("/all", e.GetAllNotes)
			notes.GET("/by-tag", e.GetNotesByTag)
			notes.PATCH("/title", e.ChangeTitleNote)

			notes.POST("/tag", e.AddTagToNote)
			notes.DELETE("/tag", e.RmTagFromNote)
		}

		blocks := api.Group("/block")
		{
			blocks.GET("", e.GetBlock)
			blocks.POST("", e.CreateBlock)
			blocks.DELETE("", e.DeleteBlock)

			blocks.POST("/op", e.OpBlock)
			blocks.PATCH("/type", e.ChangeTypeBlock)
			blocks.PATCH("/order", e.ChangeBlockOrder)
		}

		trash := api.Group("/trash")
		{
			trash.DELETE("", e.CleanTrash)
			trash.PUT("/to", e.NoteToTrash)
			trash.PUT("/from", e.NoteFromTrash)
			//trash.POST("/note/find", e.FindNoteInTrash)
			trash.GET("", e.GetNotesFromTrash)
		}

		tags := api.Group("/tag")
		{
			tags.GET("/by-user", e.GetTagsByUser)

			tags.POST("", e.CreateTag)

			tags.PATCH("/title", e.UpdateTagTitle)
			tags.PATCH("/color", e.UpdateTagColor)
			tags.PATCH("/emoji", e.UpdateTagEmoji)

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
