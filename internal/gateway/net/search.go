package net

import (
	"context"
	"encoding/json"
	"fmt"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

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
// @Router /api/note/search [get]
func (e *Echo) Search(c echo.Context) error {
	const op = "gateway.net.Search"

	w := c.Response().Writer
	r := c.Request()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported!")
	}

	api := e.bnAPI.API
	idUser, errGetId := getIdUser(c)
	if errGetId != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "bad idUser from access token"})
	}

	//start, end, resPag := getPagination(c)
	//if resPag != nil {
	//	if r, ok := resPag.(domain.Error); ok {
	//		return c.JSON(http.StatusBadRequest, r)
	//	}
	//	if r, ok := resPag.(*brzrpc.NoteParts); ok {
	//		return c.JSON(http.StatusOK, r)
	//	}
	//}

	prompt := c.QueryParam("prompt")
	if prompt == "" {
		return c.JSON(http.StatusOK, []*brzrpc.NotePart{})
	}
	prompt = strings.ToLower(prompt)

	ctx, done := context.WithTimeout(r.Context(), 10*time.Second)
	defer done()

	notes, err := api.Search(ctx, &brzrpc.SearchRequest{
		UserId: idUser,
		Prompt: prompt,
	})

	code, errRes := bNErrors(op, err)
	if code != http.StatusOK {
		return c.JSON(code, errRes)
	}

	total := 0

	for {
		note, err := notes.Recv()
		if err == io.EOF {
			// Поток завершен сервером
			break
		}
		if err != nil {
			break
		}
		total++
		jsonData, err := json.Marshal(struct {
			Np    *brzrpc.NotePart `json:"note"`
			Total int              `json:"total"`
		}{note, total})
		if err != nil {
			log.Error(op, "marshall note", err)
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}

	return nil
}
