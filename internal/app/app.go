package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Uncomment after generating proto files
	// "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	// "google.golang.org/grpc/credentials/insecure"
	// "github.com/tvttt/gokits/swagger"

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
	httpServer *http.Server
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
	grpcAddress := a.config.Server.GetServerAddress()
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcAddress, err)
	}

	// Start gRPC server in goroutine with panic recovery
	middleware.RecoverGoroutine(a.logger, "grpc-server", func() {
		a.logger.Info("gRPC server is running", zap.String("address", grpcAddress))
		if err := a.grpcServer.Serve(listener); err != nil {
			a.logger.Error("gRPC server stopped with error", zap.Error(err))
		}
	})

	// Setup HTTP Gateway + Swagger
	// TODO: Uncomment after generating proto files
	// Run: powershell -ExecutionPolicy Bypass -File .\scripts\generate-proto.ps1
	/*
		if err := a.setupHTTPGateway(); err != nil {
			return fmt.Errorf("failed to setup HTTP gateway: %w", err)
		}
	*/

	// Wait for shutdown signal
	a.waitForShutdown()

	return nil
}

// setupHTTPGateway sets up the HTTP gateway with Swagger UI
// Uncomment after generating proto files
/*
func (a *App) setupHTTPGateway() error {
	// Create gRPC Gateway mux
	gwMux := runtime.NewServeMux()

	// Register gRPC-Gateway handlers
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcEndpoint := a.config.Server.GetServerAddress()

	err := pb.RegisterIAMServiceHandlerFromEndpoint(context.Background(), gwMux, grpcEndpoint, opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	// Create HTTP mux
	mux := http.NewServeMux()

	// Register gRPC Gateway
	mux.Handle("/", corsMiddleware(gwMux))

	// Register Swagger UI
	if a.config.Swagger.Enabled {
		swaggerCfg := &swagger.Config{
			BasePath: a.config.Swagger.BasePath,
			SpecPath: a.config.Swagger.SpecPath,
			Title:    a.config.Swagger.Title,
			Enabled:  true,
		}

		// Swagger UI handler
		mux.HandleFunc(a.config.Swagger.BasePath, swagger.Handler(swaggerCfg, a.logger))

		// Swagger spec handler
		mux.HandleFunc(a.config.Swagger.SpecPath, swagger.ServeSpec("./pkg/proto/iam_gateway.swagger.json", a.logger))

		a.logger.Info("Swagger UI enabled",
			zap.String("url", fmt.Sprintf("http://%s%s", a.config.Server.GetHTTPServerAddress(), a.config.Swagger.BasePath)),
		)
	}

	// Create HTTP server
	httpAddress := a.config.Server.GetHTTPServerAddress()
	a.httpServer = &http.Server{
		Addr:         httpAddress,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in goroutine with panic recovery
	middleware.RecoverGoroutine(a.logger, "http-server", func() {
		a.logger.Info("HTTP Gateway is running",
			zap.String("address", httpAddress),
			zap.String("swagger", fmt.Sprintf("http://%s%s", httpAddress, a.config.Swagger.BasePath)),
		)

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server stopped with error", zap.Error(err))
		}
	})

	return nil
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
*/

// waitForShutdown handles graceful shutdown
func (a *App) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	a.logger.Info("Shutting down gracefully...")
	if err := a.Shutdown(context.Background()); err != nil {
		a.logger.Error("Error during shutdown", zap.Error(err))
	}
}

// Shutdown performs graceful shutdown
func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErrors []error

	// Stop HTTP server
	if a.httpServer != nil {
		a.logger.Info("Stopping HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("HTTP server shutdown error: %w", err))
		} else {
			a.logger.Info("HTTP server stopped")
		}
	}

	// Stop gRPC server
	if a.grpcServer != nil {
		a.logger.Info("Stopping gRPC server...")
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
		if syncErr := a.logger.Sync(); syncErr != nil {
			// Ignore sync errors on stdout/stderr (common on some platforms)
			// Only log if it's a real error
			a.logger.Warn("Logger sync returned error (may be harmless)", zap.Error(syncErr))
		}
	}

	a.logger.Info("Service shutdown complete")

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	return nil
}
