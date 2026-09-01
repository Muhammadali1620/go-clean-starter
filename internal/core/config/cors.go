package config

type CORSConfig struct {
	AllowOrigins []string
}

func loadCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: GetEnvAsSlice("CORS_ALLOW_ORIGINS", []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:8080",
		}),
	}
}
