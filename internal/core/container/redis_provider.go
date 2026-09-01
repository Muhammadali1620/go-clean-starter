package container

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"new_project/internal/core/config"
	"new_project/internal/core/database"
)

// ProvideRedisClient creates Redis connection automatically from config.
func ProvideRedisClient(cfg *config.Config) *redis.Client {
	if cfg.Redis.Host == "" {
		return nil
	}
	rdb, err := database.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		fmt.Printf("⚠️ Redis skipped: %v\n", err)
		return nil
	}
	return rdb
}
