package config

type JWTConfig struct {
	Secret        string
	AccessExpiry  int
	RefreshExpiry int
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:        GetEnv("JWT_SECRET", "super_secret_jwt_key_change_in_production"),
		AccessExpiry:  GetEnvAsInt("JWT_ACCESS_EXPIRY_MINUTES", 60),
		RefreshExpiry: GetEnvAsInt("JWT_REFRESH_EXPIRY_DAYS", 30),
	}
}
