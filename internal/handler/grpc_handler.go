package handler

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tvttt/iam-services/internal/service"
	pb "github.com/tvttt/iam-services/pkg/proto"
)

// GRPCHandler implements the gRPC IAMService
type GRPCHandler struct {
	pb.UnimplementedIAMServiceServer
	authService   service.AuthService
	authzService  service.AuthorizationService
	roleService   service.RoleService
	permService   service.PermissionService
	casbinService service.CasbinService
	logger        *zap.Logger
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(
	authService service.AuthService,
	authzService service.AuthorizationService,
	roleService service.RoleService,
	permService service.PermissionService,
	casbinService service.CasbinService,
	logger *zap.Logger,
) *GRPCHandler {
	return &GRPCHandler{
		authService:   authService,
		authzService:  authzService,
		roleService:   roleService,
		permService:   permService,
		casbinService: casbinService,
		logger:        logger,
	}
}

// Register handles user registration
func (h *GRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	h.logger.Info("Register request received", zap.String("username", req.Username))

	user, err := h.authService.Register(ctx, req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User registered successfully",
	}, nil
}

// Login handles user login
func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	h.logger.Info("Login request received", zap.String("username", req.Username))

	user, tokenPair, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &pb.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// RefreshToken handles token refresh
func (h *GRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	h.logger.Info("RefreshToken request received")

	tokenPair, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh token", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Logout handles user logout
func (h *GRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	h.logger.Info("Logout request received", zap.String("user_id", req.UserId))

	if err := h.authService.Logout(ctx, req.UserId); err != nil {
		h.logger.Error("Failed to logout", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to logout")
	}

	return &pb.LogoutResponse{
		Message: "Logged out successfully",
	}, nil
}

// VerifyToken handles token verification
func (h *GRPCHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	h.logger.Info("VerifyToken request received")

	userID, roles, err := h.authService.VerifyToken(ctx, req.Token)
	if err != nil {
		h.logger.Error("Failed to verify token", zap.Error(err))
		return &pb.VerifyTokenResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Valid:   true,
		UserId:  userID,
		Roles:   roles,
		Message: "Token is valid",
	}, nil
}

// AssignRole handles role assignment to user
func (h *GRPCHandler) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	h.logger.Info("AssignRole request received", zap.String("user_id", req.UserId), zap.String("role_id", req.RoleId))

	if err := h.authzService.AssignRole(ctx, req.UserId, req.RoleId); err != nil {
		h.logger.Error("Failed to assign role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to assign role: %v", err)
	}

	return &pb.AssignRoleResponse{
		Message: "Role assigned successfully",
	}, nil
}

// RemoveRole handles role removal from user
func (h *GRPCHandler) RemoveRole(ctx context.Context, req *pb.RemoveRoleRequest) (*pb.RemoveRoleResponse, error) {
	h.logger.Info("RemoveRole request received", zap.String("user_id", req.UserId), zap.String("role_id", req.RoleId))

	if err := h.authzService.RemoveRole(ctx, req.UserId, req.RoleId); err != nil {
		h.logger.Error("Failed to remove role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to remove role: %v", err)
	}

	return &pb.RemoveRoleResponse{
		Message: "Role removed successfully",
	}, nil
}

// GetUserRoles handles getting user roles
func (h *GRPCHandler) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	h.logger.Info("GetUserRoles request received", zap.String("user_id", req.UserId))

	roles, err := h.authzService.GetUserRoles(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user roles", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get user roles: %v", err)
	}

	pbRoles := make([]*pb.Role, len(roles))
	for i, role := range roles {
		pbRoles[i] = domainRoleToPB(role)
	}

	return &pb.GetUserRolesResponse{
		Roles: pbRoles,
	}, nil
}

// CheckPermission handles permission checking
func (h *GRPCHandler) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	h.logger.Info("CheckPermission request received",
		zap.String("user_id", req.UserId),
		zap.String("resource", req.Resource),
		zap.String("action", req.Action))

	allowed, err := h.authzService.CheckPermission(ctx, req.UserId, req.Resource, req.Action)
	if err != nil {
		h.logger.Error("Failed to check permission", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to check permission: %v", err)
	}

	message := "Permission denied"
	if allowed {
		message = "Permission granted"
	}

	return &pb.CheckPermissionResponse{
		Allowed: allowed,
		Message: message,
	}, nil
}

