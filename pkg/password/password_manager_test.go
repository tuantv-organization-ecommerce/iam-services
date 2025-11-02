package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPassword = "MySecurePassword123!"
)

func TestNewPasswordManager(t *testing.T) {
	manager := NewPasswordManager()
	assert.NotNil(t, manager)
}

func TestHashPassword_Success(t *testing.T) {
	manager := NewPasswordManager()
	password := testPassword

	hash, err := manager.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.True(t, len(hash) > 50) // bcrypt hashes are typically 60 chars
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	manager := NewPasswordManager()

	hash, err := manager.HashPassword("")

	// Should still work but hash will be of empty string
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPassword_Success(t *testing.T) {
	manager := NewPasswordManager()
	password := testPassword

	// Hash the password
	hash, err := manager.HashPassword(password)
	require.NoError(t, err)

	// Check with correct password
	result := manager.CheckPassword(password, hash)
	assert.True(t, result)
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	manager := NewPasswordManager()
	password := testPassword
	wrongPassword := "WrongPassword456!"

	// Hash the password
	hash, err := manager.HashPassword(password)
	require.NoError(t, err)

	// Check with wrong password
	result := manager.CheckPassword(wrongPassword, hash)
	assert.False(t, result)
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	manager := NewPasswordManager()
	password := testPassword
	invalidHash := "invalid-hash-string"

	// Check with invalid hash
	result := manager.CheckPassword(password, invalidHash)
	assert.False(t, result)
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	manager := NewPasswordManager()
	password := testPassword

	// Generate two hashes for the same password
	hash1, err := manager.HashPassword(password)
	require.NoError(t, err)

	hash2, err := manager.HashPassword(password)
	require.NoError(t, err)

	// Hashes should be different (bcrypt uses random salt)
	assert.NotEqual(t, hash1, hash2)

	// But both should verify correctly
	assert.True(t, manager.CheckPassword(password, hash1))
	assert.True(t, manager.CheckPassword(password, hash2))
}

func TestCheckPassword_CaseSensitive(t *testing.T) {
	manager := NewPasswordManager()
	password := testPassword

	hash, err := manager.HashPassword(password)
	require.NoError(t, err)

	// Check with different case
	assert.False(t, manager.CheckPassword("mysecurepassword123!", hash))
	assert.False(t, manager.CheckPassword("MYSECUREPASSWORD123!", hash))
	assert.True(t, manager.CheckPassword(password, hash))
}
