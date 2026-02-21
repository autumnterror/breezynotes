package net

import (
	"context"
	"net/http"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
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
// @Param Tag body domain.CreateTagRequest true "Tag data"
// @Success 201 {object} domain.Id
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag [post]
func (e *Echo) CreateTag(c echo.Context) error {
	const op = "gateway.net.CreateTag"
	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.CreateTagRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	newId := uid.New()

	_, err := api.CreateTag(ctx, &brzrpc.Tag{
		Id:     newId,
		Title:  r.Title,
		Color:  r.Color,
		Emoji:  r.Emoji,
		UserId: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.JSON(http.StatusCreated, domain.Id{Id: newId})
}

// UpdateTagTitle godoc
// @Summary Update tag title
// @Description Updates title of tag
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagTitleRequest body domain.UpdateTagTitleRequest true "Tag ID and title"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag/title [patch]
func (e *Echo) UpdateTagTitle(c echo.Context) error {
	const op = "gateway.net.UpdateTagTitle"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.UpdateTagTitleRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.UpdateTagTitle(ctx, &brzrpc.UpdateTagTitleRequest{
		IdTag:  r.IdTag,
		Title:  r.Title,
		IdUser: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

// UpdateTagColor godoc
// @Summary Update tag color
// @Description Updates color of tag
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagColorRequest body domain.UpdateTagColorRequest true "Tag ID and color"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag/color [patch]
func (e *Echo) UpdateTagColor(c echo.Context) error {
	const op = "gateway.net.UpdateTagColor"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.UpdateTagColorRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.UpdateTagColor(ctx, &brzrpc.UpdateTagColorRequest{
		IdTag:  r.IdTag,
		Color:  r.Color,
		IdUser: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

// UpdateTagEmoji godoc
// @Summary Update tag emoji
// @Description Updates emoji of tag
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagEmojiRequest body domain.UpdateTagEmojiRequest true "Tag ID and emoji"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag/emoji [patch]
func (e *Echo) UpdateTagEmoji(c echo.Context) error {
	const op = "gateway.net.UpdateTagEmoji"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.UpdateTagEmojiRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.UpdateTagEmoji(ctx, &brzrpc.UpdateTagEmojiRequest{
		IdTag:  r.IdTag,
		Emoji:  r.Emoji,
		IdUser: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

// UpdatePinnedEmoji godoc
// @Summary Update pinned emoji
// @Tags tag
// @Accept json
// @Produce json
// @Param UpdateTagEmojiRequest body domain.UpdatePinnedEmojiRequest true "Tag ID and emoji"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag/pinned [patch]
func (e *Echo) UpdatePinnedEmoji(c echo.Context) error {
	const op = "gateway.net.UpdatePinnedEmoji"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.UpdatePinnedEmojiRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.UpdateTagPinned(ctx, &brzrpc.UserTagId{
		TagId:  r.IdTag,
		UserId: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

// DeleteTag godoc
// @Summary delete tag
// @Description Deletes tag by ID
// @Tags tag
// @Accept json
// @Produce json
// @Param id query string true "Tag ID"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag [delete]
func (e *Echo) DeleteTag(c echo.Context) error {
	const op = "gateway.net.DeleteTag"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad param"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.DeleteTag(ctx, &brzrpc.UserTagId{TagId: id, UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

// GetTagsByUser godoc
// @Summary tags by user
// @Description Returns all tags for user
// @Tags tag
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Tag
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag/by-user [get]
func (e *Echo) GetTagsByUser(c echo.Context) error {
	const op = "gateway.net.GetTagsByUser"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	if tgs, err := e.rdsAPI.API.GetTagsByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			st, ok := status.FromError(err)
			if !ok {
				log.Error(op, "REDIS ERROR", err)
			} else {
				if st.Code() != codes.NotFound {
					log.Error(op, "REDIS ERROR", err)
				}
			}
		} else {
			if st.Code() != codes.NotFound {
				st, ok := status.FromError(err)
				if !ok {
					log.Error(op, "REDIS ERROR", err)
				} else {
					if st.Code() != codes.NotFound {
						log.Error(op, "REDIS ERROR", err)
					}
				}
			}
		}
	} else {
		if tgs != nil {
			if len(tgs.GetItems()) != 0 {

				return c.JSON(http.StatusOK, tgs.GetItems())
			}
		}
	}

	tags, err := api.GetTagsByUser(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.SetTagsByUser(ctx, &brzrpc.TagsByUser{UserId: idUser, Items: tags.GetItems()}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if len(tags.Items) == 0 {
		tags.Items = []*brzrpc.Tag{}
	}

	return c.JSON(http.StatusOK, domain.ToTags(tags).Tgs)
}

// GetPinnedTagsByUser godoc
// @Summary tags by user
// @Description Returns all tags for user
// @Tags tag
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Tag
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/tag/pinned [get]
func (e *Echo) GetPinnedTagsByUser(c echo.Context) error {
	const op = "gateway.net.GetPinnedTagsByUser"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	tags, err := api.GetPinnedTagsByUser(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	return c.JSON(http.StatusOK, domain.ToTags(tags).Tgs)
}