// CreateRole handles role creation
func (h *GRPCHandler) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	h.logger.Info("CreateRole request received", zap.String("name", req.Name))

	role, err := h.roleService.CreateRole(ctx, req.Name, req.Description, req.PermissionIds)
	if err != nil {
		h.logger.Error("Failed to create role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create role: %v", err)
	}

	return &pb.CreateRoleResponse{
		RoleId:  role.ID,
		Message: "Role created successfully",
	}, nil
}

// UpdateRole handles role update
func (h *GRPCHandler) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	h.logger.Info("UpdateRole request received", zap.String("role_id", req.RoleId))

	if err := h.roleService.UpdateRole(ctx, req.RoleId, req.Name, req.Description, req.PermissionIds); err != nil {
		h.logger.Error("Failed to update role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update role: %v", err)
	}

	return &pb.UpdateRoleResponse{
		Message: "Role updated successfully",
	}, nil
}

// DeleteRole handles role deletion
func (h *GRPCHandler) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	h.logger.Info("DeleteRole request received", zap.String("role_id", req.RoleId))

	if err := h.roleService.DeleteRole(ctx, req.RoleId); err != nil {
		h.logger.Error("Failed to delete role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete role: %v", err)
	}

	return &pb.DeleteRoleResponse{
		Message: "Role deleted successfully",
	}, nil
}

// GetRole handles getting a single role
func (h *GRPCHandler) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	h.logger.Info("GetRole request received", zap.String("role_id", req.RoleId))

	role, err := h.roleService.GetRole(ctx, req.RoleId)
	if err != nil {
		h.logger.Error("Failed to get role", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "role not found")
	}

	return &pb.GetRoleResponse{
		Role: domainRoleToPB(role),
	}, nil
}

// ListRoles handles listing roles
func (h *GRPCHandler) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	h.logger.Info("ListRoles request received")

	page := int(req.Page)
	pageSize := int(req.PageSize)

	roles, total, err := h.roleService.ListRoles(ctx, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list roles", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list roles: %v", err)
	}

	pbRoles := make([]*pb.Role, len(roles))
	for i, role := range roles {
		pbRoles[i] = domainRoleToPB(role)
	}

	return &pb.ListRolesResponse{
		Roles: pbRoles,
		Total: int32(total),
	}, nil
}

// CreatePermission handles permission creation
func (h *GRPCHandler) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.CreatePermissionResponse, error) {
	h.logger.Info("CreatePermission request received", zap.String("name", req.Name))

	permission, err := h.permService.CreatePermission(ctx, req.Name, req.Resource, req.Action, req.Description)
	if err != nil {
		h.logger.Error("Failed to create permission", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create permission: %v", err)
	}

	return &pb.CreatePermissionResponse{
		PermissionId: permission.ID,
		Message:      "Permission created successfully",
	}, nil
}

// DeletePermission handles permission deletion
func (h *GRPCHandler) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*pb.DeletePermissionResponse, error) {
	h.logger.Info("DeletePermission request received", zap.String("permission_id", req.PermissionId))

	if err := h.permService.DeletePermission(ctx, req.PermissionId); err != nil {
		h.logger.Error("Failed to delete permission", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete permission: %v", err)
	}

	return &pb.DeletePermissionResponse{
		Message: "Permission deleted successfully",
	}, nil
}

// ListPermissions handles listing permissions
func (h *GRPCHandler) ListPermissions(ctx context.Context, req *pb.ListPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	h.logger.Info("ListPermissions request received")

	page := int(req.Page)
	pageSize := int(req.PageSize)

	permissions, total, err := h.permService.ListPermissions(ctx, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list permissions", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list permissions: %v", err)
	}

	pbPermissions := make([]*pb.Permission, len(permissions))
	for i, perm := range permissions {
		pbPermissions[i] = domainPermissionToPB(perm)
	}

	return &pb.ListPermissionsResponse{
		Permissions: pbPermissions,
		Total:       int32(total),
	}, nil
}
