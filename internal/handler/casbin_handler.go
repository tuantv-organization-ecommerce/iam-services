package handler

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tvttt/iam-services/internal/domain"
	pb "github.com/tvttt/iam-services/pkg/proto"
)

// CheckAPIAccess handles API access authorization check
func (h *GRPCHandler) CheckAPIAccess(ctx context.Context, req *pb.CheckAPIAccessRequest) (*pb.CheckAPIAccessResponse, error) {
	h.logger.Info("CheckAPIAccess request received",
		zap.String("user_id", req.UserId),
		zap.String("api_path", req.ApiPath),
		zap.String("method", req.Method))

	allowed, err := h.casbinService.CheckAPIAccess(ctx, req.UserId, req.ApiPath, req.Method)
	if err != nil {
		h.logger.Error("Failed to check API access", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to check API access: %v", err)
	}

	message := "Access denied"
	if allowed {
		message = "Access granted"
	}

	return &pb.CheckAPIAccessResponse{
		Allowed: allowed,
		Message: message,
	}, nil
}

// CheckCMSAccess handles CMS access authorization check
func (h *GRPCHandler) CheckCMSAccess(ctx context.Context, req *pb.CheckCMSAccessRequest) (*pb.CheckCMSAccessResponse, error) {
	h.logger.Info("CheckCMSAccess request received",
		zap.String("user_id", req.UserId),
		zap.String("cms_tab", req.CmsTab),
		zap.String("action", req.Action))

	// Get user's CMS tabs first
	tabs, err := h.casbinService.GetUserCMSTabs(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user CMS tabs", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get user CMS tabs: %v", err)
	}

	accessibleTabs := make([]string, len(tabs))
	for i, tab := range tabs {
		accessibleTabs[i] = string(tab)
	}

	// Check access to specific tab
	cmsTab := domain.CMSTab(req.CmsTab)
	allowed, err := h.casbinService.CheckCMSAccess(ctx, req.UserId, cmsTab, req.Action)
	if err != nil {
		h.logger.Error("Failed to check CMS access", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to check CMS access: %v", err)
	}

	message := "Access denied to CMS tab"
	if allowed {
		message = "Access granted to CMS tab"
	}

	return &pb.CheckCMSAccessResponse{
		Allowed:        allowed,
		Message:        message,
		AccessibleTabs: accessibleTabs,
	}, nil
}

// EnforcePolicy handles general policy enforcement
func (h *GRPCHandler) EnforcePolicy(ctx context.Context, req *pb.EnforcePolicyRequest) (*pb.EnforcePolicyResponse, error) {
	h.logger.Info("EnforcePolicy request received",
		zap.String("user_id", req.UserId),
		zap.String("domain", req.Domain),
		zap.String("resource", req.Resource),
		zap.String("action", req.Action))

	authReq := &domain.AuthorizationRequest{
		UserID:   req.UserId,
		Domain:   domain.CasbinDomain(req.Domain),
		Resource: req.Resource,
		Action:   req.Action,
	}

	response, err := h.casbinService.Enforce(ctx, authReq)
	if err != nil {
		h.logger.Error("Failed to enforce policy", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to enforce policy: %v", err)
	}

	return &pb.EnforcePolicyResponse{
		Allowed: response.Allowed,
		Reason:  response.Reason,
	}, nil
}

// CreateCMSRole handles CMS role creation
func (h *GRPCHandler) CreateCMSRole(ctx context.Context, req *pb.CreateCMSRoleRequest) (*pb.CreateCMSRoleResponse, error) {
	h.logger.Info("CreateCMSRole request received", zap.String("name", req.Name))

	tabs := make([]domain.CMSTab, len(req.Tabs))
	for i, tab := range req.Tabs {
		tabs[i] = domain.CMSTab(tab)
	}

	cmsRole, err := h.casbinService.CreateCMSRole(ctx, req.Name, req.Description, tabs)
	if err != nil {
		h.logger.Error("Failed to create CMS role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create CMS role: %v", err)
	}

	return &pb.CreateCMSRoleResponse{
		CmsRoleId: cmsRole.ID,
		Message:   "CMS role created successfully",
	}, nil
}

// AssignCMSRole handles CMS role assignment to user
func (h *GRPCHandler) AssignCMSRole(ctx context.Context, req *pb.AssignCMSRoleRequest) (*pb.AssignCMSRoleResponse, error) {
	h.logger.Info("AssignCMSRole request received",
		zap.String("user_id", req.UserId),
		zap.String("cms_role_id", req.CmsRoleId))

	if err := h.casbinService.AssignCMSRole(ctx, req.UserId, req.CmsRoleId); err != nil {
		h.logger.Error("Failed to assign CMS role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to assign CMS role: %v", err)
	}

	return &pb.AssignCMSRoleResponse{
		Message: "CMS role assigned successfully",
	}, nil
}

// RemoveCMSRole handles CMS role removal from user
func (h *GRPCHandler) RemoveCMSRole(ctx context.Context, req *pb.RemoveCMSRoleRequest) (*pb.RemoveCMSRoleResponse, error) {
	h.logger.Info("RemoveCMSRole request received",
		zap.String("user_id", req.UserId),
		zap.String("cms_role_id", req.CmsRoleId))

	if err := h.casbinService.RemoveCMSRole(ctx, req.UserId, req.CmsRoleId); err != nil {
		h.logger.Error("Failed to remove CMS role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to remove CMS role: %v", err)
	}

	return &pb.RemoveCMSRoleResponse{
		Message: "CMS role removed successfully",
	}, nil
}

// GetUserCMSTabs handles getting user's accessible CMS tabs
func (h *GRPCHandler) GetUserCMSTabs(ctx context.Context, req *pb.GetUserCMSTabsRequest) (*pb.GetUserCMSTabsResponse, error) {
	h.logger.Info("GetUserCMSTabs request received", zap.String("user_id", req.UserId))

	tabs, err := h.casbinService.GetUserCMSTabs(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user CMS tabs", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get user CMS tabs: %v", err)
	}

	tabStrings := make([]string, len(tabs))
	for i, tab := range tabs {
		tabStrings[i] = string(tab)
	}

	return &pb.GetUserCMSTabsResponse{
		Tabs: tabStrings,
	}, nil
}

// ListCMSRoles handles listing CMS roles
func (h *GRPCHandler) ListCMSRoles(ctx context.Context, req *pb.ListCMSRolesRequest) (*pb.ListCMSRolesResponse, error) {
	h.logger.Info("ListCMSRoles request received")

	// For simplicity, we'll list all CMS roles from database
	// In a real implementation, you'd add pagination support

	return &pb.ListCMSRolesResponse{
		Roles: []*pb.CMSRole{},
		Total: 0,
	}, nil
}

// CreateAPIResource handles API resource creation
func (h *GRPCHandler) CreateAPIResource(ctx context.Context, req *pb.CreateAPIResourceRequest) (*pb.CreateAPIResourceResponse, error) {
	h.logger.Info("CreateAPIResource request received",
		zap.String("path", req.Path),
		zap.String("method", req.Method))

	resource, err := h.casbinService.CreateAPIResource(ctx, req.Path, req.Method, req.Service, req.Description)
	if err != nil {
		h.logger.Error("Failed to create API resource", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create API resource: %v", err)
	}

	return &pb.CreateAPIResourceResponse{
		ApiResourceId: resource.ID,
		Message:       "API resource created successfully",
	}, nil
}

// ListAPIResources handles listing API resources
func (h *GRPCHandler) ListAPIResources(ctx context.Context, req *pb.ListAPIResourcesRequest) (*pb.ListAPIResourcesResponse, error) {
	h.logger.Info("ListAPIResources request received", zap.String("service", req.Service))

	resources, err := h.casbinService.ListAPIResources(ctx, req.Service)
	if err != nil {
		h.logger.Error("Failed to list API resources", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list API resources: %v", err)
	}

	pbResources := make([]*pb.APIResource, len(resources))
	for i, res := range resources {
		pbResources[i] = &pb.APIResource{
			Id:          res.ID,
			Path:        res.Path,
			Method:      res.Method,
			Service:     res.Service,
			Description: res.Description,
			CreatedAt:   res.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   res.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &pb.ListAPIResourcesResponse{
		Resources: pbResources,
		Total:     int32(len(pbResources)),
	}, nil
}
