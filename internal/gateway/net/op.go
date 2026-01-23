package net

import (
	"errors"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/labstack/echo/v4"
	"strconv"
)

func getPagination(c echo.Context) (int, int, interface{}) {
	s := c.QueryParam("start")
	if s == "" {
		return 0, 0, domain.Error{Error: "bad start"}
	}
	en := c.QueryParam("end")
	if en == "" {
		return 0, 0, domain.Error{Error: "bad end"}
	}

	start, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, domain.Error{Error: "start must be int"}
	}
	end, err := strconv.Atoi(en)
	if err != nil {
		return 0, 0, domain.Error{Error: "end must be int"}
	}

	if start >= end {
		return 0, 0, &brzrpc.NoteParts{Items: []*brzrpc.NotePart{}}
	}

	if start < 0 {
		return 0, 0, domain.Error{Error: "start < 0!"}
	}

	return start, end, nil
}

func getIdUser(c echo.Context) (string, error) {
	idInt := c.Get(domain.IdFromContext)
	idUser, ok := idInt.(string)
	if !ok && idUser == "" {
		return "", errors.New("bad id")
	}
	return idUser, nil
}
