package net

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/views"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBlockAsFirst godoc
// @Summary Get block as first string
// @Description Returns block representation as string
// @Tags block
// @Accept json
// @Produce json
// @Param id query string true "Block ID"
// @Success 200 {object} brzrpc.StringResponse
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/blocks/as-first [get]
func (e *Echo) GetBlockAsFirst(c echo.Context) error {
	const op = "gateway.net.GetBlockAsFirst"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	resp, err := api.GetBlockAsFirst(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get block as first error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get block as first error"})
		}

		switch st.Code() {
		case codes.Unknown:
			log.Warn(op, "unknown block type", err)
			return c.JSON(http.StatusBadRequest, views.SWGError{Error: "unknown block type"})
		default:
			log.Error(op, "get block as first error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get block as first error"})
		}
	}

	log.Success(op, "")

	return c.JSON(http.StatusOK, resp)
}

// GetBlock godoc
// @Summary Get block
// @Description Returns block by ID
// @Tags block
// @Accept json
// @Produce json
// @Param id query string true "Block ID"
// @Success 200 {object} views.SWGBlock
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/blocks [post]
func (e *Echo) GetBlock(c echo.Context) error {
	const op = "gateway.net.GetBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	block, err := api.GetBlock(ctx, &brzrpc.Id{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get block error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get block error"})
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
// @Failure 502 {object} views.SWGError
// @Router /api/blocks [post]
func (e *Echo) CreateBlock(c echo.Context) error {
	const op = "gateway.net.CreateBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.CreateBlockRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "create block bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

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
		default:
			log.Error(op, "create block error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create block error"})
		}
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
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/blocks/op [post]
func (e *Echo) OpBlock(c echo.Context) error {
	const op = "gateway.net.OpBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.OpBlockRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "op block bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

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
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/blocks/change-type [patch]
func (e *Echo) ChangeTypeBlock(c echo.Context) error {
	const op = "gateway.net.ChangeTypeBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.ChangeTypeBlockRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "change type block bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

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
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/blocks/change-order [patch]
func (e *Echo) ChangeBlockOrder(c echo.Context) error {
	const op = "gateway.net.ChangeBlockOrder"
	log.Info(op, "")

	api := e.bnAPI.API

	var r brzrpc.ChangeBlockOrderRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "change block order bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

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

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// DeleteBlock godoc
// @Summary Delete block
// @Description Deletes block by ID
// @Tags block
// @Accept json
// @Produce json
// @Param id query string true "Block ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/blocks [delete]
func (e *Echo) DeleteBlock(c echo.Context) error {
	const op = "gateway.net.DeleteBlock"
	log.Info(op, "")

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	_, err := api.DeleteBlock(ctx, &brzrpc.Id{Id: id})
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

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}
