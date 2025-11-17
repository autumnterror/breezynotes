package net

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

// CleanTrash godoc
// @Summary Clean trash
// @Description Deletes all notes from trash for user
// @Tags trash
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/trash [delete]
func (e *Echo) CleanTrash(c echo.Context) error {
	const op = "gateway.net.CleanTrash"
	log.Info(op, "")

	api := e.bnAPI.API

	idInt := c.Get("id")
	id, ok := idInt.(string)
	if !ok && id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	_, err := api.CleanTrash(ctx, &brzrpc.UserId{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "clean trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "clean trash error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "clean trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "clean trash error"})
		}
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// NoteToTrash godoc
// @Summary Move note to trash
// @Description Moves note to trash
// @Tags trash
// @Accept json
// @Produce json
// @Param id query string true "Note ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/trash/to [put]
func (e *Echo) NoteToTrash(c echo.Context) error {
	const op = "gateway.net.NoteToTrash"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	idInt := c.Get("id")
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//AUTHORIZE
	if n, err := api.GetNote(ctx, &brzrpc.Id{Id: id}); err == nil {
		if n.GetAuthor() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "not author"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}
	//AUTHORIZE

	_, err := api.NoteToTrash(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "note to trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "note to trash error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "note to trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "note to trash error"})
		}
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// NoteFromTrash godoc
// @Summary Restore note from trash
// @Description Restores note from trash
// @Tags trash
// @Accept json
// @Produce json
// @Param id query string true "Note ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/trash/from [put]
func (e *Echo) NoteFromTrash(c echo.Context) error {
	const op = "gateway.net.NoteFromTrash"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	idInt := c.Get("id")
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//AUTHORIZE
	if n, err := api.FindNoteInTrash(ctx, &brzrpc.Id{Id: id}); err == nil {
		if n.GetAuthor() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "not author"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}
	//AUTHORIZE

	_, err := api.NoteFromTrash(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "note from trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "note from trash error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "note from trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "note from trash error"})
		}
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// GetNotesFromTrash godoc
// @Summary Get notes from trash
// @Description Returns notes from trash by user ID
// @Tags trash
// @Accept json
// @Produce json
// @Success 200 {object} brzrpc.NoteParts
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/trash [get]
func (e *Echo) GetNotesFromTrash(c echo.Context) error {
	const op = "gateway.net.GetNotesFromTrash"
	log.Info(op, "")

	api := e.bnAPI.API

	idInt := c.Get("id")
	id, ok := idInt.(string)
	if !ok && id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	nts, err := api.GetNotesFromTrash(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get notes from trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get notes from trash error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get notes from trash error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get notes from trash error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, nts)
}

// FindNoteInTrash godoc
// @Summary Find note in trash
// @Description Checks if note exists in trash
// @Tags trash
// @Accept json
// @Produce json
// @Param Id body brzrpc.Id true "Note ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/trash/note/find [post]
//func (e *Echo) FindNoteInTrash(c echo.Context) error {
//	const op = "gateway.net.FindNoteInTrash"
//	log.Info(op, "")
//
//	api := e.bnAPI.API
//
//	var r brzrpc.Id
//	if err := c.Bind(&r); err != nil {
//		log.Error(op, "find note in trash bind", err)
//		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
//	}
//
//	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
//	defer done()
//
//	_, err := api.FindNoteInTrash(ctx, &r)
//	if err != nil {
//		st, ok := status.FromError(err)
//		if !ok {
//			log.Error(op, "find note in trash error", err)
//			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "find note in trash error"})
//		}
//
//		switch st.Code() {
//		case codes.NotFound:
//			log.Warn(op, "note not found", err)
//			return c.JSON(http.StatusNotFound, views.SWGError{Error: "note not found"})
//		case codes.InvalidArgument:
//			log.Warn(op, "bad argument", err)
//			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad argument"})
//		default:
//			log.Error(op, "find note in trash error", err)
//			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "find note in trash error"})
//		}
//	}
//
//	log.Success(op, "")
//
//	return c.NoContent(http.StatusOK)
//}
