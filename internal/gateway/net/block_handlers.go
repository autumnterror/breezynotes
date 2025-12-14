package net

import (
	"context"
	"net/http"
	"time"

	"github.com/autumnterror/breezynotes/pkg/utils/alg"

	"github.com/labstack/echo/v4"

	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBlock godoc
// @Summary GetNote block
// @Description Returns block by ID
// @Tags block
// @Accept json
// @Produce json
// @Param id query string true "Block ID"
// @Success 200 {object} views.SWGBlock
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/block [post]
func (e *Echo) GetBlock(c echo.Context) error {
	const op = "gateway.net.GetBlock"
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

	block, err := api.GetBlock(ctx, &brzrpc.BlockId{BlockId: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get block error"})
		}

		switch st.Code() {
		case codes.NotFound:
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "not found"})
		default:
			log.Error(op, "get block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get block error"})
		}
	}

	if note, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: block.GetNoteId()}); err == nil {
		if note.GetAuthor() != idUser || !alg.IsIn(idUser, note.GetEditors()) || !alg.IsIn(idUser, note.GetReaders()) {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, block)
}

// CreateBlock godoc
// @Summary Create block
// @Description Creates block of given type
// @Tags block
// @Accept json
// @Produce json
// @Param CreateBlockRequest body views.SWGCreateBlockRequest true "Type and data"
// @Success 201 {object} brzrpc.Id
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/block [post]
func (e *Echo) CreateBlock(c echo.Context) error {
	const op = "gateway.net.CreateBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.CreateBlockRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "create block bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: r.GetNoteId()}); err == nil {
		if n.GetAuthor() != idUser || !alg.IsIn(idUser, n.GetEditors()) {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}

	id, err := api.CreateBlock(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "create block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create block error"})
		}

		switch st.Code() {
		case codes.Unknown:
			log.Warn(op, "unknown block type", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "unknown block type"})
		case codes.NotFound:
			log.Warn(op, "unknown block type", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
		default:
			log.Error(op, "create block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create block error"})
		}
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.GetNoteId()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	log.Success(op, "")

	return c.JSON(http.StatusCreated, id)
}

// OpBlock godoc
// @Summary Operate on block
// @Description Performs operation on block
// @Tags block
// @Accept json
// @Produce json
// @Param OpBlockRequest body views.SWGOpBlockRequest true "Block ID, operation and data"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/block/op [post]
func (e *Echo) OpBlock(c echo.Context) error {
	const op = "gateway.net.OpBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.OpBlockRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "op block bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	noteId := ""

	if block, err := api.GetBlock(ctx, &brzrpc.BlockId{BlockId: r.GetId()}); err == nil {
		if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: block.GetNoteId()}); err == nil {
			if n.GetAuthor() != idUser || !alg.IsIn(idUser, n.GetEditors()) {
				return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
			}
			noteId = n.GetId()
		} else {
			log.Error(op, "get note", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
		}
	}

	_, err := api.OpBlock(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "op block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "op block error"})
		}

		switch st.Code() {
		case codes.Unknown:
			log.Warn(op, "unknown block type", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "unknown block type"})
		default:
			log.Error(op, "op block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "op block error"})
		}
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: noteId}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// ChangeTypeBlock godoc
// @Summary Change block type
// @Description Changes block type
// @Tags block
// @Accept json
// @Produce json
// @Param ChangeTypeBlockRequest body brzrpc.ChangeTypeBlockRequest true "Block ID and new type"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/block/type [patch]
func (e *Echo) ChangeTypeBlock(c echo.Context) error {
	const op = "gateway.net.ChangeTypeBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.ChangeTypeBlockRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "change type block bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if block, err := api.GetBlock(ctx, &brzrpc.BlockId{BlockId: r.GetId()}); err == nil {
		if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: block.GetNoteId()}); err == nil {
			if n.GetAuthor() != idUser || !alg.IsIn(idUser, n.GetEditors()) {
				return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
			}
		} else {
			log.Error(op, "get note", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
		}
	}

	_, err := api.ChangeTypeBlock(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "change type block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change type block error"})
		}

		switch st.Code() {
		case codes.Unknown:
			log.Warn(op, "unknown block type", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "unknown block type"})
		default:
			log.Error(op, "change type block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change type block error"})
		}
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.GetId()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// ChangeBlockOrder godoc
// @Summary Change block order in note
// @Description Changes order of block in note
// @Tags block
// @Accept json
// @Produce json
// @Param ChangeBlockOrderRequest body brzrpc.ChangeBlockOrderRequest true "Note ID and new/old order"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/block/order [patch]
func (e *Echo) ChangeBlockOrder(c echo.Context) error {
	const op = "gateway.net.ChangeBlockOrder"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.ChangeBlockOrderRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "change block order bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: r.GetId()}); err == nil {
		if n.GetAuthor() != idUser || !alg.IsIn(idUser, n.GetEditors()) {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note id"})
	}

	_, err := api.ChangeBlockOrder(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "change block order error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change block order error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "change block order error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "change block order error"})
		}
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.GetId()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// DeleteBlock godoc
// @Summary Delete block
// @Description Deletes block by ID
// @Tags block
// @Accept json
// @Produce json
// @Param block_id query string true "Block ID"
// @Param note_id query string true "Note ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/block [delete]
func (e *Echo) DeleteBlock(c echo.Context) error {
	const op = "gateway.net.DeleteBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	blockId := c.QueryParam("block_id")
	if blockId == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "no block id"})
	}
	noteId := c.QueryParam("note_id")
	if noteId == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "no note id"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if n, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: noteId}); err == nil {
		if n.GetAuthor() != idUser || !alg.IsIn(idUser, n.GetEditors()) {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get note", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad note blockId"})
	}

	_, err := api.DeleteBlock(ctx, &brzrpc.NoteBlockId{
		NoteId:  noteId,
		BlockId: blockId,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "delete block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete block error"})
		}

		switch st.Code() {
		case codes.NotFound:
			log.Warn(op, "block not found", err)
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "block not found"})
		default:
			log.Error(op, "delete block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete block error"})
		}
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: noteId}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}
