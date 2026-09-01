package dto

import "reflect"

// BaseResponse is the standard generic response wrapper for all API responses.
type BaseResponse[T any] struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// PaginatedData holds the list of items and pagination metadata.
type PaginatedData[T any] struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Items      []T `json:"items"`
}

// NewResponse creates a standard success response with generic data.
func NewResponse[T any](data T, success bool) BaseResponse[T] {
	return BaseResponse[T]{
		Success: success,
		Data:    data,
	}
}

// NewErrorResponse creates an error response with an explicit client error code.
func NewErrorResponse(code, message string) BaseResponse[any] {
	return BaseResponse[any]{
		Success: false,
		Code:    code,
		Message: message,
	}
}

// NewPaginatedResponse creates a standard paginated response.
// Crucial: Guarantee that Items is NEVER serialized as null in JSON (always [] if empty/nil).
func NewPaginatedResponse[T any](items []T, total, limit, offset int) BaseResponse[PaginatedData[T]] {
	if items == nil {
		items = make([]T, 0)
	}

	currentPage := 1
	if limit > 0 {
		currentPage = (offset / limit) + 1
	}

	return BaseResponse[PaginatedData[T]]{
		Success: true,
		Data: PaginatedData[T]{
			TotalCount: total,
			Page:       currentPage,
			PageSize:   limit,
			Items:      items,
		},
	}
}

// NewAnyPaginatedResponse is a dynamic fallback for untyped or interface items.
func NewAnyPaginatedResponse(items any, total, limit, offset int) BaseResponse[any] {
	// Guard against nil slices converting to JSON null
	safeItems := items
	if safeItems == nil {
		safeItems = []any{}
	} else {
		val := reflect.ValueOf(safeItems)
		if val.Kind() == reflect.Slice && val.IsNil() {
			safeItems = []any{}
		}
	}

	currentPage := 1
	if limit > 0 {
		currentPage = (offset / limit) + 1
	}

	return BaseResponse[any]{
		Success: true,
		Data: map[string]any{
			"total_count": total,
			"page":        currentPage,
			"page_size":   limit,
			"items":       safeItems,
		},
	}
}
