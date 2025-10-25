package service

import "errors"

var (
	ErrPasswordHashingFailed      = errors.New("password hashing failed")
	ErrPasswordVerificationFailed = errors.New("password verification failed")
)

// PasswordService defines the contract for password operations
// This is a domain service interface
type PasswordService interface {
	// Hash hashes a plain text password
	Hash(password string) (string, error)

	// Verify verifies a password against a hash
	Verify(password, hash string) bool
}
