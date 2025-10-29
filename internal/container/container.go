package container

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/tvttt/iam-services/internal/config"
	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/handler"
	"github.com/tvttt/iam-services/internal/repository"
	"github.com/tvttt/iam-services/internal/service"
	casbinPkg "github.com/tvttt/iam-services/pkg/casbin"
	"github.com/tvttt/iam-services/pkg/jwt"
	"github.com/tvttt/iam-services/pkg/password"
)

// Container holds all application dependencies following DIP
type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB             *sql.DB
	Logger         *zap.Logger
	CasbinEnforcer *casbinPkg.Enforcer

	// Managers (external packages)
	JWTManager      *jwt.JWTManager
	PasswordManager *password.PasswordManager

	// DAOs
	DAOs *DAORegistry

	// Services
	Services *ServiceRegistry

	// Handlers
	GRPCHandler *handler.GRPCHandler
}

// DAORegistry holds all DAOs
type DAORegistry struct {
	User           dao.UserDAO
	Role           dao.RoleDAO
	Permission     dao.PermissionDAO
	UserRole       dao.UserRoleDAO
	RolePermission dao.RolePermissionDAO
	APIResource    dao.APIResourceDAO
	CMSRole        dao.CMSRoleDAO
	UserCMSRole    dao.UserCMSRoleDAO
}

// ServiceRegistry holds all services
type ServiceRegistry struct {
	// Application services (legacy - for handlers)
	Auth          service.AuthService
	Authorization service.AuthorizationService
	Role          service.RoleService
	Permission    service.PermissionService
	Casbin        service.CasbinService
}

// NewContainer creates and wires all dependencies
func NewContainer(
	cfg *config.Config,
	db *sql.DB,
	logger *zap.Logger,
	casbinEnforcer *casbinPkg.Enforcer,
) (*Container, error) {
	c := &Container{
		Config:         cfg,
		DB:             db,
		Logger:         logger,
		CasbinEnforcer: casbinEnforcer,
	}

	// Initialize managers
	c.initializeManagers()

	// Initialize DAOs
	c.initializeDAOs()

	// Initialize services
	if err := c.initializeServices(); err != nil {
		return nil, err
	}

	// Initialize handlers
	c.initializeHandlers()

	logger.Info("Dependency container initialized successfully")

	return c, nil
}

// initializeManagers creates external package managers
func (c *Container) initializeManagers() {
	c.JWTManager = jwt.NewJWTManager(
		c.Config.JWT.Secret,
		c.Config.JWT.AccessTokenDuration,
		c.Config.JWT.RefreshTokenDuration,
	)
	c.PasswordManager = password.NewPasswordManager()
}

// initializeDAOs creates all Data Access Objects
func (c *Container) initializeDAOs() {
	c.DAOs = &DAORegistry{
		User:           dao.NewUserDAO(c.DB),
		Role:           dao.NewRoleDAO(c.DB),
		Permission:     dao.NewPermissionDAO(c.DB),
		UserRole:       dao.NewUserRoleDAO(c.DB),
		RolePermission: dao.NewRolePermissionDAO(c.DB),
		APIResource:    dao.NewAPIResourceDAO(c.DB),
		CMSRole:        dao.NewCMSRoleDAO(c.DB),
		UserCMSRole:    dao.NewUserCMSRoleDAO(c.DB),
	}
}

// initializeServices creates all application and domain services
func (c *Container) initializeServices() error {
	c.Services = &ServiceRegistry{}

	// Application services (for current handlers)
	c.Services.Auth = service.NewAuthService(
		repository.NewUserRepository(c.DAOs.User),
		repository.NewAuthorizationRepository(c.DAOs.UserRole, c.DAOs.RolePermission, c.DAOs.Permission),
		c.JWTManager,
		c.PasswordManager,
	)

	c.Services.Authorization = service.NewAuthorizationService(
		repository.NewAuthorizationRepository(c.DAOs.UserRole, c.DAOs.RolePermission, c.DAOs.Permission),
		repository.NewUserRepository(c.DAOs.User),
		repository.NewRoleRepository(c.DAOs.Role, c.DAOs.RolePermission),
	)

	c.Services.Role = service.NewRoleService(
		repository.NewRoleRepository(c.DAOs.Role, c.DAOs.RolePermission),
		repository.NewAuthorizationRepository(c.DAOs.UserRole, c.DAOs.RolePermission, c.DAOs.Permission),
	)

	c.Services.Permission = service.NewPermissionService(
		repository.NewPermissionRepository(c.DAOs.Permission),
	)

	c.Services.Casbin = service.NewCasbinService(
		c.CasbinEnforcer,
		repository.NewCMSRepository(c.DAOs.CMSRole, c.DAOs.UserCMSRole),
		repository.NewAPIResourceRepository(c.DAOs.APIResource),
		repository.NewUserRepository(c.DAOs.User),
		repository.NewRoleRepository(c.DAOs.Role, c.DAOs.RolePermission),
	)

	return nil
}

// initializeHandlers creates gRPC handlers
func (c *Container) initializeHandlers() {
	c.GRPCHandler = handler.NewGRPCHandler(
		c.Services.Auth,
		c.Services.Authorization,
		c.Services.Role,
		c.Services.Permission,
		c.Services.Casbin,
		c.Logger,
	)
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.Logger != nil {
		if err := c.Logger.Sync(); err != nil {
			// Ignore sync errors on stdout/stderr (common on some platforms)
			c.Logger.Warn("Logger sync returned error (may be harmless)", zap.Error(err))
		}
	}
	return nil
}
