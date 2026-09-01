package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"new_project/internal/core/config"
)

// NewPostgresDB creates and configures a new PostgreSQL connection with retry backoff.
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

	// Logging hook for development
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Robust Retry Connection Loop (Waiting for DB container to warm up)
	var pingErr error
	maxAttempts := 15
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		pingErr = db.PingContext(ctx)
		cancel()

		if pingErr == nil {
			log.Println("🐘 PostgreSQL is ready and accepting connections!")
			return db, nil
		}

		log.Printf("⏳ Waiting for PostgreSQL to be ready... (attempt %d/%d)", attempt, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to postgres after %d attempts: %w", maxAttempts, pingErr)
}
