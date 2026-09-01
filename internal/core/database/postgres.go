package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"new_project/internal/core/config"
)

// NewPostgresDB creates and configures a new PostgreSQL connection using Bun ORM.
func NewPostgresDB(cfg config.DatabaseConfig) (*bun.DB, error) {
	safePassword := url.QueryEscape(cfg.Password)

	// Construct Data Source Name (DSN)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, safePassword, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	// Create a generic SQL database object
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Configure connection pool settings for high performance
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)                                           // Maximum number of open connections to the database
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)                                           // Maximum number of connections in the idle connection pool
	sqldb.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute) // Maximum amount of time a connection may be reused
	sqldb.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTimeMinutes) * time.Minute) // Maximum amount of time a connection may be idle

	// Initialize Bun database client
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query logging hook for debugging
	// It prints all executed SQL queries to the terminal beautifully
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Verify that the connection is actually established
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	return db, nil
}
