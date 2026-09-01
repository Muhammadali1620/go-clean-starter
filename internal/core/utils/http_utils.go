package utils

import (
	"regexp"
	"strconv"
	"strings"

	"new_project/internal/dto"
	"new_project/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

var orderRegex = regexp.MustCompile(`^[a-zA-Z0-9_\.]+(?::|\s+)(?i)(ASC|DESC)$`)

// GetPagingParams safely extracts page and pageSize (default 1 and 10, max 100).
func GetPagingParams(c echo.Context) (page int, pageSize int) {
	page = 1
	pageSize = 10

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.QueryParam("pageSize")); err == nil && ps > 0 {
		pageSize = ps
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

// GetStringQueryParam safely extracts a query string.
func GetStringQueryParam(c echo.Context, param string) string {
	return strings.TrimSpace(c.QueryParam(param))
}

// GetQueryParamInt safely extracts a query param as int64.
func GetQueryParamInt(c echo.Context, name string) (int64, error) {
	valStr := GetStringQueryParam(c, name)
	if valStr == "" {
		return 0, nil
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid integer query parameter: "+name)
	}
	return val, nil
}

// GetPathParamInt safely extracts a path param (e.g. /:id) as int64.
func GetPathParamInt(c echo.Context, name string) (int64, error) {
	valStr := c.Param(name)
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid path parameter: "+name)
	}
	return val, nil
}

// GetStringSliceQueryParam extracts comma-separated strings (e.g., ?status=ACTIVE,PENDING).
func GetStringSliceQueryParam(c echo.Context, name string) []string {
	query := c.QueryParam(name)
	if query == "" {
		return nil
	}
	parts := strings.Split(query, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

// GetIntSliceQueryParam extracts comma-separated integers (e.g., ?ids=1,2,3).
func GetIntSliceQueryParam(c echo.Context, name string) []int64 {
	strSlice := GetStringSliceQueryParam(c, name)
	if len(strSlice) == 0 {
		return nil
	}
	res := make([]int64, 0, len(strSlice))
	for _, s := range strSlice {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			res = append(res, v)
		}
	}
	return res
}

// ParsePagination safely extracts page, pageSize, and multi-order clauses.
// Supports comma-separated ordering: "?order=rating:desc,distance:asc" or "?order=created_at DESC"
func ParsePagination(c echo.Context) dto.QueryOptions {
	opts := dto.NewQueryOptions()

	page, pageSize := GetPagingParams(c)
	opts.Limit = pageSize
	opts.Offset = (page - 1) * pageSize

	orderQuery := c.QueryParam("order")
	if orderQuery != "" {
		orderPairs := strings.Split(orderQuery, ",")
		for _, pair := range orderPairs {
			pair = strings.TrimSpace(pair)
			// Normalize colon syntax (e.g. "rating:desc" -> "rating DESC")
			normalized := strings.Replace(pair, ":", " ", 1)
			if orderRegex.MatchString(normalized) {
				parts := strings.Fields(normalized)
				if len(parts) == 2 {
					opts.Orders = append(opts.Orders, dto.OrderOption{
						Field:     parts[0],
						Direction: strings.ToUpper(parts[1]),
					})
				}
			}
		}
	}

	return opts
}

// BindAndValidate binds the JSON/Form body to the struct and executes validation tags.
func BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid JSON payload")
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	return nil
}

// BuildQueryOptions enriches QueryOptions with search, typed filters, IN filters, and date ranges.
// filterTypes: map[paramName]type ("int", "string", "bool", "decimal", "int_slice", "string_slice")
func BuildQueryOptions(c echo.Context, baseOpts dto.QueryOptions, filterTypes map[string]string, searchFields []string) dto.QueryOptions {
	// 1. Search Query
	searchQuery := c.QueryParam("search")
	if searchQuery != "" && len(searchFields) > 0 {
		baseOpts.Search = &dto.SearchOption{
			Query:  searchQuery,
			Fields: searchFields,
		}
	}

	// 2. Strict / Slice Filters
	for field, expectedType := range filterTypes {
		paramValue := c.QueryParam(field)
		if paramValue == "" {
			continue
		}

		switch expectedType {
		case "int":
			if intValue, err := strconv.ParseInt(paramValue, 10, 64); err == nil {
				baseOpts.Filters[field] = intValue
			}
		case "decimal":
			if decimalValue, err := decimal.NewFromString(paramValue); err == nil {
				baseOpts.Filters[field] = decimalValue
			}
		case "bool":
			switch strings.ToLower(paramValue) {
			case "true", "1":
				baseOpts.Filters[field] = true
			case "false", "0":
				baseOpts.Filters[field] = false
			}
		case "int_slice":
			intSlice := GetIntSliceQueryParam(c, field)
			if len(intSlice) > 0 {
				sliceAny := make([]any, len(intSlice))
				for i, v := range intSlice {
					sliceAny[i] = v
				}
				baseOpts.InFilters[field] = sliceAny
			}
		case "string_slice":
			strSlice := GetStringSliceQueryParam(c, field)
			if len(strSlice) > 0 {
				sliceAny := make([]any, len(strSlice))
				for i, v := range strSlice {
					sliceAny[i] = v
				}
				baseOpts.InFilters[field] = sliceAny
			}
		default:
			baseOpts.Filters[field] = paramValue
		}
	}

	// 3. Universal Date Range (?from_date=... & ?to_date=...)
	fromDate := c.QueryParam("from_date")
	toDate := c.QueryParam("to_date")
	dateField := c.QueryParam("date_field")
	if dateField == "" {
		dateField = "created_at"
	}

	if fromDate != "" || toDate != "" {
		from, to, err := ParseDateRangeUTC(fromDate, toDate)
		if err == nil {
			baseOpts.DateRanges = append(baseOpts.DateRanges, dto.DateRangeFilter{
				Field: dateField,
				From:  from,
				To:    to,
			})
		}
	}

	return baseOpts
}
