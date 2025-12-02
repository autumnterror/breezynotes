package net

import (
	"context"
	"net/http"
	"time"

	"github.com/autumnterror/breezynotes/pkg/log"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/utils/uid"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateTag godoc
// @Summary Create tag
// @Description Creates new tag
// @Tags tag
// @Accept json
// @Produce json
// @Param Tag body views.TagReq true "Tag data"
// @Success 201 {object} brzrpc.Id
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/tag [post]
func (e *Echo) CreateTag(c echo.Context) error {
	const op = "gateway.net.CreateTag"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	var r views.TagReq
	if err := c.Bind(&r); err != nil {
		log.Error(op, "create tag bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	newId := uid.New()

	_, err := api.CreateTag(ctx, &brzrpc.Tag{
		Id:     newId,
		Title:  r.Title,
		Color:  r.Color,
		Emoji:  r.Emoji,
		UserId: idUser,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "create tag error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create tag error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "create tag error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "create tag error"})
		}
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.JSON(http.StatusCreated, brzrpc.Id{Id: newId})
}

// UpdateTagTitle godoc
// @Summary Update tag title
// @Description Updates title of tag
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagTitleRequest body brzrpc.UpdateTagTitleRequest true "Tag ID and title"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/tag/title [patch]
func (e *Echo) UpdateTagTitle(c echo.Context) error {
	const op = "gateway.net.UpdateTagTitle"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	var r brzrpc.UpdateTagTitleRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "update tag title bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if t, err := api.GetTag(ctx, &brzrpc.TagId{TagId: r.GetId()}); err == nil {
		if t.GetUserId() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get tag", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad tag id"})
	}

	_, err := api.UpdateTagTitle(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "update tag title error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update tag title error"})
		}

		switch st.Code() {
		case codes.NotFound:
			log.Warn(op, "tag not found", err)
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "tag not found"})
		default:
			log.Error(op, "update tag title error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update tag title error"})
		}
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// UpdateTagColor godoc
// @Summary Update tag color
// @Description Updates color of tag
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagColorRequest body brzrpc.UpdateTagColorRequest true "Tag ID and color"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/tag/color [patch]
func (e *Echo) UpdateTagColor(c echo.Context) error {
	const op = "gateway.net.UpdateTagColor"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	var r brzrpc.UpdateTagColorRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "update tag color bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if t, err := api.GetTag(ctx, &brzrpc.TagId{TagId: r.GetId()}); err == nil {
		if t.GetUserId() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get tag", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad tag id"})
	}

	_, err := api.UpdateTagColor(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "update tag color error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update tag color error"})
		}

		switch st.Code() {
		case codes.NotFound:
			log.Warn(op, "tag not found", err)
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "tag not found"})
		default:
			log.Error(op, "update tag color error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update tag color error"})
		}
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// UpdateTagEmoji godoc
// @Summary Update tag emoji
// @Description Updates emoji of tag
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagEmojiRequest body brzrpc.UpdateTagEmojiRequest true "Tag ID and emoji"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 404 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/tag/emoji [patch]
func (e *Echo) UpdateTagEmoji(c echo.Context) error {
	const op = "gateway.net.UpdateTagEmoji"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	var r brzrpc.UpdateTagEmojiRequest
	if err := c.Bind(&r); err != nil {
		log.Error(op, "update tag emoji bind", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if t, err := api.GetTag(ctx, &brzrpc.TagId{TagId: r.GetId()}); err == nil {
		if t.GetUserId() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get tag", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad tag id"})
	}

	_, err := api.UpdateTagEmoji(ctx, &r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "update tag emoji error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update tag emoji error"})
		}

		switch st.Code() {
		case codes.NotFound:
			log.Warn(op, "tag not found", err)
			return c.JSON(http.StatusNotFound, views.SWGError{Error: "tag not found"})
		default:
			log.Error(op, "update tag emoji error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "update tag emoji error"})
		}
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// DeleteTag godoc
// @Summary Delete tag
// @Description Deletes tag by ID
// @Tags tag
// @Accept json
// @Produce json
// @Param id query string true "Tag ID"
// @Success 200
// @Failure 400 {object} views.SWGError
// @Failure 401 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/tag [delete]
func (e *Echo) DeleteTag(c echo.Context) error {
	const op = "gateway.net.DeleteTag"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if t, err := api.GetTag(ctx, &brzrpc.TagId{TagId: id}); err == nil {
		if t.GetUserId() != idUser {
			return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "user dont have permission"})
		}
	} else {
		log.Error(op, "get tag", err)
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad tag id"})
	}

	_, err := api.DeleteTag(ctx, &brzrpc.TagId{TagId: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "delete tag error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete tag error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "delete tag error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "delete tag error"})
		}
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}

	log.Success(op, "")

	return c.NoContent(http.StatusOK)
}

// GetTagsByUser godoc
// @Summary Get tags by user
// @Description Returns all tags for user
// @Tags tag
// @Accept json
// @Produce json
// @Success 200 {object} []brzrpc.Tag
// @Failure 400 {object} views.SWGError
// @Failure 502 {object} views.SWGError
// @Router /api/tag/by-user [get]
func (e *Echo) GetTagsByUser(c echo.Context) error {
	const op = "gateway.net.GetTagsByUser"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, views.SWGError{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer done()

	if tgs, err := e.rdsAPI.API.GetTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	} else {
		if tgs != nil {
			if len(tgs.GetItems()) != 0 {
				log.Blue("read from cache")
				return c.JSON(http.StatusOK, tgs.GetItems())
			}
		}
	}

	tags, err := api.GetTagsByUser(ctx, &brzrpc.UserId{UserId: idUser})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "get tags by user error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get tags by user error"})
		}

		switch st.Code() {
		default:
			log.Error(op, "get tags by user error", err)
			return c.JSON(http.StatusBadGateway, views.SWGError{Error: "get tags by user error"})
		}
	}

	if _, err := e.rdsAPI.API.SetTagsByUser(ctx, &brzrpc.TagsByUser{UserId: idUser, Items: tags.GetItems()}); err != nil {
		log.Error(op, "REDIS ERROR", err)
	}
	if len(tags.Items) == 0 {
		tags.Items = []*brzrpc.Tag{}
	}
	log.Success(op, "")

	return c.JSON(http.StatusOK, tags.GetItems())
}
