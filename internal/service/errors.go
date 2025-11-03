package service

import (
	"fmt"

	"go.uber.org/zap"
)

// ErrorCode represents a unique error code for tracking and logging
type ErrorCode string

// Error codes for Auth Service
const (
	// Authentication errors (AUTH-xxx)
	ErrCodeInvalidCredentials     ErrorCode = "AUTH-001"
	ErrCodeUserNotFound           ErrorCode = "AUTH-002"
	ErrCodeUserInactive           ErrorCode = "AUTH-003"
	ErrCodeInvalidToken           ErrorCode = "AUTH-004"
	ErrCodeTokenExpired           ErrorCode = "AUTH-005"
	ErrCodeTokenGenerationFailed  ErrorCode = "AUTH-006"
	ErrCodePasswordHashFailed     ErrorCode = "AUTH-007"
	ErrCodeUserCreationFailed     ErrorCode = "AUTH-008"
	ErrCodeTokenRevocationFailed  ErrorCode = "AUTH-009"
	ErrCodeInvalidInput           ErrorCode = "AUTH-010"
	ErrCodeGetUserRolesFailed     ErrorCode = "AUTH-011"
)

// Error codes for Authorization Service
const (
	// Authorization errors (AUTHZ-xxx)
	ErrCodeRoleAssignmentFailed   ErrorCode = "AUTHZ-001"
	ErrCodeRoleRemovalFailed      ErrorCode = "AUTHZ-002"
	ErrCodeRoleNotFound           ErrorCode = "AUTHZ-003"
	ErrCodePermissionCheckFailed  ErrorCode = "AUTHZ-004"
	ErrCodeGetRolesFailed         ErrorCode = "AUTHZ-005"
	ErrCodeGetPermissionsFailed   ErrorCode = "AUTHZ-006"
	ErrCodeInvalidParameters      ErrorCode = "AUTHZ-007"
)

// Error codes for Role Service
const (
	// Role management errors (ROLE-xxx)
	ErrCodeRoleCreationFailed     ErrorCode = "ROLE-001"
	ErrCodeRoleUpdateFailed       ErrorCode = "ROLE-002"
	ErrCodeRoleDeletionFailed     ErrorCode = "ROLE-003"
	ErrCodeRoleGetFailed          ErrorCode = "ROLE-004"
	ErrCodeRoleListFailed         ErrorCode = "ROLE-005"
	ErrCodePermissionAssignFailed ErrorCode = "ROLE-006"
	ErrCodePermissionRemoveFailed ErrorCode = "ROLE-007"
)

// Error codes for Permission Service
const (
	// Permission management errors (PERM-xxx)
	ErrCodePermissionCreationFailed ErrorCode = "PERM-001"
	ErrCodePermissionUpdateFailed   ErrorCode = "PERM-002"
	ErrCodePermissionDeletionFailed ErrorCode = "PERM-003"
	ErrCodePermissionGetFailed      ErrorCode = "PERM-004"
	ErrCodePermissionListFailed     ErrorCode = "PERM-005"
	ErrCodePermissionNotFound       ErrorCode = "PERM-006"
)

// Error codes for Casbin Service
const (
	// Casbin policy errors (CASBIN-xxx)
	ErrCodePolicyAddFailed        ErrorCode = "CASBIN-001"
	ErrCodePolicyRemoveFailed     ErrorCode = "CASBIN-002"
	ErrCodePolicyEnforceFailed    ErrorCode = "CASBIN-003"
	ErrCodePolicyLoadFailed       ErrorCode = "CASBIN-004"
	ErrCodePolicyGetFailed        ErrorCode = "CASBIN-005"
	ErrCodeGroupingAddFailed      ErrorCode = "CASBIN-006"
	ErrCodeGroupingRemoveFailed   ErrorCode = "CASBIN-007"
)

// ServiceError represents a structured error with code for tracking
type ServiceError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError
func NewServiceError(code ErrorCode, message string, err error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// LogError logs an error with structured fields including error code
func LogError(logger *zap.Logger, err error, operation string, fields ...zap.Field) {
	if serviceErr, ok := err.(*ServiceError); ok {
		allFields := append([]zap.Field{
			zap.String("error_code", string(serviceErr.Code)),
			zap.String("operation", operation),
			zap.Error(serviceErr.Err),
		}, fields...)
		logger.Error(serviceErr.Message, allFields...)
	} else {
		allFields := append([]zap.Field{
			zap.String("operation", operation),
		}, fields...)
		logger.Error(err.Error(), allFields...)
	}
}

// LogErrorWithContext logs an error with additional context fields
func LogErrorWithContext(logger *zap.Logger, err error, operation, userID, resource string, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("user_id", userID),
		zap.String("resource", resource),
	}
	allFields := append(baseFields, fields...)
	LogError(logger, err, operation, allFields...)
}

