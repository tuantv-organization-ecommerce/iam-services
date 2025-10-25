package authorization

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/service"
	casbinPkg "github.com/tvttt/iam-services/pkg/casbin"
)

// casbinServiceImpl implements domain.service.AuthorizationService using Casbin
type casbinServiceImpl struct {
	enforcer *casbinPkg.Enforcer
}

// NewCasbinService creates a new Casbin-based authorization service
func NewCasbinService(enforcer *casbinPkg.Enforcer) service.AuthorizationService {
	return &casbinServiceImpl{
		enforcer: enforcer,
	}
}

func (s *casbinServiceImpl) Enforce(ctx context.Context, subject string, domain service.CasbinDomain, object, action string) (bool, error) {
	return s.enforcer.Enforce(subject, string(domain), object, action)
}

func (s *casbinServiceImpl) AddPolicy(ctx context.Context, subject string, domain service.CasbinDomain, object, action string) error {
	return s.enforcer.AddPolicy(subject, string(domain), object, action)
}

func (s *casbinServiceImpl) RemovePolicy(ctx context.Context, subject string, domain service.CasbinDomain, object, action string) error {
	return s.enforcer.RemovePolicy(subject, string(domain), object, action)
}

func (s *casbinServiceImpl) AddRoleForUser(ctx context.Context, userID, role string, domain service.CasbinDomain) error {
	return s.enforcer.AddRoleForUser(userID, role, string(domain))
}

func (s *casbinServiceImpl) RemoveRoleForUser(ctx context.Context, userID, role string, domain service.CasbinDomain) error {
	return s.enforcer.DeleteRoleForUser(userID, role, string(domain))
}

func (s *casbinServiceImpl) GetRolesForUser(ctx context.Context, userID string, domain service.CasbinDomain) ([]string, error) {
	return s.enforcer.GetRolesForUser(userID, string(domain))
}

func (s *casbinServiceImpl) GetPermissionsForUser(ctx context.Context, userID string, domain service.CasbinDomain) ([][]string, error) {
	return s.enforcer.GetPermissionsForUser(userID, string(domain))
}
