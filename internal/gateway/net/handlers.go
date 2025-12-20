package net

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/views"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
)

// Healthz godoc
// @Summary check health of gateway
// @Description
// @Tags healthz
// @Produce json
// @Success 200 {object} views.SWGMessage
// @Failure 502 {object} views.SWGMessage
// @Router /api/healthz [get]
func (e *Echo) Healthz(c echo.Context) error {
	const op = "gateway.net.Healthz"
	log.Info(op, "")

	ctx, done := context.WithTimeout(c.Request().Context(), time.Second)
	defer done()
	_, err := e.bnAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.SWGMessage{Message: "bad blocknote"})
	}
	_, err = e.rdsAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.SWGMessage{Message: "bad redis"})
	}
	_, err = e.authAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.SWGMessage{Message: "bad auth"})
	}

	log.Success(op, "")
	return c.JSON(http.StatusOK, views.SWGMessage{Message: "HEALTHZ"})
}

// Search godoc
// @Summary search note by title or blocks inside
// @Description if not fiend title get all block as first and fiend sames
// @Tags note
// @Accept json
// @Produce json
// @Param start query int true  "start > 0"
// @Param end query int true  "end"
// @Param prompt query string true  "prompt"
// @Success 200 {object} []brzrpc.NotePart
// @Failure 400 {object} views.SWGError
// @Failure 408 {object} []brzrpc.NotePart
// @Failure 502 {object} views.SWGError
// @Failure 500 {object} views.SWGError
// @Router /api/note/find [get]
func (e *Echo) Search(c echo.Context) error {
	const op = "gateway.net.Search"
	log.Info(op, "")

	api := e.bnAPI.API

	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusBadRequest, views.SWGError{Error: "bad idUser from access token"})
	}

	start, end, resPag := getPagination(c)
	if resPag != nil {
		if r, ok := resPag.(views.SWGError); ok {
			return c.JSON(http.StatusBadRequest, r)
		}
		if r, ok := resPag.(*brzrpc.NoteParts); ok {
			return c.JSON(http.StatusOK, r)
		}
	}

	prompt := c.QueryParam("prompt")
	if prompt == "" {
		return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	}
	prompt = strings.ToLower(prompt)

	ctx, done := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer done()

	res := make(chan interface{})
	var nts []*brzrpc.NotePart

	go func() {
		if nl, err := e.rdsAPI.API.GetNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
			log.Error(op, "REDIS ERROR", err)
		} else {
			if nl != nil {
				if nl.GetItems() != nil {
					if len(nl.GetItems()) != 0 {
						nts = nl.GetItems()
					}
				}
			}
		}

		if len(nts) == 0 {
			notes, err := api.GetAllNotes(ctx, &brzrpc.UserId{
				UserId: idUser,
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					log.Error(op, "get all notes error", err)
					res <- views.SWGError{Error: "get all notes error"}
				}

				switch st.Code() {
				default:
					log.Error(op, "get all notes error", err)
					res <- views.SWGError{Error: "get all notes error"}
				}
			}
			nts = notes.GetItems()
		}

		var finds []*brzrpc.NotePart

		for _, n := range nts {
			if strings.Contains(strings.ToLower(n.GetTitle()), prompt) {
				finds = append(finds, n)
			}
		}
		if len(finds) == 0 {
			for _, n := range nts {
				if note, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: n.GetId()}); err == nil {
					for _, b := range note.GetBlocks() {
						if block, err := api.GetBlockAsFirst(ctx, &brzrpc.BlockId{BlockId: b}); err == nil {
							if strings.Contains(strings.ToLower(block.GetValue()), prompt) {
								finds = append(finds, n)
							}
						}
					}
				}
			}
		}
		items := finds
		if start >= len(items) {
			res <- []*brzrpc.NotePart{}
		}
		if end > len(items) {
			end = len(items)
		}
		nlPad := items[start:end]
		res <- nlPad
	}()

	select {
	case <-ctx.Done():
		if len(nts) == 0 {
			nts = []*brzrpc.NotePart{}
		}
		return c.JSON(http.StatusRequestTimeout, nts)
	case r := <-res:
		if rs, ok := r.(views.SWGError); ok {
			return c.JSON(http.StatusBadGateway, rs)
		}
		if rs, ok := r.([]*brzrpc.NotePart); ok {
			log.Success(op, "")
			return c.JSON(http.StatusOK, rs)
		}
		return c.JSON(http.StatusInternalServerError, views.SWGError{Error: "bad server"})
	}
}
