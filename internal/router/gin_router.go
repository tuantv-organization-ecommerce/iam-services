package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tvttt/gokits/swagger"
	"github.com/tvttt/iam-services/internal/config"
	"github.com/tvttt/iam-services/internal/handler"
	"github.com/tvttt/iam-services/internal/middleware"
)

// SetupGinRouter sets up the Gin router with all routes
func SetupGinRouter(
	cfg *config.Config,
	ginHandler *handler.GinHandler,
	logger *zap.Logger,
) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	r := gin.New()

	// Apply global middleware
	r.Use(middleware.GinRecovery(logger))
	r.Use(middleware.GinLogger(logger))
	r.Use(middleware.GinCORS())

	// Health check endpoint
	r.GET("/health", ginHandler.Health)

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", ginHandler.Register)
			auth.POST("/login", ginHandler.Login)
			auth.POST("/refresh", ginHandler.RefreshToken)
			auth.POST("/logout", ginHandler.Logout)
			auth.POST("/verify", ginHandler.VerifyToken)
		}

		// Role management routes (requires authentication)
		roles := v1.Group("/roles")
		{
			roles.POST("", ginHandler.CreateRole)
			roles.GET("", ginHandler.ListRoles)
			roles.GET("/:id", ginHandler.GetRole)
			roles.PUT("/:id", ginHandler.UpdateRole)
			roles.DELETE("/:id", ginHandler.DeleteRole)
			roles.POST("/assign", ginHandler.AssignRole)
			roles.POST("/remove", ginHandler.RemoveRole)
		}

		// User roles routes
		users := v1.Group("/users")
		{
			users.GET("/:user_id/roles", ginHandler.GetUserRoles)
		}

		// Permission management routes
		permissions := v1.Group("/permissions")
		{
			permissions.POST("", ginHandler.CreatePermission)
			permissions.GET("", ginHandler.ListPermissions)
			permissions.DELETE("/:id", ginHandler.DeletePermission)
			permissions.POST("/check", ginHandler.CheckPermission)
		}

		// Access control routes
		access := v1.Group("/access")
		{
			access.POST("/api", ginHandler.CheckAPIAccess)
			access.POST("/cms", ginHandler.CheckCMSAccess)
		}

		// Policy routes
		policies := v1.Group("/policies")
		{
			policies.POST("/enforce", ginHandler.EnforcePolicy)
		}

		// CMS routes
		cms := v1.Group("/cms")
		{
			cmsRoles := cms.Group("/roles")
			{
				cmsRoles.POST("", ginHandler.CreateCMSRole)
				cmsRoles.GET("", ginHandler.ListCMSRoles)
				cmsRoles.POST("/assign", ginHandler.AssignCMSRole)
				cmsRoles.POST("/remove", ginHandler.RemoveCMSRole)
			}

			cmsUsers := cms.Group("/users")
			{
				cmsUsers.GET("/:user_id/tabs", ginHandler.GetUserCMSTabs)
			}
		}

		// API Resource management routes
		apiResources := v1.Group("/api/resources")
		{
			apiResources.POST("", ginHandler.CreateAPIResource)
			apiResources.GET("", ginHandler.ListAPIResources)
		}
	}

	// Setup Swagger UI if enabled
	if cfg.Swagger.Enabled {
		setupSwagger(r, cfg, logger)
	}

	return r
}

// setupSwagger configures Swagger UI for Gin
func setupSwagger(r *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	var basicAuth *swagger.BasicAuthConfig
	if cfg.Swagger.AuthUsername != "" && cfg.Swagger.AuthPassword != "" {
		basicAuth = &swagger.BasicAuthConfig{
			Username: cfg.Swagger.AuthUsername,
			Password: cfg.Swagger.AuthPassword,
			Realm:    cfg.Swagger.AuthRealm,
		}
	}

	swaggerCfg := &swagger.Config{
		BasePath:  cfg.Swagger.BasePath,
		SpecPath:  cfg.Swagger.SpecPath,
		Title:     cfg.Swagger.Title,
		Enabled:   true,
		BasicAuth: basicAuth,
	}

	// Swagger UI route with basic auth
	r.GET(cfg.Swagger.BasePath+"*filepath", gin.WrapF(swagger.Handler(swaggerCfg, logger)))

	// Swagger spec route with basic auth
	r.GET(cfg.Swagger.SpecPath, gin.WrapF(swagger.ServeSpec("./pkg/proto/iam.swagger.json", logger, basicAuth)))

	logger.Info("Swagger UI configured",
		zap.String("url", cfg.Swagger.BasePath),
		zap.String("spec", cfg.Swagger.SpecPath),
	)
}

