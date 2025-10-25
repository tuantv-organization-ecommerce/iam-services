package auth

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/application/dto"
	"github.com/tvttt/iam-services/internal/domain/repository"
	"github.com/tvttt/iam-services/internal/domain/service"
)

// LoginUseCase handles user login business logic
type LoginUseCase struct {
	userRepo    repository.UserRepository
	authzRepo   repository.AuthorizationRepository
	passwordSvc service.PasswordService
	tokenSvc    service.TokenService
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(
	userRepo repository.UserRepository,
	authzRepo repository.AuthorizationRepository,
	passwordSvc service.PasswordService,
	tokenSvc service.TokenService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		authzRepo:   authzRepo,
		passwordSvc: passwordSvc,
		tokenSvc:    tokenSvc,
	}
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Validate input
	if err := uc.validateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Find user by username
	user, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Check if user is active
	if !user.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	// 4. Verify password
	if !uc.passwordSvc.Verify(req.Password, user.PasswordHash()) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 5. Get user roles
	roles, err := uc.authzRepo.GetUserRoles(ctx, user.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name()
	}

	// 6. Generate tokens
	tokenPair, err := uc.tokenSvc.GenerateTokenPair(user.ID(), user.Username(), roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// 7. Return response
	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &dto.UserDTO{
			ID:        user.ID(),
			Username:  user.Username(),
			Email:     user.Email(),
			FullName:  user.FullName(),
			IsActive:  user.IsActive(),
			CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func (uc *LoginUseCase) validateRequest(req *dto.LoginRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
