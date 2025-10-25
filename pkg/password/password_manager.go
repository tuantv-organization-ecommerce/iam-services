package password

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordManager handles password hashing and verification
type PasswordManager struct {
	Cost int
}

// NewPasswordManager creates a new password manager
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		Cost: bcrypt.DefaultCost,
	}
}

// HashPassword hashes a plain text password
func (m *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), m.Cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword verifies a password against a hash
func (m *PasswordManager) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
