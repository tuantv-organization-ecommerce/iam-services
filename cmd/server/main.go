package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	// Shared utilities
	"github.com/tvttt/gokits/logger"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/database"

	// Infrastructure layer
	"github.com/tvttt/iam-services/internal/infrastructure/authorization"
	infraConfig "github.com/tvttt/iam-services/internal/infrastructure/config"
	"github.com/tvttt/iam-services/internal/infrastructure/persistence"
	"github.com/tvttt/iam-services/internal/infrastructure/security"

	// Old layer for compatibility (handlers)
	"github.com/tvttt/iam-services/internal/handler"
	"github.com/tvttt/iam-services/internal/repository"
	"github.com/tvttt/iam-services/internal/service"

	// External packages
	casbinPkg "github.com/tvttt/iam-services/pkg/casbin"
	"github.com/tvttt/iam-services/pkg/jwt"
	"github.com/tvttt/iam-services/pkg/password"
	pb "github.com/tvttt/iam-services/pkg/proto"
)

func main() {
	// Initialize logger using shared gokits package
	log, err := logger.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync(log)

	log.Info("Starting IAM Service (Clean Architecture)...")

	// Load configuration
	cfg, err := infraConfig.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Validate configuration
	if err := infraConfig.ValidateConfig(cfg); err != nil {
		log.Fatal("Invalid configuration", zap.Error(err))
	}

	// Connect to database
	db, err := database.Connect(cfg.Database.GetDSN(), log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close(db, log)

	// Initialize Casbin enforcer
	casbinEnforcer, err := casbinPkg.NewEnforcer(db, "configs/rbac_model.conf", log)
	if err != nil {
		log.Fatal("Failed to initialize Casbin enforcer", zap.Error(err))
	}
	log.Info("Casbin enforcer initialized successfully")

	// ============================================
	// DEPENDENCY INJECTION - CLEAN ARCHITECTURE
	// ============================================

	// 1. Initialize DAOs (Data Access Layer)
	userDAO := dao.NewUserDAO(db)
	roleDAO := dao.NewRoleDAO(db)
	permissionDAO := dao.NewPermissionDAO(db)
	userRoleDAO := dao.NewUserRoleDAO(db)
	rolePermissionDAO := dao.NewRolePermissionDAO(db)
	apiResourceDAO := dao.NewAPIResourceDAO(db)
	cmsRoleDAO := dao.NewCMSRoleDAO(db)
	userCMSRoleDAO := dao.NewUserCMSRoleDAO(db)

	// 2. Initialize Infrastructure Layer (Implements Domain Ports)

	// 2.1 Persistence - Repository implementations
	userRepoNew := persistence.NewUserRepository(userDAO)
	roleRepoNew := persistence.NewRoleRepository(roleDAO, rolePermissionDAO)
	permissionRepoNew := persistence.NewPermissionRepository(permissionDAO)
	authzRepoNew := persistence.NewAuthorizationRepository(userRoleDAO, rolePermissionDAO)
	apiResourceRepoNew := persistence.NewAPIResourceRepository(apiResourceDAO)
	cmsRepoNew := persistence.NewCMSRepository(cmsRoleDAO, userCMSRoleDAO)

	// 2.2 Security - Token and Password services
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
	passwordMgr := password.NewPasswordManager()

	tokenService := security.NewJWTService(jwtManager)
	passwordService := security.NewPasswordService(passwordMgr)

	// 2.3 Authorization - Casbin service
	casbinAuthzService := authorization.NewCasbinService(casbinEnforcer)

	// 3. Initialize Old Services (for compatibility with existing gRPC handlers)
	// TODO: Replace these with new use cases in handlers
	authService := service.NewAuthService(
		repository.NewUserRepository(userDAO),
		repository.NewAuthorizationRepository(userRoleDAO, rolePermissionDAO, permissionDAO),
		jwtManager,
		passwordMgr,
	)
	authzService := service.NewAuthorizationService(
		repository.NewAuthorizationRepository(userRoleDAO, rolePermissionDAO, permissionDAO),
		repository.NewUserRepository(userDAO),
		repository.NewRoleRepository(roleDAO, rolePermissionDAO),
	)
	roleService := service.NewRoleService(
		repository.NewRoleRepository(roleDAO, rolePermissionDAO),
		repository.NewAuthorizationRepository(userRoleDAO, rolePermissionDAO, permissionDAO),
	)
	permService := service.NewPermissionService(
		repository.NewPermissionRepository(permissionDAO),
	)
	casbinService := service.NewCasbinService(
		casbinEnforcer,
		repository.NewCMSRepository(cmsRoleDAO, userCMSRoleDAO),
		repository.NewAPIResourceRepository(apiResourceDAO),
		repository.NewUserRepository(userDAO),
		repository.NewRoleRepository(roleDAO, rolePermissionDAO),
	)

	// 4. Initialize Adapter Layer (gRPC handlers)
	grpcHandler := handler.NewGRPCHandler(
		authService,
		authzService,
		roleService,
		permService,
		casbinService,
		log,
	)

	// Log initialization of new infrastructure (for demonstration)
	log.Info("New infrastructure initialized",
		zap.String("userRepo", fmt.Sprintf("%T", userRepoNew)),
		zap.String("roleRepo", fmt.Sprintf("%T", roleRepoNew)),
		zap.String("permissionRepo", fmt.Sprintf("%T", permissionRepoNew)),
		zap.String("authzRepo", fmt.Sprintf("%T", authzRepoNew)),
		zap.String("apiResourceRepo", fmt.Sprintf("%T", apiResourceRepoNew)),
		zap.String("cmsRepo", fmt.Sprintf("%T", cmsRepoNew)),
		zap.String("tokenService", fmt.Sprintf("%T", tokenService)),
		zap.String("passwordService", fmt.Sprintf("%T", passwordService)),
		zap.String("casbinAuthzService", fmt.Sprintf("%T", casbinAuthzService)),
	)

	// Suppress unused variable warnings (these will be used when handlers are updated)
	_ = userRepoNew
	_ = roleRepoNew
	_ = permissionRepoNew
	_ = authzRepoNew
	_ = apiResourceRepoNew
	_ = cmsRepoNew
	_ = tokenService
	_ = passwordService
	_ = casbinAuthzService

	// ============================================
	// gRPC SERVER SETUP
	// ============================================

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterIAMServiceServer(grpcServer, grpcHandler)

	// Register reflection service (useful for tools like grpcurl)
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcAddress := cfg.Server.GetServerAddress()
	grpcListener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("Failed to listen on gRPC address", zap.String("address", grpcAddress), zap.Error(err))
	}

	// Start gRPC server in goroutine
	go func() {
		log.Info("gRPC server is running", zap.String("address", grpcAddress))
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// ============================================
	// HTTP GATEWAY SERVER SETUP
	// ============================================

	// Create gRPC client connection for gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup gRPC Gateway mux
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register gateway handlers
	err = pb.RegisterIAMGatewayServiceHandlerFromEndpoint(ctx, gwMux, grpcAddress, opts)
	if err != nil {
		log.Fatal("Failed to register gateway handler", zap.Error(err))
	}

	// Create HTTP server
	httpAddress := cfg.Server.GetHTTPServerAddress()
	httpServer := &http.Server{
		Addr:    httpAddress,
		Handler: corsMiddleware(gwMux),
	}

	// Start HTTP gateway in goroutine
	go func() {
		log.Info("HTTP gateway server is running", zap.String("address", httpAddress))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to serve HTTP gateway", zap.Error(err))
		}
	}()

	// ============================================
	// GRACEFUL SHUTDOWN
	// ============================================

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down gracefully...")

	// Shutdown HTTP gateway
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Failed to shutdown HTTP server gracefully", zap.Error(err))
	} else {
		log.Info("HTTP gateway stopped")
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()
	log.Info("gRPC server stopped")

	log.Info("Service shutdown complete")
}

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
