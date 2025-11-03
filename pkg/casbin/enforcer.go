package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Enforcer wraps Casbin enforcer with helper methods
type Enforcer struct {
	enforcer *casbin.Enforcer
	logger   *zap.Logger
}

// NewEnforcer creates a new Casbin enforcer instance
// Uses DSN to create a separate GORM connection to avoid prepared statement conflicts
func NewEnforcer(dsn string, modelPath string, logger *zap.Logger) (*Enforcer, error) {
	// Create a separate GORM DB instance for Casbin adapter to avoid connection conflicts
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false, // Disable prepared statements to avoid parameter mismatch errors
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM DB: %w", err)
	}

	// Create Casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin adapter: %w", err)
	}

	// Load model from file
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load Casbin model: %w", err)
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %w", err)
	}

	// Load policies from database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	logger.Info("Casbin enforcer initialized successfully")

	return &Enforcer{
		enforcer: enforcer,
		logger:   logger,
	}, nil
}

// Enforce checks if a subject can perform an action on an object in a domain
func (e *Enforcer) Enforce(subject, domain, object, action string) (bool, error) {
	allowed, err := e.enforcer.Enforce(subject, domain, object, action)
	if err != nil {
		e.logger.Error("Failed to enforce policy",
			zap.String("subject", subject),
			zap.String("domain", domain),
			zap.String("object", object),
			zap.String("action", action),
			zap.Error(err))
		return false, err
	}

	e.logger.Debug("Policy enforcement result",
		zap.String("subject", subject),
		zap.String("domain", domain),
		zap.String("object", object),
		zap.String("action", action),
		zap.Bool("allowed", allowed))

	return allowed, nil
}

// AddPolicy adds a policy rule
func (e *Enforcer) AddPolicy(subject, domain, object, action string) error {
	added, err := e.enforcer.AddPolicy(subject, domain, object, action)
	if err != nil {
		return fmt.Errorf("failed to add policy: %w", err)
	}
	if !added {
		e.logger.Warn("Policy already exists",
			zap.String("subject", subject),
			zap.String("domain", domain),
			zap.String("object", object),
			zap.String("action", action))
	}
	return nil
}

// RemovePolicy removes a policy rule
func (e *Enforcer) RemovePolicy(subject, domain, object, action string) error {
	removed, err := e.enforcer.RemovePolicy(subject, domain, object, action)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}
	if !removed {
		e.logger.Warn("Policy not found",
			zap.String("subject", subject),
			zap.String("domain", domain),
			zap.String("object", object),
			zap.String("action", action))
	}
	return nil
}

// AddRoleForUser adds a role for a user in a domain
func (e *Enforcer) AddRoleForUser(user, role, domain string) error {
	added, err := e.enforcer.AddRoleForUser(user, role, domain)
	if err != nil {
		return fmt.Errorf("failed to add role for user: %w", err)
	}
	if !added {
		e.logger.Warn("Role assignment already exists",
			zap.String("user", user),
			zap.String("role", role),
			zap.String("domain", domain))
	}
	return nil
}

// DeleteRoleForUser deletes a role for a user in a domain
func (e *Enforcer) DeleteRoleForUser(user, role, domain string) error {
	deleted, err := e.enforcer.DeleteRoleForUser(user, role, domain)
	if err != nil {
		return fmt.Errorf("failed to delete role for user: %w", err)
	}
	if !deleted {
		e.logger.Warn("Role assignment not found",
			zap.String("user", user),
			zap.String("role", role),
			zap.String("domain", domain))
	}
	return nil
}

// GetRolesForUser gets all roles for a user in a domain
func (e *Enforcer) GetRolesForUser(user, domain string) ([]string, error) {
	roles, err := e.enforcer.GetRolesForUser(user, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for user: %w", err)
	}
	return roles, nil
}

// GetUsersForRole gets all users that have a role in a domain
func (e *Enforcer) GetUsersForRole(role, domain string) ([]string, error) {
	users, err := e.enforcer.GetUsersForRole(role, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for role: %w", err)
	}
	return users, nil
}

// GetPermissionsForUser gets all permissions for a user in a domain
func (e *Enforcer) GetPermissionsForUser(user, domain string) ([][]string, error) {
	permissions, err := e.enforcer.GetPermissionsForUser(user, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for user: %w", err)
	}
	return permissions, nil
}

// LoadPolicy reloads the policy from database
func (e *Enforcer) LoadPolicy() error {
	return e.enforcer.LoadPolicy()
}

// SavePolicy saves all policy rules to database
func (e *Enforcer) SavePolicy() error {
	return e.enforcer.SavePolicy()
}

// GetAllSubjects gets all subjects (users/roles)
func (e *Enforcer) GetAllSubjects() []string {
	return e.enforcer.GetAllSubjects()
}

// GetAllObjects gets all objects (resources)
func (e *Enforcer) GetAllObjects() []string {
	return e.enforcer.GetAllObjects()
}

// GetAllActions gets all actions
func (e *Enforcer) GetAllActions() []string {
	return e.enforcer.GetAllActions()
}

// GetAllRoles gets all roles
func (e *Enforcer) GetAllRoles() []string {
	return e.enforcer.GetAllRoles()
}
