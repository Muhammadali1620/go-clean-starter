package utils

import (
	"time"

	"new_project/internal/models"
)

// NowUTC returns the current time strictly in UTC+0.
func NowUTC() time.Time {
	return time.Now().UTC()
}

// ParseISO8601 parses a standard ISO8601 / RFC3339 string into UTC.
func ParseISO8601(val string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid ISO8601 timestamp (expected RFC3339, e.g. 2026-03-30T15:04:05Z)")
	}
	return t.UTC(), nil
}

// ParseDateUTC parses a simple YYYY-MM-DD date and sets time to 00:00:00 UTC.
func ParseDateUTC(val string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return time.Time{}, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid date format (expected YYYY-MM-DD)")
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
}

// ParseDateTimeUTC parses "YYYY-MM-DD HH:mm:ss" or "YYYY-MM-DD" in UTC.
func ParseDateTimeUTC(val string) (time.Time, error) {
	if val == "" {
		return time.Time{}, nil
	}
	layout := "2006-01-02 15:04:05"
	if len(val) == 10 {
		layout = "2006-01-02"
	}

	t, err := time.Parse(layout, val)
	if err != nil {
		return time.Time{}, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid datetime format (expected YYYY-MM-DD HH:mm:ss or YYYY-MM-DD)")
	}
	return t.UTC(), nil
}

// ParseDateRangeUTC parses from_date and to_date (YYYY-MM-DD) into start-of-day and end-of-day UTC times.
func ParseDateRangeUTC(fromStr, toStr string) (*time.Time, *time.Time, error) {
	var from, to *time.Time

	if fromStr != "" {
		parsedFrom, err := ParseDateUTC(fromStr)
		if err != nil {
			return nil, nil, err
		}
		from = &parsedFrom
	}

	if toStr != "" {
		parsedTo, err := ParseDateUTC(toStr)
		if err != nil {
			return nil, nil, err
		}
		// Set to end of day 23:59:59.999999999 UTC
		endOfDay := time.Date(parsedTo.Year(), parsedTo.Month(), parsedTo.Day(), 23, 59, 59, 999999999, time.UTC)
		to = &endOfDay
	}

	if from != nil && to != nil && from.After(*to) {
		return nil, nil, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidDateRange, "from_date cannot be after to_date")
	}

	return from, to, nil
}

// ParseDateTimeRangeUTC parses from and to strings (ISO8601 or YYYY-MM-DD HH:mm:ss) into exact UTC timestamps.
func ParseDateTimeRangeUTC(fromStr, toStr string) (*time.Time, *time.Time, error) {
	var from, to *time.Time

	if fromStr != "" {
		parsedFrom, err := ParseDateTimeUTC(fromStr)
		if err != nil {
			return nil, nil, err
		}
		from = &parsedFrom
	}

	if toStr != "" {
		parsedTo, err := ParseDateTimeUTC(toStr)
		if err != nil {
			return nil, nil, err
		}
		to = &parsedTo
	}

	if from != nil && to != nil && from.After(*to) {
		return nil, nil, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidDateRange, "from_datetime cannot be after to_datetime")
	}

	return from, to, nil
}
