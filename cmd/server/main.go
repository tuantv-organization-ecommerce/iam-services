package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tvttt/iam-services/internal/config"
	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/database"
	"github.com/tvttt/iam-services/internal/handler"
	"github.com/tvttt/iam-services/internal/repository"
	"github.com/tvttt/iam-services/internal/service"
	"github.com/tvttt/iam-services/pkg/jwt"
	"github.com/tvttt/iam-services/pkg/password"
	pb "github.com/tvttt/iam-services/pkg/proto"
)

func main() {
	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting IAM Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Connect to database
	db, err := database.Connect(cfg.Database.GetDSN(), logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close(db, logger)

	// Initialize DAOs
	userDAO := dao.NewUserDAO(db)
	roleDAO := dao.NewRoleDAO(db)
	permissionDAO := dao.NewPermissionDAO(db)
	userRoleDAO := dao.NewUserRoleDAO(db)
	rolePermissionDAO := dao.NewRolePermissionDAO(db)

	// Initialize Repositories
	userRepo := repository.NewUserRepository(userDAO)
	roleRepo := repository.NewRoleRepository(roleDAO, rolePermissionDAO)
	permissionRepo := repository.NewPermissionRepository(permissionDAO)
	authzRepo := repository.NewAuthorizationRepository(userRoleDAO, rolePermissionDAO, permissionDAO)

	// Initialize JWT manager and password manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
	passwordMgr := password.NewPasswordManager()

	// Initialize Services
	authService := service.NewAuthService(userRepo, authzRepo, jwtManager, passwordMgr)
	authzService := service.NewAuthorizationService(authzRepo, userRepo, roleRepo)
	roleService := service.NewRoleService(roleRepo, authzRepo)
	permService := service.NewPermissionService(permissionRepo)

	// Initialize gRPC handler
	grpcHandler := handler.NewGRPCHandler(authService, authzService, roleService, permService, logger)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterIAMServiceServer(grpcServer, grpcHandler)

	// Register reflection service (useful for tools like grpcurl)
	reflection.Register(grpcServer)

	// Start listening
	address := cfg.Server.GetServerAddress()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("Failed to listen", zap.String("address", address), zap.Error(err))
	}

	logger.Info("IAM Service is running", zap.String("address", address))

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	return config.Build()
}
