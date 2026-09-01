package repositories

import (
	"context"
	"regexp"
	"strings"

	"github.com/uptrace/bun"

	"new_project/internal/core/database"
	"new_project/internal/dto"
)

// Safe regex for column names: allows letters, numbers, underscores, and jsonb arrows (e.g. name->>'uz')
var validIdentRegex = regexp.MustCompile(`^[a-zA-Z0-9_\.]+(?:->>'?[a-zA-Z0-9_]+'?)?$`)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int64) error
	DeleteBy(ctx context.Context, field string, value any) error
	FindByID(ctx context.Context, id int64, relations ...string) (*T, error)
	FindOneBy(ctx context.Context, field string, value any, relations ...string) (*T, error)
	FindAll(ctx context.Context, opts dto.QueryOptions) ([]T, int, error)
	FindByIDForUpdate(ctx context.Context, id int64, relations ...string) (*T, error)
	FindOneByForUpdate(ctx context.Context, field string, value any, relations ...string) (*T, error)
	Conn(ctx context.Context) bun.IDB
	BuildQuery(ctx context.Context, query *bun.SelectQuery, opts dto.QueryOptions) *bun.SelectQuery
}

type baseRepository[T any] struct {
	db *bun.DB
}

func NewBaseRepository[T any](db *bun.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

// Conn extracts the active transaction from context if present, or returns the base DB client.
func (r *baseRepository[T]) Conn(ctx context.Context) bun.IDB {
	if tx, ok := database.ExtractTx(ctx); ok {
		return tx
	}
	return r.db
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	_, err := r.Conn(ctx).NewInsert().Model(entity).Exec(ctx)
	return database.MapDBError(err)
}

func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	_, err := r.Conn(ctx).NewUpdate().Model(entity).WherePK().Exec(ctx)
	return database.MapDBError(err)
}

func (r *baseRepository[T]) Delete(ctx context.Context, id int64) error {
	_, err := r.Conn(ctx).NewDelete().Model((*T)(nil)).Where("?TableAlias.id = ?", id).Exec(ctx)
	return database.MapDBError(err)
}

func (r *baseRepository[T]) DeleteBy(ctx context.Context, field string, value any) error {
	if !validIdentRegex.MatchString(field) {
		return nil
	}
	_, err := r.Conn(ctx).NewDelete().Model((*T)(nil)).Where("?TableAlias.? = ?", bun.Ident(field), value).Exec(ctx)
	return database.MapDBError(err)
}

func (r *baseRepository[T]) FindByID(ctx context.Context, id int64, relations ...string) (*T, error) {
	return r.FindOneBy(ctx, "id", id, relations...)
}

