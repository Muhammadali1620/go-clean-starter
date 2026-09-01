package database

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type Transactor interface {
	RunInTx(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type txKey struct{}

func InjectTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func ExtractTx(ctx context.Context) (bun.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(bun.Tx)
	return tx, ok
}

type bunTransactor struct {
	db *bun.DB
}

// @inject
func NewBunTransactor(db *bun.DB) Transactor {
	return &bunTransactor{db: db}
}

func (t *bunTransactor) RunInTx(ctx context.Context, fn func(ctxTx context.Context) error) error {
	return t.db.RunInTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted}, func(ctx context.Context, tx bun.Tx) error {
		ctxWithTx := InjectTx(ctx, tx)
		return fn(ctxWithTx)
	})
}
