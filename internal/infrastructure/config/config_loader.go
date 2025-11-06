package config

import (
	"fmt"

	"github.com/tvttt/iam-services/internal/config"
)

// LoadConfig loads configuration from config.yml using Viper
func LoadConfig() (*config.Config, error) {
	// Use the main config loader which supports YAML and environment variables
	return config.Load()
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
