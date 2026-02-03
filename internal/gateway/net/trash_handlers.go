package net

import (
	"context"
	"net/http"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CleanTrash godoc
// @Summary Clean trash
// @Description Deletes all notes from trash for user
// @Tags trash
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/trash [delete]
func (e *Echo) CleanTrash(c echo.Context) error {
	const op = "gateway.net.CleanTrash"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.CleanTrash(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

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
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/trash/to [put]
func (e *Echo) NoteToTrash(c echo.Context) error {
	const op = "gateway.net.NoteToTrash"

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad param"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.NoteToTrash(ctx, &brzrpc.UserNoteId{NoteId: id, UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: id}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

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
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/trash/from [put]
func (e *Echo) NoteFromTrash(c echo.Context) error {
	const op = "gateway.net.NoteFromTrash"

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.NoteFromTrash(ctx, &brzrpc.UserNoteId{NoteId: id, UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.RmNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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

	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: id}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

// GetNotesFromTrash godoc
// @Summary GetNote notes from trash
// @Description Returns notes from trash by user ID
// @Tags trash
// @Accept json
// @Produce json
// @Success 200 {object} []brzrpc.NotePart
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/trash [get]
func (e *Echo) GetNotesFromTrash(c echo.Context) error {
	const op = "gateway.net.GetNotesFromTrash"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	if ntsT, err := e.rdsAPI.API.GetNotesFromTrashByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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
		if ntsT != nil {
			if ntsT.GetItems() != nil {
				if len(ntsT.GetItems()) != 0 {

					return c.JSON(http.StatusOK, ntsT.GetItems())
				}
			}
		}
	}

	nts, err := api.GetNotesFromTrash(ctx, &brzrpc.UserId{UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.SetNotesFromTrashByUser(ctx, &brzrpc.NoteListByUser{UserId: idUser, Items: nts.GetItems()}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	if nts.GetItems() == nil {
		return c.JSON(http.StatusOK, []brzrpc.NotePart{})
	}
	if len(nts.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []brzrpc.NotePart{})
	}

	return c.JSON(http.StatusOK, nts.GetItems())
}
