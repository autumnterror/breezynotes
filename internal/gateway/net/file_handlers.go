package net

import (
	"errors"
	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/log"

	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	FilesDir       = "./files"
	MaxUploadBytes = 10 << 20 // 10 MB
)

// UploadFile godoc
// @Tags files
// @Produce json
// @Success 200 {object} domain.Name
// @Failure 400 {object} domain.Error "Неверный формат данных"
// @Failure 500 {object} domain.Error "Ошибка на сервере"
// @Router /api/files [post]
func (e *Echo) UploadFile(c echo.Context) error {
	const op = "handlers.UploadFile"
	log.Blue(op)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "file field 'file' is required"})
	}

	src, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "cannot open uploaded file"})
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "file extension is required"})
	}

	filename := uuid.NewString() + ext
	dstPath := filepath.Join(FilesDir, filename)

	if err := saveMultipartFile(src, dstPath, fileHeader); err != nil {
		var status int
		var msg string

		if errors.Is(err, errTooLarge) {
			status = http.StatusBadRequest
			msg = "file is too large"
		} else {
			status = http.StatusInternalServerError
			msg = "failed to save file"
		}

		log.Error(op, "", err)
		return c.JSON(status, domain.Error{Error: msg})
	}

	return c.JSON(http.StatusOK, domain.Name{
		Name: filename,
	})
}

var errTooLarge = errors.New("too large")

func saveMultipartFile(src multipart.File, dstPath string, fh *multipart.FileHeader) error {
	if fh.Size > MaxUploadBytes {
		return errTooLarge
	}

	if err := os.MkdirAll(FilesDir, 0755); err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return err
	}

	if err := os.Chmod(dstPath, 0644); err != nil {
		return err
	}

	return nil
}

// DeleteFile godoc
// @Tags files
// @Produce json
// @Param title query string true "название файла"
// @Success 204
// @Failure 400 {object} domain.Error "Неверное название"
// @Failure 500 {object} domain.Error "Ошибка на сервере"
// @Router /api/files [delete]
func (e *Echo) DeleteFile(c echo.Context) error {
	const op = "handlers.DeleteFile"
	log.Blue(op)

	filename := c.QueryParam("title")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, domain.Error{Error: "empty filename"})
	}

	if err := deleteFile(filename); err != nil {
		switch {
		case errors.Is(err, ErrFileName):
			return c.JSON(http.StatusBadRequest, domain.Error{Error: "invalid filename"})
		case errors.Is(err, ErrFileNotFound):
			return c.JSON(http.StatusNotFound, domain.Error{Error: "file not found"})
		default:
			log.Error(op, "", err)
			return c.JSON(http.StatusInternalServerError, domain.Error{Error: "check logs"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

var ErrFileName = errors.New("invalid filename")
var ErrFileNotFound = errors.New("file not found")

func deleteFile(filename string) error {

	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return ErrFileName
	}

	fullPath := filepath.Clean(filepath.Join(filepath.Clean(FilesDir), filename))

	err := os.Remove(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}
