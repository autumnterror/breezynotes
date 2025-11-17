package net

import (
	"context"
	"github.com/autumnterror/breezynotes/pkg/utils/id"
	"net/http"
	"time"

	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChangeTitleNote godoc
// @Summary Change note title
// @Description Changes title of existing note
// @Tags note
// @Accept json
// @Produce json
// @Param ChangeTitleNoteRequest body brzrpc.ChangeTitleNoteRequest true "Note ID and new title"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/change-title [patch]
func (e *Echo) ChangeTitleNote(c echo.Context) error {
	const op = "gateway.net.ChangeTitleNote"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.ChangeTitleNoteRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "change title bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	_, err := api.ChangeTitleNote(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "change title error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change title error"})
		}

		switch st.Code() {
		case codes.NotFound:
			log.Warn(op, "note not found", err)
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "note not found"})
		default:
			log.Error(op, "change title error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change title error"})
		}
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// GetNote godoc
// @Summary Get note
// @Description Returns note by ID
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true  "Note ID"
// @Success 200 {object} brzrpc.Note
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes [get]
func (e *Echo) GetNote(c echo.Context) error {
	const op = "gateway.net.GetNote"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	note, err := api.GetNote(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get note error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get note error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, note)
}

// GetAllNotes godoc
// @Summary Get all notes of user
// @Description Returns all notes by user ID
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true  "User ID"
// @Success 200 {object} brzrpc.Notes
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/all [get]
func (e *Echo) GetAllNotes(c echo.Context) error {
	const op = "gateway.net.GetAllNotes"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	notes, err := api.GetAllNotes(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get all notes error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get all notes error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get all notes error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get all notes error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, notes)
}

// GetNotesByTag godoc
// @Summary Get notes by tag
// @Description Returns all notes that contain given tag
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true  "Tag ID"
// @Success 200 {object} brzrpc.Notes
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/by-tag [get]
func (e *Echo) GetNotesByTag(c echo.Context) error {
	const op = "gateway.net.GetNotesByTag"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	notes, err := api.GetNotesByTag(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get notes by tag error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get notes by tag error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get notes by tag error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get notes by tag error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, notes)
}

// CreateNote godoc
// @Summary Create note
// @Description Creates new note
// @Tags note
// @Accept json
// @Produce json
// @Param Note body views.NoteReq true "Note info"
// @Success 201 {object} brzrpc.Id
// @Failure 400 {object} views.SWGError
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes [post]
func (e *Echo) CreateNote(c echo.Context) error {
	const op = "gateway.net.CreateNote"
	log.Info(op, "")

	api := e.bnAPI.API

	var r views.NoteReq
	if err := c.Bind(&r); err != nil {
		log.Error(op, "create note bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	t, err := api.GetTag(ctx, &brzrpc.Id{Id: r.TagId})
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad tag"})
	}
	id := id.New()
	_, err = api.CreateNote(ctx, &brzrpc.Note{
		Id:        id,
		Title:     r.Title,
		CreatedAt: 0,
		UpdatedAt: 0,
		Tag:       t,
		Author:    r.Author,
		Editors:   r.Editors,
		Readers:   r.Readers,
		Blocks:    r.Blocks,
		Status:    brzrpc.Statuses(r.Status),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "create note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create note error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "create note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create note error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusCreated, brzrpc.Id{Id: id})
}

// GetAllBlocksInNote godoc
// @Summary Get all blocks in note
// @Description Returns all blocks of given note
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true "Note ID"
// @Success 200 {object} views.SWGBlocks
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/blocks [get]
func (e *Echo) GetAllBlocksInNote(c echo.Context) error {
	const op = "gateway.net.GetAllBlocksInNote"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	blocks, err := api.GetAllBlocksInNote(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get all blocks error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get all blocks error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get all blocks error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get all blocks error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, blocks)
}

// AddTagToNote godoc
// @Summary Add tag to note
// @Description Attaches tag to note
// @Tags note
// @Accept json
// @Produce json
// @Param AddTagToNoteRequest body brzrpc.AddTagToNoteRequest true "Note ID and Tag ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/add-tag [post]
func (e *Echo) AddTagToNote(c echo.Context) error {
	const op = "gateway.net.AddTagToNote"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.AddTagToNoteRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "add tag to note bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	_, err := api.AddTagToNote(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "add tag to note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "add tag to note error"})
		}

		switch st.Code() {
		case codes.InvalidArgument:
			log.Warn(op, "bad argument", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad argument"})
		default:
			log.Error(op, "add tag to note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "add tag to note error"})
		}
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}
