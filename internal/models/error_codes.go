package models

// ErrorCode is a strongly typed domain error code for client i18n mapping.
type ErrorCode string

const (
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists       ErrorCode = "ALREADY_EXISTS"
	ErrCodeValidationFailed    ErrorCode = "VALIDATION_FAILED"
	ErrCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrCodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"
	ErrCodeInternalError       ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidCredentials  ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeAccountBlocked      ErrorCode = "ACCOUNT_BLOCKED"
	ErrCodeInvalidDateRange    ErrorCode = "INVALID_DATE_RANGE"
	ErrCodeInvalidInput        ErrorCode = "INVALID_INPUT"
	ErrCodeConflict            ErrorCode = "CONFLICT"
	ErrCodeForeignKeyViolation ErrorCode = "FOREIGN_KEY_VIOLATION"
)
