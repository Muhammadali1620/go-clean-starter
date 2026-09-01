package dto

import "time"

// SearchOption defines parameters for text search across multiple fields.
type SearchOption struct {
	Query  string   // The search keyword
	Fields []string // Column names or JSONB expressions (e.g. []string{"name->>'uz'", "email"})
}

// DateRangeFilter defines a range between two UTC timestamps for a specific column.
type DateRangeFilter struct {
	Field string     // Column name (e.g. "created_at", "booking_time")
	From  *time.Time // GTE condition (>= From)
	To    *time.Time // LTE condition (<= To)
}

// OrderOption defines a single column order direction.
type OrderOption struct {
	Field     string // Column name
	Direction string // "ASC" or "DESC"
}

// QueryOptions contains flexible, DB-agnostic querying parameters.
type QueryOptions struct {
	Filters      map[string]any    // Strict equality: WHERE col = val
	InFilters    map[string][]any  // Inclusion: WHERE col IN (v1, v2)
	NotInFilters map[string][]any  // Exclusion: WHERE col NOT IN (v1, v2)
	DateRanges   []DateRangeFilter // Ranges: WHERE col >= from AND col <= to
	Search       *SearchOption     // Multi-field ILIKE search
	Limit        int               // Limit
	Offset       int               // Offset
	Orders       []OrderOption     // Ordered sorting list: ORDER BY col1 ASC, col2 DESC
	Relations    []string          // Relations to preload
}

// NewQueryOptions returns an initialized QueryOptions with allocated maps.
func NewQueryOptions() QueryOptions {
	return QueryOptions{
		Filters:      make(map[string]any),
		InFilters:    make(map[string][]any),
		NotInFilters: make(map[string][]any),
		DateRanges:   make([]DateRangeFilter, 0),
		Orders:       make([]OrderOption, 0),
		Relations:    make([]string, 0),
	}
}
