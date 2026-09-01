package server

import (
	"errors"
	"net/http"

	"new_project/internal/dto"
	"new_project/internal/models"

	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler catches errors and maps them to standard BaseResponse with error codes.
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	httpCode := http.StatusInternalServerError
	errorCode := models.ErrCodeInternalError
	message := "internal server error"

	var appErr *models.AppError
	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			c.Logger().Errorf("[DOMAIN ERROR] %s: %v", appErr.Code, appErr.Err)
		}

		switch appErr.Type {
		case models.ErrorTypeNotFound:
			httpCode = http.StatusNotFound
		case models.ErrorTypeConflict:
			httpCode = http.StatusConflict
		case models.ErrorTypeValidation:
			httpCode = http.StatusBadRequest
		case models.ErrorTypeUnauthorized:
			httpCode = http.StatusUnauthorized
		case models.ErrorTypeForbidden:
			httpCode = http.StatusForbidden
		case models.ErrorTypeTooManyRequests:
			httpCode = http.StatusTooManyRequests
		default:
			httpCode = http.StatusInternalServerError
		}

		errorCode = appErr.Code
		message = appErr.Message
	} else if he, ok := err.(*echo.HTTPError); ok {
		httpCode = he.Code
		if he.Code == http.StatusNotFound {
			errorCode = models.ErrCodeNotFound
		} else if he.Code == http.StatusUnauthorized {
			errorCode = models.ErrCodeUnauthorized
		}
		if msgStr, ok := he.Message.(string); ok {
			message = msgStr
		}
	} else {
		c.Logger().Error(err)
	}

	_ = c.JSON(httpCode, dto.NewErrorResponse(string(errorCode), message))
}
