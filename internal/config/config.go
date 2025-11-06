package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
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

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
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

// Load loads configuration from YAML file
func Load() (*Config, error) {
	// Get current working directory
	currentDir := getCurrentDir()
	log.Printf("Current working directory: %s", currentDir)

	// Setup Viper
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add multiple search paths
	v.AddConfigPath(".")                  // Current directory
	v.AddConfigPath("./configs")          // configs directory
	v.AddConfigPath("../../")             // From cmd/server/
	v.AddConfigPath("../../../")          // Extra fallback
	v.AddConfigPath("/etc/iam-services/") // System config (Linux)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		log.Printf("WARNING: Failed to read config file: %v", err)
		log.Printf("Will use default values")
		// Don't return error - will use defaults
	} else {
		log.Printf("SUCCESS: Config file loaded from: %s", v.ConfigFileUsed())
	}

	// Allow environment variables to override config file
	v.AutomaticEnv()
	v.SetEnvPrefix("IAM") // will look for env vars like IAM_SERVER_HOST

	// Set defaults FIRST
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "50051")
	v.SetDefault("server.http_host", "0.0.0.0")
	v.SetDefault("server.http_port", "8080")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.dbname", "iam_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.secret", "your-secret-key-change-this-in-production")
	v.SetDefault("jwt.access_token_expiration_hours", 24)
	v.SetDefault("jwt.refresh_token_expiration_hours", 168)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.encoding", "json")
	v.SetDefault("swagger.enabled", true)
	v.SetDefault("swagger.base_path", "/swagger/")
	v.SetDefault("swagger.spec_path", "/swagger.json")
	v.SetDefault("swagger.title", "IAM Service API Documentation")
	v.SetDefault("swagger.auth.username", "admin")
	v.SetDefault("swagger.auth.password", "changeme")
	v.SetDefault("swagger.auth.realm", "IAM Service API Documentation")

	// Parse configuration into struct (matching config.yml structure)
	config := &Config{
		Server: ServerConfig{
			Host:     v.GetString("server.host"),
			Port:     v.GetString("server.port"),
			HTTPHost: v.GetString("server.http_host"), // Match config.yml: http_host
			HTTPPort: v.GetString("server.http_port"), // Match config.yml: http_port
		},
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetString("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			DBName:   v.GetString("database.dbname"),
			SSLMode:  v.GetString("database.sslmode"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("redis.host"),
			Port:     v.GetString("redis.port"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:               v.GetString("jwt.secret"),
			AccessTokenDuration:  time.Duration(v.GetInt("jwt.access_token_expiration_hours")) * time.Hour,
			RefreshTokenDuration: time.Duration(v.GetInt("jwt.refresh_token_expiration_hours")) * time.Hour,
		},
		Log: LogConfig{
			Level:    v.GetString("log.level"),
			Encoding: v.GetString("log.encoding"),
		},
		Swagger: SwaggerConfig{
			Enabled:      v.GetBool("swagger.enabled"),
			BasePath:     v.GetString("swagger.base_path"),
			SpecPath:     v.GetString("swagger.spec_path"),
			Title:        v.GetString("swagger.title"),
			AuthUsername: v.GetString("swagger.auth.username"),
			AuthPassword: v.GetString("swagger.auth.password"),
			AuthRealm:    v.GetString("swagger.auth.realm"),
		},
	}

	// CRITICAL: Ensure HTTP host and port are never empty
	if config.Server.HTTPHost == "" {
		log.Printf("WARNING: server.http.host not set, using default: 0.0.0.0")
		config.Server.HTTPHost = "0.0.0.0"
	}
	if config.Server.HTTPPort == "" {
		log.Printf("WARNING: server.http.port not set, using default: 8080")
		config.Server.HTTPPort = "8080"
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == "" {
		config.Server.Port = "50051"
	}

	// Log final configuration
	log.Printf("✓ Server Configuration:")
	log.Printf("  - gRPC: %s:%s", config.Server.Host, config.Server.Port)
	log.Printf("  - HTTP: %s:%s", config.Server.HTTPHost, config.Server.HTTPPort)
	log.Printf("✓ Database: %s@%s:%s/%s", config.Database.User, config.Database.Host, config.Database.Port, config.Database.DBName)
	log.Printf("✓ Redis: %s:%s (DB:%d)", config.Redis.Host, config.Redis.Port, config.Redis.DB)
	log.Printf("✓ Log Level: %s", config.Log.Level)

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

// GetAddress returns the Redis server address
func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getCurrentDir returns the current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
