package net

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/autumnterror/breezynotes/pkg/utils/alg"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"

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

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//AUTHORIZE
	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: r.GetId()}); err == nil {
		if n.GetAuthor() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "not author"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}
	//AUTHORIZE

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

	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.GetId()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
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
// @Param id query string true "Note ID"
// @Success 200 {object} views.SWGNoteWithBlocks
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

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//---------------REDIS---------------
	if note, err := e.rdsAPI.API.GetNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: id}); err != nil {
		if strings.Contains(err.Error(), "not found in cache") {
			log.Blue("note not found in cache")
		} else {
			log.Error(op, "REDIS ERROR", err)
		}
	} else {
		if note != nil {
			if note.GetAuthor() != idUser || alg.IsIn(idUser, note.GetEditors()) || alg.IsIn(idUser, note.GetReaders()) {
				return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
			}
			log.Blue("read from cache")
			return c.JSON(http.StatusOK, note)
		}
	}
	//---------------REDIS---------------

	note, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get note error"})
		}

		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "not found"})
		default:
			log.Error(op, "get note error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get note error"})
		}
	}
	if note.GetAuthor() != idUser || alg.IsIn(idUser, note.GetEditors()) || alg.IsIn(idUser, note.GetReaders()) {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
	}

	blocks, err := api.GetAllBlocksInNote(ctx, &brzrpc.Strings{Values: note.Blocks})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get blocks error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get blocks error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get blocks error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get blocks error"})
		}
	}

	if blocks == nil || len(blocks.Items) == 0 {
		blocks = &brzrpc.Blocks{Items: []*brzrpc.Block{}}
	}

	log.Success(op, "")

	nwb := &brzrpc.NoteWithBlocks{
		Id:        note.GetId(),
		Title:     note.GetTitle(),
		Blocks:    blocks.GetItems(),
		Author:    note.GetAuthor(),
		Readers:   note.GetReaders(),
		Editors:   note.GetEditors(),
		CreatedAt: note.GetCreatedAt(),
		UpdatedAt: note.GetUpdatedAt(),
		Tag:       note.GetTag(),
	}

	if _, err := e.rdsAPI.API.SetNoteByUser(ctx, &brzrpc.NoteByUser{UserId: idUser, Note: nwb}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	return c.JSON(http.StatusOK, nwb)
}

