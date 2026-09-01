package config

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     GetEnv("REDIS_HOST", "localhost"),
		Port:     GetEnv("REDIS_PORT", "6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetEnvAsInt("REDIS_DB", 0),
	}
}
