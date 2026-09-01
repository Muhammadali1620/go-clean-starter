package config

type DevConfig struct {
	BasicAuthUser     string
	BasicAuthPassword string
}

func loadDevConfig() DevConfig {
	return DevConfig{
		BasicAuthUser:     GetEnv("DEV_USER", "admin"),
		BasicAuthPassword: GetEnv("DEV_PASSWORD", "admin123"),
	}
}
