package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"log"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	Swagger  SwaggerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host     string
	Port     string
	HTTPHost string
	HTTPPort string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level    string
	Encoding string
}

// SwaggerConfig holds Swagger UI configuration
type SwaggerConfig struct {
	Enabled      bool
	BasePath     string
	SpecPath     string
	Title        string
	AuthUsername string
	AuthPassword string
	AuthRealm    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Host:     getEnv("SERVER_HOST", "0.0.0.0"),
			Port:     getEnv("SERVER_PORT", "50051"),
			HTTPHost: getEnv("HTTP_HOST", "0.0.0.0"),
			HTTPPort: getEnv("HTTP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "iam_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
			AccessTokenDuration:  getDurationEnv("JWT_EXPIRATION_HOURS", 24) * time.Hour,
			RefreshTokenDuration: getDurationEnv("JWT_REFRESH_EXPIRATION_HOURS", 168) * time.Hour,
		},
		Log: LogConfig{
			Level:    getEnv("LOG_LEVEL", "info"),
			Encoding: getEnv("LOG_ENCODING", "json"),
		},
		Swagger: SwaggerConfig{
			Enabled:      getBoolEnv("SWAGGER_ENABLED", true),
			BasePath:     getEnv("SWAGGER_BASE_PATH", "/swagger/"),
			SpecPath:     getEnv("SWAGGER_SPEC_PATH", "/swagger.json"),
			Title:        getEnv("SWAGGER_TITLE", "IAM Service API Documentation"),
			AuthUsername: getEnv("SWAGGER_AUTH_USERNAME", "admin"),
			AuthPassword: getEnv("SWAGGER_AUTH_PASSWORD", "changeme"),
			AuthRealm:    getEnv("SWAGGER_AUTH_REALM", "IAM Service API Documentation"),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetServerAddress returns the gRPC server address
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// GetHTTPServerAddress returns the HTTP gateway server address
func (c *ServerConfig) GetHTTPServerAddress() string {
	return fmt.Sprintf("%s:%s", c.HTTPHost, c.HTTPPort)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getDurationEnv(key string, defaultValue int) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return time.Duration(defaultValue)
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.Duration(defaultValue)
	}
	return time.Duration(value)
}

func getBoolEnv(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
