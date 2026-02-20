package net

import (
	"context"
	"net/http"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/labstack/echo/v4"

	"github.com/autumnterror/utils_go/pkg/log"
)

// GetRegisteredTypes godoc
// @Summary Get all registered types
// @Description get
// @Tags block
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Router /api/block/types [get]
func (e *Echo) GetRegisteredTypes(c echo.Context) error {
	types, _ := e.bnAPI.API.GetRegisteredBlocks(c.Request().Context(), nil)

	return c.JSON(http.StatusOK, types.Values)
}

// GetBlock godoc
// @Summary GetNote block
// @Description Returns block by ID
// @Tags block
// @Accept json
// @Produce json
// @Param GetBlockRequest body domain.BlockNoteId true "Note ID and Block ID"
// @Success 200 {object} domain.Block
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/block [get]
func (e *Echo) GetBlock(c echo.Context) error {
	const op = "gateway.net.GetBlock"

	api := e.bnAPI.API

	var r domain.BlockNoteId
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	block, err := api.GetBlock(ctx, &brzrpc.NoteBlockUserId{
		NoteId:  r.NoteId,
		BlockId: r.BlockId,
		UserId:  idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.JSON(http.StatusOK, block)
}

// CreateBlock godoc
// @Summary Create block
// @Description Creates block of given type
// @Tags block
// @Accept json
// @Produce json
// @Param CreateBlockRequest body domain.CreateBlockRequest true "Type and data"
// @Success 201 {object} domain.Id
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/block [post]
func (e *Echo) CreateBlock(c echo.Context) error {
	const op = "gateway.net.CreateBlock"

	api := e.bnAPI.API

	var r domain.CreateBlockRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	s, err := structpb.NewStruct(r.Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad data"})
	}

	id, err := api.CreateBlock(ctx, &brzrpc.CreateBlockRequest{
		Type:   r.Type,
		NoteId: r.NoteId,
		Pos:    int32(r.Pos),
		Data:   s,
		UserId: idUser,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.NoteId}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.JSON(http.StatusCreated, id)
}

// OpBlock godoc
// @Summary Operate on block
// @Description Performs operation on block
// @Tags block
// @Accept json
// @Produce json
// @Param OpBlockRequest body domain.OpBlockRequest true "Block ID, operation and data"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/block/op [post]
func (e *Echo) OpBlock(c echo.Context) error {
	const op = "gateway.net.OpBlock"

	api := e.bnAPI.API

	var r domain.OpBlockRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	s, err := structpb.NewStruct(r.Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad data"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err = api.OpBlock(ctx, &brzrpc.OpBlockRequest{
		BlockId: r.BlockId,
		Op:      r.Op,
		Data:    s,
		UserId:  idUser,
		NoteId:  r.NoteId,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.NoteId}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeTypeBlock godoc
// @Summary Change block type
// @Description Changes block type
// @Tags block
// @Accept json
// @Produce json
// @Param ChangeTypeBlockRequest body domain.ChangeTypeBlockRequest true "Block ID and new type"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/block/type [patch]
func (e *Echo) ChangeTypeBlock(c echo.Context) error {
	const op = "gateway.net.ChangeTypeBlock"

	api := e.bnAPI.API

	var r domain.ChangeTypeBlockRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.ChangeTypeBlock(ctx, &brzrpc.ChangeTypeBlockRequest{
		BlockId: r.BlockId,
		NewType: r.NewType,
		UserId:  idUser,
		NoteId:  r.NoteId,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.NoteId}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeBlockOrder godoc
// @Summary Change block order in note
// @Description Changes order of block in note
// @Tags block
// @Accept json
// @Produce json
// @Param ChangeBlockOrderRequest body domain.ChangeBlockOrderRequest true "Note ID and new/old order"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/block/order [patch]
func (e *Echo) ChangeBlockOrder(c echo.Context) error {
	const op = "gateway.net.ChangeBlockOrder"

	api := e.bnAPI.API

	var r domain.ChangeBlockOrderRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.ChangeBlockOrder(ctx, &brzrpc.ChangeBlockOrderRequest{
		NoteId:   r.NoteId,
		OldOrder: int32(r.OldOrder),
		NewOrder: int32(r.NewOrder),
		UserId:   idUser,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.NoteId}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteBlock godoc
// @Summary delete block
// @Description Deletes block by ID
// @Tags block
// @Accept json
// @Produce json
// @Param block_id query string true "Block ID"
// @Param note_id query string true "Note ID"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/block [delete]
func (e *Echo) DeleteBlock(c echo.Context) error {
	const op = "gateway.net.DeleteBlock"

	api := e.bnAPI.API

	blockId := c.QueryParam("block_id")
	if blockId == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "no block id"})
	}
	noteId := c.QueryParam("note_id")
	if noteId == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "no note id"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.DeleteBlock(ctx, &brzrpc.NoteBlockUserId{
		NoteId:  noteId,
		BlockId: blockId,
		UserId:  idUser,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: noteId}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}
