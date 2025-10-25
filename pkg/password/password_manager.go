package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Custom errors
var (
	ErrPasswordTooShort = errors.New("password is too short (minimum 8 characters)")
	ErrPasswordTooLong  = errors.New("password is too long (maximum 72 characters)")
	ErrPasswordMismatch = errors.New("password does not match")
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

// VerifyPassword verifies a password against a hash and returns an error if it doesn't match
func (m *PasswordManager) VerifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return err
	}
	return nil
}
