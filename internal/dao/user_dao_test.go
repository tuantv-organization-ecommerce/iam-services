package dao

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tvttt/iam-services/internal/domain"
)

// getTestDB returns a test database connection
// It uses environment variables set by CI/CD or defaults for local testing
func getTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// Build DSN from individual env vars (for CI)
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		user := getEnvOrDefault("DB_USER", "postgres")
		password := getEnvOrDefault("DB_PASSWORD", "postgres")
		dbname := getEnvOrDefault("DB_NAME", "iam_db_test")
		sslmode := getEnvOrDefault("DB_SSL_MODE", "disable")

		dsn = "host=" + host + " port=" + port + " user=" + user +
			" password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
	}

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func TestUserDAO_Create_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userDAO.Create(ctx, user)
	require.NoError(t, err)

	// Cleanup
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
}

func TestUserDAO_FindByID_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userDAO.Create(ctx, user)
	require.NoError(t, err)

	// Get by ID
	retrieved, err := userDAO.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Username, retrieved.Username)
	assert.Equal(t, user.Email, retrieved.Email)

	// Cleanup
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
}

func TestUserDAO_FindByUsername_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userDAO.Create(ctx, user)
	require.NoError(t, err)

	// Get by username
	retrieved, err := userDAO.FindByUsername(ctx, user.Username)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Username, retrieved.Username)

	// Cleanup
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
}

func TestUserDAO_FindByEmail_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userDAO.Create(ctx, user)
	require.NoError(t, err)

	// Get by email
	retrieved, err := userDAO.FindByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Email, retrieved.Email)

	// Cleanup
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
}

func TestUserDAO_Update_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userDAO.Create(ctx, user)
	require.NoError(t, err)

	// Update user
	user.FullName = "Updated Name"
	user.IsActive = false
	user.UpdatedAt = time.Now()

	err = userDAO.Update(ctx, user)
	require.NoError(t, err)

	// Verify update
	retrieved, err := userDAO.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, "Updated Name", retrieved.FullName)
	assert.False(t, retrieved.IsActive)

	// Cleanup
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
}

func TestUserDAO_FindByID_NotFound(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	result, err := userDAO.FindByID(ctx, nonExistentID)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestUserDAO_Delete_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userDAO := NewUserDAO(db)
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userDAO.Create(ctx, user)
	require.NoError(t, err)

	// Delete user
	err = userDAO.Delete(ctx, user.ID)
	require.NoError(t, err)

	// Verify deletion
	result, err := userDAO.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Nil(t, result)
}
