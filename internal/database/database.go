package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Connect establishes a connection to the database
func Connect(dsn string, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully")
	return db, nil
}

// Close closes the database connection
func Close(db *sql.DB, logger *zap.Logger) {
	if err := db.Close(); err != nil {
		logger.Error("Failed to close database connection", zap.Error(err))
	} else {
		logger.Info("Database connection closed")
	}
}