func (r *baseRepository[T]) FindByIDForUpdate(ctx context.Context, id int64, relations ...string) (*T, error) {
	var entity T
	query := r.Conn(ctx).NewSelect().Model(&entity).Where("?TableAlias.id = ?", id).For("UPDATE")

	for _, rel := range relations {
		query.Relation(rel)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, database.MapDBError(err)
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindOneBy(ctx context.Context, field string, value any, relations ...string) (*T, error) {
	var entity T
	query := r.Conn(ctx).NewSelect().Model(&entity)

	if strings.Contains(field, "->>") {
		query.Where("?TableAlias.? = ?", bun.Safe(field), value)
	} else {
		query.Where("?TableAlias.? = ?", bun.Ident(field), value)
	}

	for _, rel := range relations {
		query.Relation(rel)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, database.MapDBError(err)
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindOneByForUpdate(ctx context.Context, field string, value any, relations ...string) (*T, error) {
	var entity T
	query := r.Conn(ctx).NewSelect().Model(&entity).For("UPDATE")

	if strings.Contains(field, "->>") {
		query.Where("?TableAlias.? = ?", bun.Safe(field), value)
	} else {
		query.Where("?TableAlias.? = ?", bun.Ident(field), value)
	}

	for _, rel := range relations {
		query.Relation(rel)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, database.MapDBError(err)
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindAll(ctx context.Context, opts dto.QueryOptions) ([]T, int, error) {
	// Initialize empty slice so scanning 0 items never results in nil JSON
	entities := make([]T, 0)
	query := r.Conn(ctx).NewSelect().Model(&entities)

	query = r.BuildQuery(ctx, query, opts)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, database.MapDBError(err)
	}

	return entities, count, nil
}

// BuildQuery applies equality filters, IN, NOT IN, date ranges, search, ordering, relations, and pagination.
func (r *baseRepository[T]) BuildQuery(ctx context.Context, query *bun.SelectQuery, opts dto.QueryOptions) *bun.SelectQuery {
	// 1. Strict Equality Filters: WHERE field = value
	for field, value := range opts.Filters {
		if !validIdentRegex.MatchString(field) {
			continue
		}
		if strings.Contains(field, ".") {
			if strings.Contains(field, "->>") {
				query.Where("? = ?", bun.Safe(field), value)
			} else {
				query.Where("? = ?", bun.Ident(field), value)
			}
		} else {
			if strings.Contains(field, "->>") {
				query.Where("?TableAlias.? = ?", bun.Safe(field), value)
			} else {
				query.Where("?TableAlias.? = ?", bun.Ident(field), value)
			}
		}
	}

	// 2. IN Filters: WHERE field IN (...)
	for field, values := range opts.InFilters {
		if len(values) == 0 || !validIdentRegex.MatchString(field) {
			continue
		}
		if strings.Contains(field, ".") {
			query.Where("? IN (?)", bun.Ident(field), bun.In(values))
		} else {
			query.Where("?TableAlias.? IN (?)", bun.Ident(field), bun.In(values))
		}
	}

	// 3. NOT IN Filters (Exclude): WHERE field NOT IN (...)
	for field, values := range opts.NotInFilters {
		if len(values) == 0 || !validIdentRegex.MatchString(field) {
			continue
		}
		if strings.Contains(field, ".") {
			query.Where("? NOT IN (?)", bun.Ident(field), bun.In(values))
		} else {
			query.Where("?TableAlias.? NOT IN (?)", bun.Ident(field), bun.In(values))
		}
	}

	// 4. Date Ranges: WHERE field >= from AND field <= to
	for _, dr := range opts.DateRanges {
		if dr.Field == "" || !validIdentRegex.MatchString(dr.Field) {
			continue
		}
		if dr.From != nil {
			if strings.Contains(dr.Field, ".") {
				query.Where("? >= ?", bun.Ident(dr.Field), *dr.From)
			} else {
				query.Where("?TableAlias.? >= ?", bun.Ident(dr.Field), *dr.From)
			}
		}
		if dr.To != nil {
			if strings.Contains(dr.Field, ".") {
				query.Where("? <= ?", bun.Ident(dr.Field), *dr.To)
			} else {
				query.Where("?TableAlias.? <= ?", bun.Ident(dr.Field), *dr.To)
			}
		}
	}

	// 5. Search (Case-Insensitive OR inside AND group)
	if opts.Search != nil && opts.Search.Query != "" && len(opts.Search.Fields) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			cleanSearch := escapeSQLWildcards(opts.Search.Query)
			searchTerm := "%" + cleanSearch + "%"
			for _, field := range opts.Search.Fields {
				if !validIdentRegex.MatchString(field) {
					continue
				}
				if strings.Contains(field, ".") {
					if strings.Contains(field, "->>") {
						q.WhereOr("? ILIKE ?", bun.Safe(field), searchTerm)
					} else {
						q.WhereOr("? ILIKE ?", bun.Ident(field), searchTerm)
					}
				} else {
					if strings.Contains(field, "->>") {
						q.WhereOr("?TableAlias.? ILIKE ?", bun.Safe(field), searchTerm)
					} else {
						q.WhereOr("?TableAlias.? ILIKE ?", bun.Ident(field), searchTerm)
					}
				}
			}
			return q
		})
	}

	// 6. Preload Relations
	for _, rel := range opts.Relations {
		query.Relation(rel)
	}

	// 7. Multi-Order Sorting
	hasIDOrder := false
	if len(opts.Orders) > 0 {
		for _, ord := range opts.Orders {
			if !validIdentRegex.MatchString(ord.Field) {
				continue
			}
			dir := "ASC"
			if strings.EqualFold(ord.Direction, "DESC") {
				dir = "DESC"
			}

			if strings.EqualFold(ord.Field, "id") || strings.HasSuffix(ord.Field, ".id") {
				hasIDOrder = true
			}

			if strings.Contains(ord.Field, ".") {
				query.OrderExpr("? "+dir, bun.Ident(ord.Field))
			} else {
				query.OrderExpr("?TableAlias.? "+dir, bun.Ident(ord.Field))
			}
		}
	}

	// Deterministic sort fallback to prevent pagination duplicate glitches
	if !hasIDOrder {
		query.OrderExpr("?TableAlias.created_at DESC, ?TableAlias.id DESC")
	}

	// 8. Pagination
	if opts.Limit > 0 {
		query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query.Offset(opts.Offset)
	}

	return query
}

func escapeSQLWildcards(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
