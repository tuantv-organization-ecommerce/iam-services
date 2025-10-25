package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/repository"
	casbinPkg "github.com/tvttt/iam-services/pkg/casbin"
)

// CasbinService handles Casbin-based authorization
type CasbinService interface {
	// Authorization checks
	CheckAPIAccess(ctx context.Context, userID, apiPath, method string) (bool, error)
	CheckCMSAccess(ctx context.Context, userID string, cmsTab domain.CMSTab, action string) (bool, error)
	Enforce(ctx context.Context, req *domain.AuthorizationRequest) (*domain.AuthorizationResponse, error)

	// Role management with Casbin
	AssignUserRole(ctx context.Context, userID, roleName string, domain domain.CasbinDomain) error
	RemoveUserRole(ctx context.Context, userID, roleName string, domain domain.CasbinDomain) error
	GetUserRolesInDomain(ctx context.Context, userID string, domain domain.CasbinDomain) ([]string, error)

	// Policy management
	AddPolicy(ctx context.Context, role string, domain domain.CasbinDomain, resource, action string) error
	RemovePolicy(ctx context.Context, role string, domain domain.CasbinDomain, resource, action string) error
	GetPoliciesForRole(ctx context.Context, role string, domain domain.CasbinDomain) ([][]string, error)

	// CMS Role management
	CreateCMSRole(ctx context.Context, name, description string, tabs []domain.CMSTab) (*domain.CMSRole, error)
	AssignCMSRole(ctx context.Context, userID, cmsRoleID string) error
	RemoveCMSRole(ctx context.Context, userID, cmsRoleID string) error
	GetUserCMSTabs(ctx context.Context, userID string) ([]domain.CMSTab, error)

	// API Resource management
	CreateAPIResource(ctx context.Context, path, method, service, description string) (*domain.APIResource, error)
	ListAPIResources(ctx context.Context, service string) ([]*domain.APIResource, error)
}

type casbinService struct {
	enforcer        *casbinPkg.Enforcer
	cmsRepo         repository.CMSRepository
	apiResourceRepo repository.APIResourceRepository
	userRepo        repository.UserRepository
	roleRepo        repository.RoleRepository
}

// NewCasbinService creates a new instance of CasbinService
func NewCasbinService(
	enforcer *casbinPkg.Enforcer,
	cmsRepo repository.CMSRepository,
	apiResourceRepo repository.APIResourceRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) CasbinService {
	return &casbinService{
		enforcer:        enforcer,
		cmsRepo:         cmsRepo,
		apiResourceRepo: apiResourceRepo,
		userRepo:        userRepo,
		roleRepo:        roleRepo,
	}
}

func (s *casbinService) CheckAPIAccess(ctx context.Context, userID, apiPath, method string) (bool, error) {
	// Check in API domain
	allowed, err := s.enforcer.Enforce(userID, string(domain.DomainAPI), apiPath, method)
	if err != nil {
		return false, fmt.Errorf("failed to enforce policy: %w", err)
	}
	return allowed, nil
}

func (s *casbinService) CheckCMSAccess(ctx context.Context, userID string, cmsTab domain.CMSTab, action string) (bool, error) {
	// Convert CMS tab to resource path
	resource := fmt.Sprintf("/cms/%s/*", string(cmsTab))

	// Check in CMS domain
	allowed, err := s.enforcer.Enforce(userID, string(domain.DomainCMS), resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to enforce policy: %w", err)
	}
	return allowed, nil
}

func (s *casbinService) Enforce(ctx context.Context, req *domain.AuthorizationRequest) (*domain.AuthorizationResponse, error) {
	allowed, err := s.enforcer.Enforce(
		req.UserID,
		string(req.Domain),
		req.Resource,
		req.Action,
	)

	if err != nil {
		return &domain.AuthorizationResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("Authorization check failed: %v", err),
		}, err
	}

	reason := "Permission granted"
	if !allowed {
		reason = "Permission denied"
	}

	return &domain.AuthorizationResponse{
		Allowed: allowed,
		Reason:  reason,
	}, nil
}

func (s *casbinService) AssignUserRole(ctx context.Context, userID, roleName string, dom domain.CasbinDomain) error {
	// Verify user exists
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify role exists in the domain
	if dom == domain.DomainUser || dom == domain.DomainAPI {
		if _, err := s.roleRepo.GetRoleByName(ctx, roleName); err != nil {
			return fmt.Errorf("role not found: %w", err)
		}
	}

	// Add role in Casbin
	return s.enforcer.AddRoleForUser(userID, roleName, string(dom))
}

