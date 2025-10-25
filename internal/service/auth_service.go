package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	userRepo    repository.UserRepository
	authzRepo   repository.AuthorizationRepository
	jwtManager  *jwt.JWTManager
	passwordMgr *password.PasswordManager
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	authzRepo repository.AuthorizationRepository,
	jwtManager *jwt.JWTManager,
	passwordMgr *password.PasswordManager,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		authzRepo:   authzRepo,
		jwtManager:  jwtManager,
		passwordMgr: passwordMgr,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password, fullName string) (*domain.User, error) {
	// Validate input
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email, and password are required")
	}

	// Hash password
	hashedPassword, err := s.passwordMgr.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (*domain.User, *domain.TokenPair, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if !s.passwordMgr.CheckPassword(password, user.PasswordHash) {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Get user roles
	roles, err := s.authzRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, roleNames)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenPair := &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtManager.AccessTokenDuration.Seconds()),
	}

	return user, tokenPair, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	// Verify refresh token
	claims, err := s.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Get user roles
	roles, err := s.authzRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenPair := &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtManager.AccessTokenDuration.Seconds()),
	}

	return tokenPair, nil
}

func (s *authService) VerifyToken(ctx context.Context, token string) (string, []string, error) {
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return "", nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims.UserID, claims.Roles, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	// In a production system, you would invalidate the token here
	// For now, we just return nil
	// You could store revoked tokens in Redis or a database
	return nil
}
