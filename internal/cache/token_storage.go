package cache

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// TokenStorage manages token storage in Redis
type TokenStorage struct {
	redis  *RedisClient
	logger *zap.Logger
}

// NewTokenStorage creates a new TokenStorage instance
func NewTokenStorage(redis *RedisClient, logger *zap.Logger) *TokenStorage {
	return &TokenStorage{
		redis:  redis,
		logger: logger,
	}
}

// StoreAccessToken stores an access token with TTL
func (ts *TokenStorage) StoreAccessToken(ctx context.Context, userID, token string, expiration time.Duration) error {
	key := fmt.Sprintf("access_token:%s", userID)
	if err := ts.redis.Set(ctx, key, token, expiration); err != nil {
		ts.logger.Error("Failed to store access token",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to store access token: %w", err)
	}

	ts.logger.Debug("Access token stored successfully",
		zap.String("user_id", userID),
		zap.Duration("expiration", expiration))

	return nil
}

// StoreRefreshToken stores a refresh token with TTL
func (ts *TokenStorage) StoreRefreshToken(ctx context.Context, userID, token string, expiration time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", userID)
	if err := ts.redis.Set(ctx, key, token, expiration); err != nil {
		ts.logger.Error("Failed to store refresh token",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	ts.logger.Debug("Refresh token stored successfully",
		zap.String("user_id", userID),
		zap.Duration("expiration", expiration))

	return nil
}

// GetAccessToken retrieves an access token
func (ts *TokenStorage) GetAccessToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("access_token:%s", userID)
	token, err := ts.redis.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	return token, nil
}

// GetRefreshToken retrieves a refresh token
func (ts *TokenStorage) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("refresh_token:%s", userID)
	token, err := ts.redis.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	return token, nil
}

// RevokeAccessToken removes an access token from storage
func (ts *TokenStorage) RevokeAccessToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("access_token:%s", userID)
	if err := ts.redis.Delete(ctx, key); err != nil {
		ts.logger.Error("Failed to revoke access token",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to revoke access token: %w", err)
	}

	ts.logger.Info("Access token revoked successfully",
		zap.String("user_id", userID))

	return nil
}

// RevokeRefreshToken removes a refresh token from storage
func (ts *TokenStorage) RevokeRefreshToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("refresh_token:%s", userID)
	if err := ts.redis.Delete(ctx, key); err != nil {
		ts.logger.Error("Failed to revoke refresh token",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	ts.logger.Info("Refresh token revoked successfully",
		zap.String("user_id", userID))

	return nil
}

// RevokeAllTokens removes both access and refresh tokens for a user
func (ts *TokenStorage) RevokeAllTokens(ctx context.Context, userID string) error {
	accessKey := fmt.Sprintf("access_token:%s", userID)
	refreshKey := fmt.Sprintf("refresh_token:%s", userID)

	if err := ts.redis.Delete(ctx, accessKey, refreshKey); err != nil {
		ts.logger.Error("Failed to revoke all tokens",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to revoke all tokens: %w", err)
	}

	ts.logger.Info("All tokens revoked successfully",
		zap.String("user_id", userID))

	return nil
}

// IsAccessTokenValid checks if an access token exists and is valid
func (ts *TokenStorage) IsAccessTokenValid(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("access_token:%s", userID)
	exists, err := ts.redis.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check access token validity: %w", err)
	}
	return exists, nil
}

// IsRefreshTokenValid checks if a refresh token exists and is valid
func (ts *TokenStorage) IsRefreshTokenValid(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%s", userID)
	exists, err := ts.redis.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check refresh token validity: %w", err)
	}
	return exists, nil
}

// GetAccessTokenTTL returns the remaining time for an access token
func (ts *TokenStorage) GetAccessTokenTTL(ctx context.Context, userID string) (time.Duration, error) {
	key := fmt.Sprintf("access_token:%s", userID)
	ttl, err := ts.redis.TTL(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get access token TTL: %w", err)
	}
	return ttl, nil
}

// GetRefreshTokenTTL returns the remaining time for a refresh token
func (ts *TokenStorage) GetRefreshTokenTTL(ctx context.Context, userID string) (time.Duration, error) {
	key := fmt.Sprintf("refresh_token:%s", userID)
	ttl, err := ts.redis.TTL(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get refresh token TTL: %w", err)
	}
	return ttl, nil
}

