package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tvttt/gokits/logger"

	"github.com/tvttt/iam-services/internal/config"
	"github.com/tvttt/iam-services/internal/container"
	"github.com/tvttt/iam-services/internal/database"
	infraConfig "github.com/tvttt/iam-services/internal/infrastructure/config"
	"github.com/tvttt/iam-services/internal/middleware"
	casbinPkg "github.com/tvttt/iam-services/pkg/casbin"
	pb "github.com/tvttt/iam-services/pkg/proto"
)

// App represents the IAM service application
type App struct {
	config     *config.Config
	logger     *zap.Logger
	container  *container.Container
	db         *sql.DB
	grpcServer *grpc.Server
}

// New creates a new App instance
func New() (*App, error) {
	app := &App{}

	// Initialize logger
	log, err := logger.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	app.logger = log

	app.logger.Info("Starting IAM Service...")

	// Load configuration
	cfg, err := infraConfig.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	app.config = cfg

	// Validate configuration
	if err := infraConfig.ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	app.logger.Info("Configuration loaded successfully")

	return app, nil
}

// Initialize sets up all dependencies and services
func (a *App) Initialize() error {
	// Wrap initialization with panic recovery
	return middleware.RecoverFunc(a.logger, "Initialize", func() error {
		// Connect to database
		db, err := database.Connect(a.config.Database.GetDSN(), a.logger)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		a.db = db

		// Initialize Casbin enforcer
		casbinEnforcer, err := casbinPkg.NewEnforcer(db, "configs/rbac_model.conf", a.logger)
		if err != nil {
			return fmt.Errorf("failed to initialize Casbin enforcer: %w", err)
		}

		a.logger.Info("Casbin enforcer initialized successfully")

		// Create dependency container
		c, err := container.NewContainer(a.config, db, a.logger, casbinEnforcer)
		if err != nil {
			return fmt.Errorf("failed to create dependency container: %w", err)
		}
		a.container = c

		// Setup gRPC server with panic recovery interceptors
		a.grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				middleware.RecoveryUnaryInterceptor(a.logger),
			),
			grpc.ChainStreamInterceptor(
				middleware.RecoveryStreamInterceptor(a.logger),
			),
		)
		pb.RegisterIAMServiceServer(a.grpcServer, c.GRPCHandler)
		reflection.Register(a.grpcServer)

		a.logger.Info("Application initialized successfully")

		return nil
	})
}

// Run starts the application
func (a *App) Run() error {
	// Start gRPC server
	address := a.config.Server.GetServerAddress()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	// Start gRPC server in goroutine with panic recovery
	middleware.RecoverGoroutine(a.logger, "grpc-server", func() {
		a.logger.Info("gRPC server is running", zap.String("address", address))
		if err := a.grpcServer.Serve(listener); err != nil {
			a.logger.Error("gRPC server stopped with error", zap.Error(err))
		}
	})

	// Wait for shutdown signal
	a.waitForShutdown()

	return nil
}

// waitForShutdown handles graceful shutdown
func (a *App) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	a.logger.Info("Shutting down gracefully...")
	a.Shutdown(context.Background())
}

// Shutdown performs graceful shutdown
func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErrors []error

	// Stop gRPC server
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
		a.logger.Info("gRPC server stopped")
	}

	// Close container resources
	if a.container != nil {
		if err := a.container.Close(); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("container close error: %w", err))
		}
	}

	// Close database connection
	if a.db != nil {
		database.Close(a.db, a.logger)
	}

	// Sync logger
	if a.logger != nil {
		_ = a.logger.Sync()
	}

	a.logger.Info("Service shutdown complete")

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	return nil
}
