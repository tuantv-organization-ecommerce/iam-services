package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tvttt/iam-services/internal/cache"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/repository"
	"github.com/tvttt/iam-services/pkg/jwt"
	"github.com/tvttt/iam-services/pkg/password"
)

// AuthService handles authentication business logic
type AuthService interface {
	Register(ctx context.Context, username, email, password, fullName string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, *domain.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	VerifyToken(ctx context.Context, token string) (string, []string, error)
	Logout(ctx context.Context, userID string) error
}

type authService struct {
	userRepo     repository.UserRepository
	authzRepo    repository.AuthorizationRepository
	jwtManager   *jwt.JWTManager
	passwordMgr  *password.PasswordManager
	tokenStorage *cache.TokenStorage
	logger       *zap.Logger
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	authzRepo repository.AuthorizationRepository,
	jwtManager *jwt.JWTManager,
	passwordMgr *password.PasswordManager,
	tokenStorage *cache.TokenStorage,
	logger *zap.Logger,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		authzRepo:    authzRepo,
		jwtManager:   jwtManager,
		passwordMgr:  passwordMgr,
		tokenStorage: tokenStorage,
		logger:       logger,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password, fullName string) (*domain.User, error) {
	// Validate input
	if username == "" || email == "" || password == "" {
		err := NewServiceError(ErrCodeInvalidInput, "username, email, and password are required", nil)
		LogError(s.logger, err, "Register", zap.String("username", username), zap.String("email", email))
		return nil, err
	}

	// Hash password
	hashedPassword, err := s.passwordMgr.HashPassword(password)
	if err != nil {
		serviceErr := NewServiceError(ErrCodePasswordHashFailed, "failed to hash password", err)
		LogError(s.logger, serviceErr, "Register", zap.String("username", username))
		return nil, serviceErr
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		serviceErr := NewServiceError(ErrCodeUserCreationFailed, "failed to create user", err)
		LogError(s.logger, serviceErr, "Register", zap.String("username", username), zap.String("email", email))
		return nil, serviceErr
	}

	s.logger.Info("User registered successfully",
		zap.String("user_id", user.ID),
		zap.String("username", username),
		zap.String("email", email))

	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (*domain.User, *domain.TokenPair, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeInvalidCredentials, "invalid credentials", err)
		LogError(s.logger, serviceErr, "Login", zap.String("username", username))
		return nil, nil, serviceErr
	}

	// Check if user is active
	if !user.IsActive {
		serviceErr := NewServiceError(ErrCodeUserInactive, "user account is inactive", nil)
		LogError(s.logger, serviceErr, "Login", zap.String("username", username), zap.String("user_id", user.ID))
		return nil, nil, serviceErr
	}

	// Verify password
	if !s.passwordMgr.CheckPassword(password, user.PasswordHash) {
		serviceErr := NewServiceError(ErrCodeInvalidCredentials, "invalid credentials", nil)
		LogError(s.logger, serviceErr, "Login", zap.String("username", username), zap.String("user_id", user.ID))
		return nil, nil, serviceErr
	}

	// Get user roles
	roles, err := s.authzRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeGetUserRolesFailed, "failed to get user roles", err)
		LogError(s.logger, serviceErr, "Login", zap.String("username", username), zap.String("user_id", user.ID))
		return nil, nil, serviceErr
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, roleNames)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeTokenGenerationFailed, "failed to generate access token", err)
		LogError(s.logger, serviceErr, "Login", zap.String("username", username), zap.String("user_id", user.ID))
		return nil, nil, serviceErr
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeTokenGenerationFailed, "failed to generate refresh token", err)
		LogError(s.logger, serviceErr, "Login", zap.String("username", username), zap.String("user_id", user.ID))
		return nil, nil, serviceErr
	}

	tokenPair := &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtManager.AccessTokenDuration.Seconds()),
	}

	// Store tokens in Redis with TTL
	if s.tokenStorage != nil {
		if err := s.tokenStorage.StoreAccessToken(ctx, user.ID, accessToken, s.jwtManager.AccessTokenDuration); err != nil {
			// Log error but don't fail the login - token storage is not critical
			s.logger.Warn("Failed to store access token in cache",
				zap.String("user_id", user.ID),
				zap.Error(err))
		}

		if err := s.tokenStorage.StoreRefreshToken(ctx, user.ID, refreshToken, s.jwtManager.RefreshTokenDuration); err != nil {
			// Log error but don't fail the login
			s.logger.Warn("Failed to store refresh token in cache",
				zap.String("user_id", user.ID),
				zap.Error(err))
		}
	}

	s.logger.Info("User logged in successfully",
		zap.String("user_id", user.ID),
		zap.String("username", username),
		zap.Strings("roles", roleNames))

	return user, tokenPair, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	// Verify refresh token
	claims, err := s.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeInvalidToken, "invalid refresh token", err)
		LogError(s.logger, serviceErr, "RefreshToken")
		return nil, serviceErr
	}

	// Get user
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeUserNotFound, "user not found", err)
		LogError(s.logger, serviceErr, "RefreshToken", zap.String("user_id", claims.UserID))
		return nil, serviceErr
	}

	// Check if user is active
	if !user.IsActive {
		serviceErr := NewServiceError(ErrCodeUserInactive, "user account is inactive", nil)
		LogError(s.logger, serviceErr, "RefreshToken", zap.String("user_id", user.ID))
		return nil, serviceErr
	}

	// Get user roles
	roles, err := s.authzRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeGetUserRolesFailed, "failed to get user roles", err)
		LogError(s.logger, serviceErr, "RefreshToken", zap.String("user_id", user.ID))
		return nil, serviceErr
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, roleNames)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeTokenGenerationFailed, "failed to generate access token", err)
		LogError(s.logger, serviceErr, "RefreshToken", zap.String("user_id", user.ID))
		return nil, serviceErr
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeTokenGenerationFailed, "failed to generate refresh token", err)
		LogError(s.logger, serviceErr, "RefreshToken", zap.String("user_id", user.ID))
		return nil, serviceErr
	}

	tokenPair := &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtManager.AccessTokenDuration.Seconds()),
	}

	// Revoke old tokens and store new ones in Redis
	if s.tokenStorage != nil {
		// Revoke old tokens
		if err := s.tokenStorage.RevokeAllTokens(ctx, user.ID); err != nil {
			// Log error but don't fail the refresh
			s.logger.Warn("Failed to revoke old tokens from cache",
				zap.String("user_id", user.ID),
				zap.Error(err))
		}

		// Store new tokens with TTL
		if err := s.tokenStorage.StoreAccessToken(ctx, user.ID, accessToken, s.jwtManager.AccessTokenDuration); err != nil {
			// Log error but don't fail the refresh
			s.logger.Warn("Failed to store access token in cache",
				zap.String("user_id", user.ID),
				zap.Error(err))
		}

		if err := s.tokenStorage.StoreRefreshToken(ctx, user.ID, newRefreshToken, s.jwtManager.RefreshTokenDuration); err != nil {
			// Log error but don't fail the refresh
			s.logger.Warn("Failed to store refresh token in cache",
				zap.String("user_id", user.ID),
				zap.Error(err))
		}
	}

	s.logger.Info("Token refreshed successfully",
		zap.String("user_id", user.ID),
		zap.Strings("roles", roleNames))

	return tokenPair, nil
}

func (s *authService) VerifyToken(ctx context.Context, token string) (string, []string, error) {
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeInvalidToken, "invalid token", err)
		LogError(s.logger, serviceErr, "VerifyToken")
		return "", nil, serviceErr
	}

	return claims.UserID, claims.Roles, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	// Revoke all tokens from Redis
	if s.tokenStorage != nil {
		if err := s.tokenStorage.RevokeAllTokens(ctx, userID); err != nil {
			serviceErr := NewServiceError(ErrCodeTokenRevocationFailed, "failed to revoke tokens", err)
			LogError(s.logger, serviceErr, "Logout", zap.String("user_id", userID))
			return serviceErr
		}
	}

	s.logger.Info("User logged out successfully", zap.String("user_id", userID))
	return nil
}
