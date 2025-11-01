package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tvttt/iam-services/internal/application/dto"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/service"
)

// GinHandler handles HTTP requests using Gin framework
type GinHandler struct {
	authService   service.AuthService
	authzService  service.AuthorizationService
	roleService   service.RoleService
	permService   service.PermissionService
	casbinService service.CasbinService
	logger        *zap.Logger
}

// NewGinHandler creates a new Gin HTTP handler
func NewGinHandler(
	authService service.AuthService,
	authzService service.AuthorizationService,
	roleService service.RoleService,
	permService service.PermissionService,
	casbinService service.CasbinService,
	logger *zap.Logger,
) *GinHandler {
	return &GinHandler{
		authService:   authService,
		authzService:  authzService,
		roleService:   roleService,
		permService:   permService,
		casbinService: casbinService,
		logger:        logger,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// sendError sends an error response
func (h *GinHandler) sendError(c *gin.Context, code int, err error, message string) {
	h.logger.Error(message, zap.Error(err))
	c.JSON(code, ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Code:    code,
	})
}

// sendSuccess sends a success response
func (h *GinHandler) sendSuccess(c *gin.Context, code int, data interface{}, message string) {
	response := SuccessResponse{
		Data:    data,
		Message: message,
	}
	c.JSON(code, response)
}

// Register handles user registration
func (h *GinHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to register user")
		return
	}

	response := dto.RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User registered successfully",
	}

	h.sendSuccess(c, http.StatusCreated, response, "")
}

// Login handles user login
func (h *GinHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	user, tokenPair, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err, "Invalid credentials")
		return
	}

	response := dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &dto.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}

	h.sendSuccess(c, http.StatusOK, response, "")
}

// RefreshToken handles token refresh
func (h *GinHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err, "Invalid refresh token")
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}

	h.sendSuccess(c, http.StatusOK, response, "")
}

// Logout handles user logout
func (h *GinHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	err := h.authService.Logout(c.Request.Context(), req.UserID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to logout")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "Logout successful")
}

// VerifyToken handles token verification
func (h *GinHandler) VerifyToken(c *gin.Context) {
	var req dto.VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	userID, roles, err := h.authService.VerifyToken(c.Request.Context(), req.Token)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	response := dto.VerifyTokenResponse{
		Valid:   true,
		UserID:  userID,
		Roles:   roles,
		Message: "Token is valid",
	}

	h.sendSuccess(c, http.StatusOK, response, "")
}

// Health handles health check
func (h *GinHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "iam-service",
	})
}

// Role Management Handlers

// CreateRole handles role creation
func (h *GinHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name          string   `json:"name" binding:"required"`
		Description   string   `json:"description"`
		PermissionIDs []string `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create role")
		return
	}

	h.sendSuccess(c, http.StatusCreated, gin.H{
		"role_id": role.ID,
		"name":    role.Name,
	}, "Role created successfully")
}

// ListRoles handles listing all roles
func (h *GinHandler) ListRoles(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil {
			pageSize = parsed
		}
	}

	roles, total, err := h.roleService.ListRoles(c.Request.Context(), page, pageSize)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to list roles")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"roles": roles,
		"total": total,
	}, "")
}

// GetRole handles getting a role by ID
func (h *GinHandler) GetRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Role ID is required")
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), roleID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Role not found")
		return
	}

	h.sendSuccess(c, http.StatusOK, role, "")
}

// UpdateRole handles role update
func (h *GinHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	var req struct {
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		PermissionIDs []string `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	err := h.roleService.UpdateRole(c.Request.Context(), roleID, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to update role")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "Role updated successfully")
}

// DeleteRole handles role deletion
func (h *GinHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Role ID is required")
		return
	}

	err := h.roleService.DeleteRole(c.Request.Context(), roleID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to delete role")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "Role deleted successfully")
}

// AssignRole handles assigning role to user
func (h *GinHandler) AssignRole(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		RoleID string `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	err := h.authzService.AssignRole(c.Request.Context(), req.UserID, req.RoleID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to assign role")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "Role assigned successfully")
}

// RemoveRole handles removing role from user
func (h *GinHandler) RemoveRole(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		RoleID string `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	err := h.authzService.RemoveRole(c.Request.Context(), req.UserID, req.RoleID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to remove role")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "Role removed successfully")
}

// GetUserRoles handles getting user's roles
func (h *GinHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "User ID is required")
		return
	}

	roles, err := h.authzService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get user roles")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{"roles": roles}, "")
}

// Permission Management Handlers

// CreatePermission handles permission creation
func (h *GinHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Resource    string `json:"resource" binding:"required"`
		Action      string `json:"action" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	permission, err := h.permService.CreatePermission(c.Request.Context(), req.Name, req.Resource, req.Action, req.Description)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create permission")
		return
	}

	h.sendSuccess(c, http.StatusCreated, gin.H{
		"permission_id": permission.ID,
		"name":          permission.Name,
	}, "Permission created successfully")
}

// ListPermissions handles listing all permissions
func (h *GinHandler) ListPermissions(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil {
			pageSize = parsed
		}
	}

	permissions, total, err := h.permService.ListPermissions(c.Request.Context(), page, pageSize)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to list permissions")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"permissions": permissions,
		"total":       total,
	}, "")
}

// DeletePermission handles permission deletion
func (h *GinHandler) DeletePermission(c *gin.Context) {
	permissionID := c.Param("id")
	if permissionID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Permission ID is required")
		return
	}

	err := h.permService.DeletePermission(c.Request.Context(), permissionID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to delete permission")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "Permission deleted successfully")
}

// CheckPermission handles permission check
func (h *GinHandler) CheckPermission(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		Resource string `json:"resource" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	allowed, err := h.authzService.CheckPermission(c.Request.Context(), req.UserID, req.Resource, req.Action)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to check permission")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"allowed": allowed,
		"message": "Permission check completed",
	}, "")
}

