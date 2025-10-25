package dto

// CheckAPIAccessRequest represents API access check input
type CheckAPIAccessRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	APIPath string `json:"api_path" validate:"required"`
	Method  string `json:"method" validate:"required"`
}

// CheckAPIAccessResponse represents API access check output
type CheckAPIAccessResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message"`
}

// CheckCMSAccessRequest represents CMS access check input
type CheckCMSAccessRequest struct {
	UserID string `json:"user_id" validate:"required"`
	CMSTab string `json:"cms_tab" validate:"required"`
	Action string `json:"action" validate:"required"`
}

// CheckCMSAccessResponse represents CMS access check output
type CheckCMSAccessResponse struct {
	Allowed        bool     `json:"allowed"`
	Message        string   `json:"message"`
	AccessibleTabs []string `json:"accessible_tabs"`
}

// CreateCMSRoleRequest represents CMS role creation input
type CreateCMSRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Tabs        []string `json:"tabs" validate:"required,min=1"`
}

// CreateCMSRoleResponse represents CMS role creation output
type CreateCMSRoleResponse struct {
	CMSRoleID string `json:"cms_role_id"`
	Message   string `json:"message"`
}

// CMSRoleDTO represents CMS role data transfer object
type CMSRoleDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tabs        []string `json:"tabs"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// APIResourceDTO represents API resource data transfer object
type APIResourceDTO struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	Service     string `json:"service"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
