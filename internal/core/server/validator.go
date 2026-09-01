package server

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"new_project/internal/models"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, ve := range validationErrs {
				field := strings.ToLower(ve.Field())
				switch ve.Tag() {
				case "required":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is required", field))
				case "min":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be at least %s", field, ve.Param()))
				case "max":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' cannot exceed %s", field, ve.Param()))
				case "oneof":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be one of [%s]", field, ve.Param()))
				case "gt":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be greater than %s", field, ve.Param()))
				case "gte":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be greater than or equal to %s", field, ve.Param()))
				case "numeric":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be a valid number", field))
				default:
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' failed on rule '%s'", field, ve.Tag()))
				}
			}
			return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeValidationFailed, strings.Join(errMsgs, "; "))
		}
		return models.NewAppErrorWrap(models.ErrorTypeValidation, models.ErrCodeValidationFailed, "validation failed", err)
	}
	return nil
}
