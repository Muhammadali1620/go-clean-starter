package handlers

import (
	"net/http"

	"new_project/internal/dto"
	"new_project/internal/models"
	"new_project/internal/services"

	"github.com/labstack/echo/v4"
)

type MediaHandler struct {
	mediaService services.MediaService
}

// @inject
func NewMediaHandler(mediaService services.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// @Summary      Upload a file
// @Description  Upload a file to the server
// @Tags         Media
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "File to upload"
// @Param        folder  formData  string  false  "Folder to upload to"
// @Success      200  {object}  dto.BaseResponse[dto.MediaUploadResponse]
// @Failure      400  {object}  models.AppError
// @Failure      500  {object}  models.AppError
// @Router       /api/v1/media/upload [post]
// @Security     BearerAuth
func (h *MediaHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "File is required under 'file' key")
	}

	folder := c.FormValue("folder")
	if folder == "" {
		folder = "misc"
	}
	filePath, err := h.mediaService.Upload(c.Request().Context(), file, folder)
	if err != nil {
		return err
	}

	response := dto.MediaUploadResponse{
		URL: filePath,
	}

	return c.JSON(http.StatusOK, dto.NewResponse(response, true))
}
