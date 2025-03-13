package photo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/response"
	"github.com/labstack/echo/v4"
	"image/png"
	"net/http"
	"strings"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

// GetById
// @Tags Photo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /photo/{id} [get]
func (h Handler) GetById(c echo.Context) error {
	var err error

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	tPhoto, err := h.usecase.GetById(id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	// Decode Base64 ke byte array
	imgData, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(tPhoto.PhotoPath, "data:image/png;base64,"))
	if err != nil {
		return response.Error(http.StatusBadRequest, fmt.Errorf("error decoding base64: %v", err).Error(), err, nil).SendJSON(c)
	}

	// Konversi ke format PNG (opsional)
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error decoding image")
	}

	// Simpan ke buffer
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return c.String(http.StatusInternalServerError, "Error encoding image")
	}

	// Kirim gambar sebagai response
	return c.Blob(http.StatusOK, "image/png", buf.Bytes())
}
