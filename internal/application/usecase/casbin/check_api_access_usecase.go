package casbin

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/application/dto"
	"github.com/tvttt/iam-services/internal/domain/service"
)

// CheckAPIAccessUseCase handles API access check business logic
type CheckAPIAccessUseCase struct {
	authzSvc service.AuthorizationService
}

// NewCheckAPIAccessUseCase creates a new CheckAPIAccessUseCase
func NewCheckAPIAccessUseCase(authzSvc service.AuthorizationService) *CheckAPIAccessUseCase {
	return &CheckAPIAccessUseCase{
		authzSvc: authzSvc,
	}
}

// Execute executes the check API access use case
func (uc *CheckAPIAccessUseCase) Execute(ctx context.Context, req *dto.CheckAPIAccessRequest) (*dto.CheckAPIAccessResponse, error) {
	// 1. Validate input
	if err := uc.validateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Check access using Casbin authorization service
	allowed, err := uc.authzSvc.Enforce(
		ctx,
		req.UserID,
		service.DomainAPI,
		req.APIPath,
		req.Method,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check API access: %w", err)
	}

	// 3. Return response
	message := "Access denied"
	if allowed {
		message = "Access granted"
	}

	return &dto.CheckAPIAccessResponse{
		Allowed: allowed,
		Message: message,
	}, nil
}

func (uc *CheckAPIAccessUseCase) validateRequest(req *dto.CheckAPIAccessRequest) error {
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.APIPath == "" {
		return fmt.Errorf("api_path is required")
	}
	if req.Method == "" {
		return fmt.Errorf("method is required")
	}
	return nil
}
