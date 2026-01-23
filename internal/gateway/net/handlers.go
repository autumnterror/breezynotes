package net

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// Healthz godoc
// @Summary check health of gateway
// @Description
// @Tags healthz
// @Produce json
// @Success 200 {object} domain.Message
// @Failure 502 {object} domain.Message
// @Router /api/healthz [get]
func (e *Echo) Healthz(c echo.Context) error {
	const op = "gateway.net.Healthz"

	ctx, done := context.WithTimeout(c.Request().Context(), time.Second)
	defer done()
	_, err := e.bnAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, domain.Message{Message: "bad blocknote"})
	}
	_, err = e.rdsAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, domain.Message{Message: "bad redis"})
	}
	_, err = e.authAPI.API.Healthz(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusBadGateway, domain.Message{Message: "bad auth"})
	}

	return c.JSON(http.StatusOK, domain.Message{Message: "HEALTHZ"})
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
// @Failure 400 {object} domain.Error
// @Failure 408 {object} []brzrpc.NotePart
// @Failure 502 {object} domain.Error
// @Failure 504 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /api/note/find [get]
func (e *Echo) Search(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
	//const op = "gateway.net.Search"
	//
	//api := e.bnAPI.API
	//
	//idUser, errGetId := getIdUser(c)
	//if errGetId != nil {
	//	return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad idUser from access token"})
	//}
	//
	//start, end, resPag := getPagination(c)
	//if resPag != nil {
	//	if r, ok := resPag.(domain.Error); ok {
	//		return c.JSON(http.StatusBadRequest, r)
	//	}
	//	if r, ok := resPag.(*brzrpc.NoteParts); ok {
	//		return c.JSON(http.StatusOK, r)
	//	}
	//}
	//
	//prompt := c.QueryParam("prompt")
	//if prompt == "" {
	//	return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	//}
	//prompt = strings.ToLower(prompt)
	//
	//ctx, done := context.WithTimeout(c.Request().Context(), 10*time.Second)
	//defer done()
	//
	//res := make(chan interface{})
	//var nts []*brzrpc.NotePart
	//
	//go func() {
	//	if nl, err := e.rdsAPI.API.GetNoteListByUser(ctx, &brzrpc.UserId{UserId: idUser}); err != nil {
	//		log.Error(op, "REDIS ERROR", err)
	//	} else {
	//		if nl != nil {
	//			if nl.GetItems() != nil {
	//				if len(nl.GetItems()) != 0 {
	//					nts = nl.GetItems()
	//				}
	//			}
	//		}
	//	}
	//
	//	if len(nts) == 0 {
	//		notes, err := api.GetAllNotes(ctx, &brzrpc.UserId{
	//			UserId: idUser,
	//		})
	//		if err != nil {
	//			st, ok := status.FromError(err)
	//			if !ok {
	//				log.Error(op, "get all notes error", err)
	//				res <- domain.Error{Error: "get all notes error"}
	//			}
	//
	//			switch st.Code() {
	//			default:
	//				log.Error(op, "get all notes error", err)
	//				res <- domain.Error{Error: "get all notes error"}
	//			}
	//		}
	//		nts = notes.GetItems()
	//	}
	//
	//	var finds []*brzrpc.NotePart
	//
	//	for _, n := range nts {
	//		if strings.Contains(strings.ToLower(n.GetTitle()), prompt) {
	//			finds = append(finds, n)
	//		}
	//	}
	//	if len(finds) == 0 {
	//		for _, n := range nts {
	//			if note, err := api.GetNote(ctx, &brzrpc.NoteId{NoteId: n.GetId()}); err == nil {
	//				for _, b := range note.GetBlocks() {
	//					if block, err := api.GetBlockAsFirst(ctx, &brzrpc.BlockId{BlockId: b}); err == nil {
	//						if strings.Contains(strings.ToLower(block.GetValue()), prompt) {
	//							finds = append(finds, n)
	//						}
	//					}
	//				}
	//			}
	//		}
	//	}
	//	items := finds
	//	if start >= len(items) {
	//		res <- []*brzrpc.NotePart{}
	//	}
	//	if end > len(items) {
	//		end = len(items)
	//	}
	//	nlPad := items[start:end]
	//	res <- nlPad
	//}()
	//
	//select {
	//case <-ctx.Done():
	//	if len(nts) == 0 {
	//		nts = []*brzrpc.NotePart{}
	//	}
	//	return c.JSON(http.StatusRequestTimeout, nts)
	//case r := <-res:
	//	if rs, ok := r.(domain.Error); ok {
	//		return c.JSON(http.StatusBadGateway, rs)
	//	}
	//	if rs, ok := r.([]*brzrpc.NotePart); ok {
	//
	//		return c.JSON(http.StatusOK, rs)
	//	}
	//	return c.JSON(http.StatusInternalServerError, domain.Error{Error: "bad server"})
	//}
}
