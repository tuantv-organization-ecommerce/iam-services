package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tvttt/iam-services/internal/config"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*config.Config, error) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "50051"),
		},
		Database: config.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "iam_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: config.RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: config.JWTConfig{
			Secret:               getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenDuration:  parseDuration(getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"), 15*time.Minute),
			RefreshTokenDuration: parseDuration(getEnv("JWT_REFRESH_TOKEN_DURATION", "168h"), 168*time.Hour),
		},
		Log: config.LogConfig{
			Level:    getEnv("LOG_LEVEL", "info"),
			Encoding: getEnv("LOG_ENCODING", "json"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string or returns a default value
func parseDuration(value string, defaultValue time.Duration) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable or returns a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// ValidateConfig validates the configuration
func ValidateConfig(cfg *config.Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if cfg.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if cfg.JWT.Secret == "your-secret-key" {
		fmt.Println("WARNING: Using default JWT secret. Please change it in production!")
	}

	return nil
}