func (s *casbinService) RemoveUserRole(ctx context.Context, userID, roleName string, dom domain.CasbinDomain) error {
	return s.enforcer.DeleteRoleForUser(userID, roleName, string(dom))
}

func (s *casbinService) GetUserRolesInDomain(ctx context.Context, userID string, dom domain.CasbinDomain) ([]string, error) {
	return s.enforcer.GetRolesForUser(userID, string(dom))
}

func (s *casbinService) AddPolicy(ctx context.Context, role string, dom domain.CasbinDomain, resource, action string) error {
	return s.enforcer.AddPolicy(role, string(dom), resource, action)
}

func (s *casbinService) RemovePolicy(ctx context.Context, role string, dom domain.CasbinDomain, resource, action string) error {
	return s.enforcer.RemovePolicy(role, string(dom), resource, action)
}

func (s *casbinService) GetPoliciesForRole(ctx context.Context, role string, dom domain.CasbinDomain) ([][]string, error) {
	return s.enforcer.GetPermissionsForUser(role, string(dom))
}

func (s *casbinService) CreateCMSRole(ctx context.Context, name, description string, tabs []domain.CMSTab) (*domain.CMSRole, error) {
	cmsRole := &domain.CMSRole{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Tabs:        tabs,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.cmsRepo.CreateCMSRole(ctx, cmsRole); err != nil {
		return nil, fmt.Errorf("failed to create CMS role: %w", err)
	}

	// Add Casbin policies for each tab
	for _, tab := range tabs {
		resource := fmt.Sprintf("/cms/%s/*", string(tab))
		// Default: allow GET (read) access
		if err := s.enforcer.AddPolicy(name, string(domain.DomainCMS), resource, "GET"); err != nil {
			return nil, fmt.Errorf("failed to add policy for tab %s: %w", tab, err)
		}
	}

	return cmsRole, nil
}

func (s *casbinService) AssignCMSRole(ctx context.Context, userID, cmsRoleID string) error {
	// Get CMS role
	cmsRole, err := s.cmsRepo.GetCMSRoleByID(ctx, cmsRoleID)
	if err != nil {
		return fmt.Errorf("CMS role not found: %w", err)
	}

	// Assign in database
	if err := s.cmsRepo.AssignCMSRoleToUser(ctx, userID, cmsRoleID); err != nil {
		return fmt.Errorf("failed to assign CMS role: %w", err)
	}

	// Add role in Casbin
	return s.enforcer.AddRoleForUser(userID, cmsRole.Name, string(domain.DomainCMS))
}

func (s *casbinService) RemoveCMSRole(ctx context.Context, userID, cmsRoleID string) error {
	// Get CMS role
	cmsRole, err := s.cmsRepo.GetCMSRoleByID(ctx, cmsRoleID)
	if err != nil {
		return fmt.Errorf("CMS role not found: %w", err)
	}

	// Remove from database
	if err := s.cmsRepo.RemoveCMSRoleFromUser(ctx, userID, cmsRoleID); err != nil {
		return fmt.Errorf("failed to remove CMS role: %w", err)
	}

	// Remove role in Casbin
	return s.enforcer.DeleteRoleForUser(userID, cmsRole.Name, string(domain.DomainCMS))
}

func (s *casbinService) GetUserCMSTabs(ctx context.Context, userID string) ([]domain.CMSTab, error) {
	return s.cmsRepo.GetUserCMSTabs(ctx, userID)
}

func (s *casbinService) CreateAPIResource(ctx context.Context, path, method, service, description string) (*domain.APIResource, error) {
	resource := &domain.APIResource{
		ID:          uuid.New().String(),
		Path:        path,
		Method:      method,
		Service:     service,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.apiResourceRepo.CreateAPIResource(ctx, resource); err != nil {
		return nil, fmt.Errorf("failed to create API resource: %w", err)
	}

	return resource, nil
}

func (s *casbinService) ListAPIResources(ctx context.Context, service string) ([]*domain.APIResource, error) {
	if service != "" {
		return s.apiResourceRepo.ListAPIResourcesByService(ctx, service)
	}

	resources, _, err := s.apiResourceRepo.ListAPIResources(ctx, 1, 1000)
	return resources, err
}
