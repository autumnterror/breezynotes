package net

import (
	"context"
	"net/http"
	"time"

	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
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

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	_, err := api.CleanTrash(ctx, &brzrpc.UserId{UserId: idUser})
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

	if _, err := e.rdsAPI.API.RmNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
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
// @Failure 401 {object} views.SWGError
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

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//AUTHORIZE
	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: id}); err == nil {
		if n.GetAuthor() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "not author"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}
	//AUTHORIZE

	_, err := api.NoteToTrash(ctx, &brzrpc.NoteId{NoteId: id})
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

	if _, err := e.rdsAPI.API.RmNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: id}); err != nil {
		log.Error(op, "REDIS ERROR", err)
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
// @Failure 401 {object} views.SWGError
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

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//AUTHORIZE
	if n, err := api.FindNoteInTrash(ctx, &brzrpc.NoteId{NoteId: id}); err == nil {
		if n.GetAuthor() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "not author"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}
	//AUTHORIZE

	_, err := api.NoteFromTrash(ctx, &brzrpc.NoteId{NoteId: id})
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

	if _, err := e.rdsAPI.API.RmNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: id}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// GetNotesFromTrash godoc
// @Summary GetNote notes from trash
// @Description Returns notes from trash by user ID
// @Tags trash
// @Accept json
// @Produce json
// @Success 200 {object} []brzrpc.NotePart
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/trash [get]
func (e *Echo) GetNotesFromTrash(c echo.Context) error {
	const op = "gateway.net.GetNotesFromTrash"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if ntsT, err := e.rdsAPI.API.GetNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	} else {
		if ntsT != nil {
			if ntsT.GetItems() != nil {
				if len(ntsT.GetItems()) != 0 {
					log.Blue("read from cache")
					return c.JSON(http.StatusOK, ntsT.GetItems())
				}
			}
		}
	}

	nts, err := api.GetNotesFromTrash(ctx, &brzrpc.UserId{UserId: idUser})
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

	if _, err := e.rdsAPI.API.SetNotesFromTrashByUser(ctx, &brzrpc.NoteListByUser{UserId: idUser, Items: nts.GetItems()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	if nts.GetItems() == nil {
		return c.JSON(http.StatusOK, []brzrpc.NotePart{})
	}
	if len(nts.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []brzrpc.NotePart{})
	}

	return c.JSON(http.StatusOK, nts.GetItems())
}