// Casbin Handlers

// CheckAPIAccess handles API access check
func (h *GinHandler) CheckAPIAccess(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id" binding:"required"`
		APIPath string `json:"api_path" binding:"required"`
		Method  string `json:"method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	allowed, err := h.casbinService.CheckAPIAccess(c.Request.Context(), req.UserID, req.APIPath, req.Method)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to check API access")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"allowed": allowed,
		"message": "API access check completed",
	}, "")
}

// CheckCMSAccess handles CMS access check
func (h *GinHandler) CheckCMSAccess(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		CMSTab string `json:"cms_tab" binding:"required"`
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Get user's CMS tabs first
	tabs, err := h.casbinService.GetUserCMSTabs(c.Request.Context(), req.UserID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get user CMS tabs")
		return
	}

	// Convert tabs to strings
	tabStrings := make([]string, len(tabs))
	for i, tab := range tabs {
		tabStrings[i] = string(tab)
	}

	// Check if user has access to requested tab
	allowed, err := h.casbinService.CheckCMSAccess(c.Request.Context(), req.UserID, domain.CMSTab(req.CMSTab), req.Action)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to check CMS access")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"allowed":         allowed,
		"accessible_tabs": tabStrings,
		"message":         "CMS access check completed",
	}, "")
}

// EnforcePolicy handles policy enforcement
func (h *GinHandler) EnforcePolicy(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		Domain   string `json:"domain" binding:"required"`
		Resource string `json:"resource" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	authReq := &domain.AuthorizationRequest{
		UserID:   req.UserID,
		Domain:   domain.CasbinDomain(req.Domain),
		Resource: req.Resource,
		Action:   req.Action,
	}

	response, err := h.casbinService.Enforce(c.Request.Context(), authReq)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to enforce policy")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"allowed": response.Allowed,
		"reason":  response.Reason,
	}, "")
}

// CMS Handlers

// CreateCMSRole handles CMS role creation
func (h *GinHandler) CreateCMSRole(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Tabs        []string `json:"tabs" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Convert string tabs to domain.CMSTab
	tabs := make([]domain.CMSTab, len(req.Tabs))
	for i, tab := range req.Tabs {
		tabs[i] = domain.CMSTab(tab)
	}

	role, err := h.casbinService.CreateCMSRole(c.Request.Context(), req.Name, req.Description, tabs)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create CMS role")
		return
	}

	h.sendSuccess(c, http.StatusCreated, gin.H{
		"cms_role_id": role.ID,
		"name":        role.Name,
	}, "CMS role created successfully")
}

// ListCMSRoles handles listing CMS roles
func (h *GinHandler) ListCMSRoles(c *gin.Context) {
	// Note: Full implementation pending - returning empty list for now
	h.sendSuccess(c, http.StatusOK, gin.H{
		"roles": []interface{}{},
		"total": 0,
	}, "")
}

// AssignCMSRole handles assigning CMS role to user
func (h *GinHandler) AssignCMSRole(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id" binding:"required"`
		CMSRoleID string `json:"cms_role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	err := h.casbinService.AssignCMSRole(c.Request.Context(), req.UserID, req.CMSRoleID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to assign CMS role")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "CMS role assigned successfully")
}

// RemoveCMSRole handles removing CMS role from user
func (h *GinHandler) RemoveCMSRole(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id" binding:"required"`
		CMSRoleID string `json:"cms_role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	err := h.casbinService.RemoveCMSRole(c.Request.Context(), req.UserID, req.CMSRoleID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to remove CMS role")
		return
	}

	h.sendSuccess(c, http.StatusOK, nil, "CMS role removed successfully")
}

// GetUserCMSTabs handles getting user's CMS tabs
func (h *GinHandler) GetUserCMSTabs(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "User ID is required")
		return
	}

	tabs, err := h.casbinService.GetUserCMSTabs(c.Request.Context(), userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get user CMS tabs")
		return
	}

	// Convert tabs to strings
	tabStrings := make([]string, len(tabs))
	for i, tab := range tabs {
		tabStrings[i] = string(tab)
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"tabs": tabStrings,
	}, "")
}

// API Resource Handlers

// CreateAPIResource handles API resource creation
func (h *GinHandler) CreateAPIResource(c *gin.Context) {
	var req struct {
		Path        string `json:"path" binding:"required"`
		Method      string `json:"method" binding:"required"`
		Service     string `json:"service" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	resource, err := h.casbinService.CreateAPIResource(c.Request.Context(), req.Path, req.Method, req.Service, req.Description)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create API resource")
		return
	}

	h.sendSuccess(c, http.StatusCreated, gin.H{
		"api_resource_id": resource.ID,
		"path":            resource.Path,
		"method":          resource.Method,
	}, "API resource created successfully")
}

// ListAPIResources handles listing API resources
func (h *GinHandler) ListAPIResources(c *gin.Context) {
	service := c.Query("service")

	resources, err := h.casbinService.ListAPIResources(c.Request.Context(), service)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to list API resources")
		return
	}

	h.sendSuccess(c, http.StatusOK, gin.H{
		"resources": resources,
		"total":     len(resources),
	}, "")
}
