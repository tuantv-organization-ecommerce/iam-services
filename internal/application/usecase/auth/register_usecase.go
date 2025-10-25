package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tvttt/iam-services/internal/application/dto"
	"github.com/tvttt/iam-services/internal/domain/model"
	"github.com/tvttt/iam-services/internal/domain/repository"
	"github.com/tvttt/iam-services/internal/domain/service"
)

// RegisterUseCase handles user registration business logic
type RegisterUseCase struct {
	userRepo    repository.UserRepository
	passwordSvc service.PasswordService
}

// NewRegisterUseCase creates a new RegisterUseCase
func NewRegisterUseCase(
	userRepo repository.UserRepository,
	passwordSvc service.PasswordService,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
	}
}

// Execute executes the register use case
func (uc *RegisterUseCase) Execute(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 1. Validate input
	if err := uc.validateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Check if username or email already exists
	existsUsername, err := uc.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if existsUsername {
		return nil, fmt.Errorf("username already exists")
	}

	existsEmail, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if existsEmail {
		return nil, fmt.Errorf("email already exists")
	}

	// 3. Create user entity
	userID := uuid.New().String()
	user := model.NewUser(userID, req.Username, req.Email, req.FullName)

	// 4. Validate domain entity
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// 5. Hash password using domain service
	hashedPassword, err := uc.passwordSvc.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := user.SetPasswordHash(hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	// 6. Save user to repository
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// 7. Return response
	return &dto.RegisterResponse{
		UserID:   user.ID(),
		Username: user.Username(),
		Email:    user.Email(),
		Message:  "User registered successfully",
	}, nil
}

func (uc *RegisterUseCase) validateRequest(req *dto.RegisterRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}
