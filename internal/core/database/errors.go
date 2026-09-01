package database

import (
	"database/sql"
	"errors"

	"github.com/uptrace/bun/driver/pgdriver"

	"new_project/internal/models"
)

// MapDBError converts raw database errors into our application's structured domain errors.
func MapDBError(err error) error {
	if err == nil {
		return nil
	}

	// 1. Not Found
	if errors.Is(err, sql.ErrNoRows) {
		return models.NewAppError(models.ErrorTypeNotFound, models.ErrCodeNotFound, "resource not found")
	}

	// 2. PostgreSQL specific errors
	var pgErr pgdriver.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Field('C') {
		case "23505": // Unique constraint violation
			return models.NewAppError(models.ErrorTypeConflict, models.ErrCodeAlreadyExists, "resource already exists")
		case "23503": // Foreign key violation (используем твой новый код!)
			return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeForeignKeyViolation, "related resource not found (invalid reference)")
		case "42703": // Undefined column
			return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid query parameter: field does not exist")
		case "22001": // Value too long
			return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeValidationFailed, "input text exceeds maximum allowed length")
		}
	}

	return models.NewAppErrorWrap(models.ErrorTypeInternal, models.ErrCodeInternalError, "database query execution failed", err)
}
