package app

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"

	"new_project/internal/core/config"
	"new_project/internal/core/container"
	"new_project/internal/core/database"
)

type App struct {
	Config    *config.Config
	DB        *bun.DB
	RedisDB   *redis.Client
	RedisOpt  asynq.RedisClientOpt
	Container *container.Container
}

func Bootstrap() *App {
	cfg := config.Load()

	// 1. PostgreSQL (Required)
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// 2. Redis Options (Optional)
	var redisOpt asynq.RedisClientOpt
	if cfg.Redis.Host != "" {
		redisOpt = asynq.RedisClientOpt{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}
	}

	// 3. DI Container
	appContainer, err := container.InitializeContainer(db, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize DI Container: %v", err)
	}

	// 4. Register Telegram Bot (if token exists)
	if appContainer.BotService != nil && appContainer.BotHandlers.Start != nil {
		appContainer.BotService.RegisterAllRoutes(appContainer.BotHandlers.Start)
	}

	return &App{
		Config:    cfg,
		DB:        db,
		RedisOpt:  redisOpt,
		Container: appContainer,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		_ = a.DB.Close()
	}
}
