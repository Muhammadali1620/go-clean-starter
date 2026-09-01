package config

type DatabaseConfig struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	Name                   string
	SSLMode                string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetimeMinutes int
	ConnMaxIdleTimeMinutes int
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:                   GetEnv("DB_HOST", "localhost"),
		Port:                   GetEnv("DB_PORT", "5432"),
		User:                   GetEnv("POSTGRES_USER", "postgres"),
		Password:               GetEnv("POSTGRES_PASSWORD", "postgres"),
		Name:                   GetEnv("POSTGRES_DB", "app_db"),
		SSLMode:                GetEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:           GetEnvAsInt("DB_MAX_OPEN_CONNS", 50),
		MaxIdleConns:           GetEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetimeMinutes: GetEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 60),
		ConnMaxIdleTimeMinutes: GetEnvAsInt("DB_CONN_MAX_IDLE_TIME_MINUTES", 15),
	}
}
