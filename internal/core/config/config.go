package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	Server   ServerConfig
	CORS     CORSConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Bot      BotConfig
	Dev      DevConfig
}

// Load reads .env and aggregates all modular configs.
func Load() *Config {
	if os.Getenv("CONTAINER_MODE") != "1" {
		_ = godotenv.Load("env/.env.local")
		_ = godotenv.Load("env/.env")
	} else {
		log.Println("Running in container mode")
	}

	return &Config{
		AppEnv:   GetEnv("APP_ENV", "development"),
		Server:   loadServerConfig(),
		CORS:     loadCORSConfig(),
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		JWT:      loadJWTConfig(),
		Bot:      loadBotConfig(),
		Dev:      loadDevConfig(),
	}
}
