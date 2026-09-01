package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/uptrace/bun/migrate"

	"new_project/internal/core/config"
	"new_project/internal/core/database"
	"new_project/internal/core/utils"
)

func main() {
	// 1. Load configuration and connect to the database
	appConfig := config.Load()
	db, err := database.NewPostgresDB(appConfig.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 2. Discover SQL migrations from the "migrations" directory
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(os.DirFS("migrations")); err != nil {
		log.Fatalf("Failed to discover migrations: %v", err)
	}

	// 3. Initialize the migrator
	migrator := migrate.NewMigrator(db, migrations)
	ctx := context.Background()

	// Ensure the bun_migrations table exists in the database
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	// 4. Parse CLI commands
	if len(os.Args) < 2 {
		log.Fatal("Please provide a command: up, down, status, or create")
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		if group.IsZero() {
			fmt.Println("Database is already up to date!")
		} else {
			fmt.Printf("Migrated %d files\n", len(group.Migrations))
		}

	case "down":
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		if group.IsZero() {
			fmt.Println("There are no applied migrations to rollback.")
		} else {
			fmt.Printf("Rolled back %d files\n", len(group.Migrations))
		}

	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		fmt.Printf("Migrations status:\n%s\n", ms)

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a name for the migration. Example: go run cmd/migrator/main.go create init_schema")
		}
		name := strings.Join(os.Args[2:], "_")

		// Generate timestamp in Bun's format (YYYYMMDDHHMMSS)
		timestamp := utils.NowUTC().Format("20060102150405")

		upFile := fmt.Sprintf("migrations/%s_%s.up.sql", timestamp, name)
		downFile := fmt.Sprintf("migrations/%s_%s.down.sql", timestamp, name)

		// Create empty files in the correct directory
		if err := os.WriteFile(upFile, []byte("-- Write your UP migration here\n"), 0644); err != nil {
			log.Fatalf("Failed to create up file: %v", err)
		}
		if err := os.WriteFile(downFile, []byte("-- Write your DOWN migration here\n"), 0644); err != nil {
			log.Fatalf("Failed to create down file: %v", err)
		}

		fmt.Printf("Created migration files in migrations/ directory:\n%s\n%s\n", upFile, downFile)

	default:
		log.Fatalf("Unknown command: %s. Use 'up', 'down', 'status', or 'create'", cmd)
	}
}
