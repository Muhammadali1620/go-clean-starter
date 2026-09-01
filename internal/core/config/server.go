package config

type ServerConfig struct {
	Port         string
	MaxBodySize  string
	MaxRateLimit int
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         GetEnv("PORT", "8080"),
		MaxBodySize:  GetEnv("MAX_BODY_SIZE", "10M"),
		MaxRateLimit: GetEnvAsInt("MAX_RATE_LIMIT", 60),
	}
}
