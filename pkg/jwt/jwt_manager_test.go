package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID   = "user-123"
	testUsername = "testuser"
)

func TestNewJWTManager(t *testing.T) {
	secret := "test-secret-key-min-32-chars-long"
	accessDuration := time.Hour
	refreshDuration := time.Hour * 24

	manager := NewJWTManager(secret, accessDuration, refreshDuration)

	assert.NotNil(t, manager)
	assert.Equal(t, secret, manager.SecretKey)
	assert.Equal(t, accessDuration, manager.AccessTokenDuration)
	assert.Equal(t, refreshDuration, manager.RefreshTokenDuration)
}

func TestGenerateAccessToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-chars-long", time.Hour, time.Hour*24)

	userID := testUserID
	username := testUsername
	roles := []string{"user", "admin"}

	token, err := manager.GenerateAccessToken(userID, username, roles)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-chars-long", time.Hour, time.Hour*24)

	userID := testUserID

	token, err := manager.GenerateRefreshToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken_Success(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-chars-long", time.Hour, time.Hour*24)

	userID := testUserID
	username := testUsername
	roles := []string{"user", "admin"}

	// Generate token
	token, err := manager.GenerateAccessToken(userID, username, roles)
	require.NoError(t, err)

	// Verify token
	claims, err := manager.VerifyToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, roles, claims.Roles)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-chars-long", time.Hour, time.Hour*24)

	// Try to verify invalid token
	_, err := manager.VerifyToken("invalid.token.here")
	assert.Error(t, err)
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	// Create manager with very short expiration
	manager := NewJWTManager("test-secret-key-min-32-chars-long", time.Millisecond, time.Millisecond)

	userID := testUserID
	username := testUsername
	roles := []string{"user"}

	// Generate token
	token, err := manager.GenerateAccessToken(userID, username, roles)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 10)

	// Try to verify expired token
	_, err = manager.VerifyToken(token)
	assert.Error(t, err)
}

func TestVerifyToken_WrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret-key-1-min-32-chars-long", time.Hour, time.Hour*24)
	manager2 := NewJWTManager("secret-key-2-min-32-chars-long", time.Hour, time.Hour*24)

	// Generate token with manager1
	token, err := manager1.GenerateAccessToken(testUserID, testUsername, []string{"user"})
	require.NoError(t, err)

	// Try to verify with manager2 (different secret)
	_, err = manager2.VerifyToken(token)
	assert.Error(t, err)
}
