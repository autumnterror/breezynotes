package net

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
	"github.com/autumnterror/utils_go/pkg/utils/uid"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/labstack/echo/v4"
)

// ChangeTitleNote godoc
// @Summary Change note title
// @Description Changes title of existing note
// @Tags note
// @Accept json
// @Produce json
// @Param ChangeTitleNoteRequest body domain.ChangeTitleNoteRequest true "Note ID and new title"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/title [patch]
func (e *Echo) ChangeTitleNote(c echo.Context) error {
	const op = "gateway.net.ChangeTitleNote"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.ChangeTitleNoteRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.ChangeTitleNote(ctx, &brzrpc.ChangeTitleNoteRequest{
		IdNote: r.Id,
		Title:  r.Title,
		IdUser: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
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
	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.Id}); err != nil {
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

// GetNote godoc
// @Summary GetNote note
// @Description Returns note by ID
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true "Note ID"
// @Success 200 {object} domain.NoteWithBlocks
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note [get]
func (e *Echo) GetNote(c echo.Context) error {
	const op = "gateway.net.Get"

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

	if note, err := e.rdsAPI.API.GetNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: id}); err != nil {
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
		if note != nil {
			if note.GetAuthor() != idUser && !alg.IsIn(idUser, note.GetEditors()) && !alg.IsIn(idUser, note.GetReaders()) {
				return c.JSON(http.StatusUnauthorized, domain.Error{Error: "user dont have permission"})
			}
			return c.JSON(http.StatusOK, note)
		}
	}

	note, err := api.GetNote(ctx, &brzrpc.UserNoteId{NoteId: id, UserId: idUser})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.SetNoteByUser(ctx, &brzrpc.NoteByUser{UserId: idUser, Note: note}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	return c.JSON(http.StatusOK, note)
}

// GetAllNotes godoc
// @Summary GetNote all notes of user
// @Description Returns all notes by user ID
// @Tags note
// @Accept json
// @Produce json
// @Param start query int true  "start > 0"
// @Param end query int true  "end"
// @Success 200 {object} domain.NoteListPaginationResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/all [get]
func (e *Echo) GetAllNotes(c echo.Context) error {
	const op = "gateway.net.GetAllNotes"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	start, end, resPag := getPagination(c)
	if resPag != nil {
		if r, ok := resPag.(domain.Error); ok {
			return c.JSON(http.StatusBadRequest, r)
		}
		if r, ok := resPag.(*brzrpc.NoteParts); ok {
			return c.JSON(http.StatusOK, r)
		}
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	//---------------REDIS---------------
	if nl, err := e.rdsAPI.API.GetNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
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
		if nl != nil {
			if nl.GetItems() != nil {
				if len(nl.GetItems()) != 0 {
					items := nl.GetItems()
					if start >= len(items) {
						return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
					}
					if end > len(items) {
						end = len(items)
					}
					nlPad := items[start:end]

					return c.JSON(http.StatusOK, nlPad)
				}
			}
		}
	}
	//---------------REDIS---------------

	notes, err := api.GetAllNotes(ctx, &brzrpc.UserId{
		UserId: idUser,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	if _, err := e.rdsAPI.API.SetNoteListByUser(ctx, &brzrpc.NoteListByUser{UserId: idUser, Items: notes.GetItems()}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}

	items := notes.GetItems()
	if start >= len(items) {
		return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	}
	if end > len(items) {
		end = len(items)
	}
	nlPag := items[start:end]

	return c.JSON(http.StatusOK, domain.NoteListPaginationResponse{
		Items: domain.ToNotePartList(nlPag),
		Total: len(items),
	})
}

// GetNotesByTag godoc
// @Summary GetNote notes by tag
// @Description Returns all notes that contain given tag
// @Tags note
// @Accept json
// @Produce json
// @Param id query string true  "Tag ID"
// @Param start query int true  "start > 0"
// @Param end query int true  "end"
// @Success 200 {object} domain.NoteListPaginationResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/by-tag [get]
func (e *Echo) GetNotesByTag(c echo.Context) error {
	const op = "gateway.net.GetNotesByTag"

	api := e.bnAPI.API

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad param"})
	}

	start, end, resPag := getPagination(c)
	if resPag != nil {
		if r, ok := resPag.(domain.Error); ok {
			return c.JSON(http.StatusBadRequest, r)
		}
		if r, ok := resPag.(*brzrpc.NoteParts); ok {
			return c.JSON(http.StatusOK, r)
		}
	}

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	notes, err := api.GetNotesByTag(ctx, &brzrpc.UserTagId{
		TagId:  id,
		UserId: idUser,
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	items := notes.GetItems()
	if start >= len(items) {
		return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	}
	if end > len(items) {
		end = len(items)
	}
	nlPag := items[start:end]

	return c.JSON(http.StatusOK, domain.NoteListPaginationResponse{
		Items: domain.ToNotePartList(nlPag),
		Total: len(items),
	})
}

// CreateNote godoc
// @Summary Create note
// @Description Creates new note
// @Tags note
// @Accept json
// @Produce json
// @Param Note body domain.CreateNoteRequest true "Note info"
// @Success 201 {object} brzrpc.Id
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note [post]
func (e *Echo) CreateNote(c echo.Context) error {
	const op = "gateway.net.CreateNote"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.CreateNoteRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	id := uid.New()
	_, err := api.CreateNote(ctx, &brzrpc.Note{
		Id:        id,
		Title:     r.Title,
		CreatedAt: time.Now().UTC().Unix(),
		UpdatedAt: time.Now().UTC().Unix(),
		Tag:       nil,
		Author:    idUser,
		Editors:   []string{},
		Readers:   []string{},
		Blocks:    []string{},
	})
	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
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

	return c.JSON(http.StatusCreated, brzrpc.Id{Id: id})
}

// AddTagToNote godoc
// @Summary Add tag to note
// @Description Attaches tag to note
// @Tags note
// @Accept json
// @Produce json
// @Param AddTagToNoteRequest body domain.NoteTagId true "Note ID and Tag ID"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/tag [post]
func (e *Echo) AddTagToNote(c echo.Context) error {
	const op = "gateway.net.AddTagToNote"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.NoteTagId
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.AddTagToNote(ctx, &brzrpc.NoteTagUserId{
		NoteId: r.NoteId,
		TagId:  r.TagId,
		UserId: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
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

	return c.NoContent(http.StatusOK)
}

// RmTagFromNote godoc
// @Summary Remove tag from note
// @Description Remove tag from note
// @Tags note
// @Accept json
// @Produce json
// @Param AddTagToNoteRequest body domain.NoteTagId true "Note ID and Tag ID"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/tag [delete]
func (e *Echo) RmTagFromNote(c echo.Context) error {
	const op = "gateway.net.RmTagFromNote"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.NoteId
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	_, err := api.RemoveTagFromNote(ctx, &brzrpc.UserNoteId{
		NoteId: r.NoteId,
		UserId: idUser,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
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

	return c.NoContent(http.StatusOK)
}

// ShareNote godoc
// @Summary share note
// @Description add new user to list of editors or readers
// @Tags note
// @Accept json
// @Produce json
// @Param Note body domain.ShareNoteRequest true "share info"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/share [patch]
func (e *Echo) ShareNote(c echo.Context) error {
	const op = "gateway.net.ShareNote"

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
	}

	var r domain.ShareNoteRequest
	if err := c.Bind(&r); err != nil {

		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
	}

	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
	defer done()

	id, err := e.authAPI.API.GetIdFromLogin(ctx, &brzrpc.String{Value: r.Login})
	code, errRes := authErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	_, err = api.ShareNote(ctx, &brzrpc.ShareNoteRequest{
		NoteId:        r.NoteId,
		UserIdToShare: id.Id,
		Role:          r.Role,
		UserId:        idUser,
	})
	code, errRes = bNErrors(op, err)
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
	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: id.Id, NoteId: r.NoteId}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if st.Code() != codes.NotFound {
				log.Error(op, "REDIS ERROR", err)
			}
		}
	}
	if _, err := e.rdsAPI.API.RmNoteListByUser(ctx, &brzrpc.UserId{UserId: id.Id}); err != nil {
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

// ChangeUserRole godoc
// @Summary share note
// @Description replace user to list of editors or readers
// @Tags note
// @Accept json
// @Produce json
// @Param Note body domain.ChangeRoleRequest true "share info"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Router /api/note/role [patch]
// func (e *Echo) ChangeUserRole(c echo.Context) error {
// 	const op = "gateway.net.ShareNote"

// 	api := e.bnAPI.API

// 	idUser, errGetId := getIdUser(c)
// 	if errGetId != nil {
// 		return c.JSON(http.StatusUnauthorized, domain.Error{Error: "bad idUser from access token"})
// 	}

// 	var r domain.ChangeRoleRequest
// 	if err := c.Bind(&r); err != nil {

// 		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad JSON"})
// 	}

// 	ctx, done := context.WithTimeout(c.Request().Context(), domain.WaitTime)
// 	defer done()

// 	id, err := e.authAPI.API.GetIdFromLogin(ctx, &brzrpc.String{Value: r.Login})
// 	code, errRes := authErrors(op, err)
// 	if code != http.StatusOK {
// 		return c.JSON(code, errRes)
// 	}

// 	_, err = api.ChangeUserRole(ctx, &brzrpc.ChangeUserRoleRequest{
// 		UserIdToChange: id.Id,
// 		NewRole:        r.NewRole,
// 		NoteId:         r.NoteId,
// 		UserId:         idUser,
// 	})
// 	code, errRes = bNErrors(op, err)
// 	if code != http.StatusOK {
// 		return c.JSON(code, errRes)
// 	}

// 	if _, err := e.rdsAPI.API.RmNoteByUser(ctx, &brzrpc.UserNoteId{UserId: idUser, NoteId: r.NoteId}); err != nil {
// 		st, ok := status.FromError(err)
// 		if !ok {
// 			log.Error(op, "REDIS ERROR", err)
// 		} else {
// 			if st.Code() != codes.NotFound {
// 				log.Error(op, "REDIS ERROR", err)
// 			}
// 		}
// 	}

// 	return c.NoContent(http.StatusOK)
// }
