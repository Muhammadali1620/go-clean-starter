package models

import (
	"errors"
	"fmt"
)

type ErrorType string

const (
	ErrorTypeNotFound        ErrorType = "NOT_FOUND"         // 404
	ErrorTypeConflict        ErrorType = "CONFLICT"          // 409
	ErrorTypeValidation      ErrorType = "VALIDATION"        // 400
	ErrorTypeUnauthorized    ErrorType = "UNAUTHORIZED"      // 401
	ErrorTypeForbidden       ErrorType = "FORBIDDEN"         // 403
	ErrorTypeTooManyRequests ErrorType = "TOO_MANY_REQUESTS" // 429
	ErrorTypeInternal        ErrorType = "INTERNAL"          // 500
)

// AppError represents a structured, domain-level business error.
type AppError struct {
	Type    ErrorType // Maps to HTTP status code
	Code    ErrorCode // Strongly typed code for client i18n
	Message string    // Human-readable message in English
	Err     error     // Underlying wrapped error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a business error with an explicit typed code.
func NewAppError(errType ErrorType, code ErrorCode, message string) *AppError {
	if code == "" {
		code = ErrorCode(errType)
	}
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWrap wraps an existing error with domain context and typed code.
func NewAppErrorWrap(errType ErrorType, code ErrorCode, message string, err error) *AppError {
	var existingAppErr *AppError
	if errors.As(err, &existingAppErr) {
		return existingAppErr
	}

	if code == "" {
		code = ErrorCode(errType)
	}
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Err:     err,
	}
}
