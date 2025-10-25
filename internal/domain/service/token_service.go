package service

import (
	"errors"
	"time"
)

var (
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token expired")
)

// TokenClaims represents token claims
type TokenClaims struct {
	UserID    string
	Username  string
	Roles     []string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}

// TokenService defines the contract for token operations
// This is a domain service interface
type TokenService interface {
	// GenerateAccessToken generates an access token
	GenerateAccessToken(userID, username string, roles []string) (string, error)

	// GenerateRefreshToken generates a refresh token
	GenerateRefreshToken(userID string) (string, error)

	// GenerateTokenPair generates both access and refresh tokens
	GenerateTokenPair(userID, username string, roles []string) (*TokenPair, error)

	// VerifyToken verifies and parses a token
	VerifyToken(token string) (*TokenClaims, error)

	// RefreshAccessToken generates a new access token from refresh token
	RefreshAccessToken(refreshToken string) (string, error)
}
