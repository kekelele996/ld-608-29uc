package config

import (
	"os"
)

// Config is read from env in main.go and injected into services.
type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	JWTSecret    string
	TokenTTLHour string
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Load reads the scattered deployment config (docker-compose -> env -> defaults).
func Load() Config {
	return Config{
		Port:         getenv("PORT", "3000"),
		DBHost:       getenv("DB_HOST", "db"),
		DBPort:       getenv("DB_PORT", "3306"),
		DBName:       getenv("DB_NAME", "app_db"),
		DBUser:       getenv("DB_USER", "app_user"),
		DBPassword:   getenv("DB_PASSWORD", "app_password"),
		JWTSecret:    getenv("JWT_SECRET", "local-dev-secret"),
		TokenTTLHour: getenv("JWT_TTL_HOURS", "12"),
	}
}

// DSN builds the MySQL DSN consumed by GORM.
func (c Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=UTC"
}