// GetAllNotes godoc
// @Summary Get all notes of user
// @Description Returns all notes by user ID
// @Tags note
// @Accept json
// @Produce json
// @Param start query int true  "start > 0"
// @Param end query int true  "end"
// @Success 200 {object} brzrpc.NoteParts
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/all [get]
func (e *Echo) GetAllNotes(c echo.Context) error {
	const op = "gateway.net.GetAllNotes"
	log.Info(op, "")

	api := e.bnAPI.API

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok && idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad idUser from access token"})
	}

	s := c.QueryParam("start")
	if s == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad start"})
	}
	en := c.QueryParam("end")
	if en == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad end"})
	}

	start, err := strconv.Atoi(s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "start must be int"})
	}
	end, err := strconv.Atoi(en)
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "end must be int"})
	}

	if start > end {
		return c.JSON(http.StatusOK, brzrpc.NoteParts{Items: []*brzrpc.NotePart{}})
	}

	if start < 0 {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "start < 0!"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	//---------------REDIS---------------
	if nl, err := e.rdsAPI.API.GetNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	} else {
		if nl != nil {
			if nl.GetItems() != nil {
				if len(nl.GetItems()) != 0 {
					nlPad := make([]*brzrpc.NotePart, end-start)
					for i, n := range nl.GetItems() {
						if i < start {
							continue
						}
						if i >= end {
							break
						}
						nlPad = append(nlPad, n)
					}
					log.Blue("read from cache")
					return c.JSON(http.StatusOK, nlPad)
				}
			}
		}
	}
	//---------------REDIS---------------

	notes, err := api.GetAllNotes(ctx, &brzrpc.UserId{
		UserId: idUser,
	})
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

	if _, err := e.rdsAPI.API.SetNoteListByUser(ctx, &brzrpc.NoteListByUser{UserId: idUser, Items: notes.GetItems()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	items := notes.GetItems()
	if start >= len(items) {
		return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	}
	if end > len(items) {
		end = len(items)
	}
	nlPad := items[start:end]

	log.Success(op, "")

	return c.JSON(http.StatusOK, nlPad)
}

// GetNotesByTag godoc
// @Summary Get notes by tag
// @Description Returns all notes that contain given tag
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true  "Tag ID"
// @Param start query int true  "start > 0"
// @Param end query int true  "end"
// @Success 200 {object} brzrpc.NoteParts
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/by-tag [get]
func (e *Echo) GetNotesByTag(c echo.Context) error {
	const op = "gateway.net.GetNotesByTag"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id"})
	}
	s := c.QueryParam("start")
	if s == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad start"})
	}
	en := c.QueryParam("end")
	if en == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad end"})
	}

	start, err := strconv.Atoi(s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "start must be int"})
	}
	end, err := strconv.Atoi(en)
	if err != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "end must be int"})
	}

	if start > end {
		return c.JSON(http.StatusOK, brzrpc.NoteParts{Items: []*brzrpc.NotePart{}})
	}

	if start < 0 {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "start < 0!"})
	}

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	notes, err := api.GetNotesByTag(ctx, &brzrpc.UserTagId{
		TagId:  id,
		UserId: idUser,
	})
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

	items := notes.GetItems()
	if start >= len(items) {
		return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	}
	if end > len(items) {
		end = len(items)
	}
	nlPad := items[start:end]

	log.Success(op, "")

	return c.JSON(http.StatusOK, nlPad)
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

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var r views.NoteReq
	if err := c.Bind(&r); err != nil {
		log.Error(op, "create note bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	id := uid.New()
	_, err := api.CreateNote(ctx, &brzrpc.Note{
		Id:        id,
		Title:     r.Title,
		CreatedAt: 0,
		UpdatedAt: 0,
		Tag:       nil,
		Author:    idUser,
		Editors:   []string{},
		Readers:   []string{},
		Blocks:    []string{},
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

	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.JSON(http.StatusCreated, brzrpc.Id{Id: id})
}

// AddTagToNote godoc
// @Summary Add tag to note
// @Description Attaches tag to note
// @Tags note
// @Accept json
// @Produce json
// @Param AddTagToNoteRequest body brzrpc.NoteTagId true "Note ID and Tag ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/tag [post]
func (e *Echo) AddTagToNote(c echo.Context) error {
	const op = "gateway.net.AddTagToNote"
	log.Info(op, "")

	api := e.bnAPI.API

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var r brzrpc.NoteTagId
	if err := c.Bind(&r); err != nil {
		log.Error(op, "add tag to note bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: r.GetNoteId()}); err == nil {
		if n.GetAuthor() != idUser || alg.IsIn(idUser, n.GetEditors()) || alg.IsIn(idUser, n.GetReaders()) {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}

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

	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.GetNoteId()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// RmTagFromNote godoc
// @Summary Remove tag from note
// @Description Remove tag from note
// @Tags note
// @Accept json
// @Produce json
// @Param AddTagToNoteRequest body brzrpc.NoteTagId true "Note ID and Tag ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/notes/tag [delete]
func (e *Echo) RmTagFromNote(c echo.Context) error {
	const op = "gateway.net.AddTagToNote"
	log.Info(op, "")

	api := e.bnAPI.API

	idInt := c.Get(IdFromContext)
	idUser, ok := idInt.(string)
	if !ok || idUser == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad id from access token"})
	}

	var r brzrpc.NoteTagId
	if err := c.Bind(&r); err != nil {
		log.Error(op, "add tag to note bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: r.GetNoteId()}); err == nil {
		if n.GetAuthor() != idUser || alg.IsIn(idUser, n.GetEditors()) || alg.IsIn(idUser, n.GetReaders()) {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}

	_, err := api.RemoveTagFromNote(ctx, &r)
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

	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.GetNoteId()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}
