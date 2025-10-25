package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryUnaryInterceptor recovers from panics in unary RPCs
func RecoveryUnaryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Panic recovery
		defer func() {
			if r := recover(); r != nil {
				// Log panic with stack trace
				logger.Error("Panic recovered in gRPC handler",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)

				// Return error to client
				err = status.Errorf(codes.Internal, "Internal server error: %v", r)
			}
		}()

		// Execute handler
		return handler(ctx, req)
	}
}

// RecoveryStreamInterceptor recovers from panics in stream RPCs
func RecoveryStreamInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		// Panic recovery
		defer func() {
			if r := recover(); r != nil {
				// Log panic with stack trace
				logger.Error("Panic recovered in gRPC stream handler",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)

				// Return error to client
				err = status.Errorf(codes.Internal, "Internal server error: %v", r)
			}
		}()

		// Execute handler
		return handler(srv, ss)
	}
}

// RecoverGoroutine wraps a goroutine with panic recovery
func RecoverGoroutine(logger *zap.Logger, name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered in goroutine",
					zap.String("goroutine", name),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)
			}
		}()

		fn()
	}()
}

// RecoverFunc wraps any function with panic recovery
func RecoverFunc(logger *zap.Logger, name string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered in function",
				zap.String("function", name),
				zap.Any("panic", r),
				zap.String("stack", string(debug.Stack())),
			)
			err = fmt.Errorf("panic in %s: %v", name, r)
		}
	}()

	return fn()
}
