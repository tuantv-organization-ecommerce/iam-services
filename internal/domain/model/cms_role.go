package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidCMSRole = errors.New("invalid CMS role")
	ErrEmptyTabs      = errors.New("CMS role must have at least one tab")
)

// CMSTab represents a CMS tab/section
type CMSTab string

const (
	CMSTabProduct   CMSTab = "product"
	CMSTabInventory CMSTab = "inventory"
	CMSTabOrder     CMSTab = "order"
	CMSTabUser      CMSTab = "user"
	CMSTabReport    CMSTab = "report"
	CMSTabSetting   CMSTab = "setting"
)

// CMSRole represents a CMS role entity with tab-based permissions
type CMSRole struct {
	id          string
	name        string
	description string
	tabs        []CMSTab
	createdAt   time.Time
	updatedAt   time.Time
}

// NewCMSRole creates a new CMSRole entity
func NewCMSRole(id, name, description string, tabs []CMSTab) *CMSRole {
	now := time.Now()
	return &CMSRole{
		id:          id,
		name:        name,
		description: description,
		tabs:        tabs,
		createdAt:   now,
		updatedAt:   now,
	}
}

// ReconstructCMSRole reconstructs a CMSRole from persistence
func ReconstructCMSRole(id, name, description string, tabs []CMSTab, createdAt, updatedAt time.Time) *CMSRole {
	return &CMSRole{
		id:          id,
		name:        name,
		description: description,
		tabs:        tabs,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters
func (c *CMSRole) ID() string           { return c.id }
func (c *CMSRole) Name() string         { return c.name }
func (c *CMSRole) Description() string  { return c.description }
func (c *CMSRole) Tabs() []CMSTab       { return c.tabs }
func (c *CMSRole) CreatedAt() time.Time { return c.createdAt }
func (c *CMSRole) UpdatedAt() time.Time { return c.updatedAt }

// UpdateDetails updates CMS role details
func (c *CMSRole) UpdateDetails(name, description string, tabs []CMSTab) {
	if name != "" {
		c.name = name
	}
	c.description = description
	if len(tabs) > 0 {
		c.tabs = tabs
	}
	c.updatedAt = time.Now()
}

// HasTab checks if role has access to a specific tab
func (c *CMSRole) HasTab(tab CMSTab) bool {
	for _, t := range c.tabs {
		if t == tab {
			return true
		}
	}
	return false
}

// AddTab adds a tab to the role
func (c *CMSRole) AddTab(tab CMSTab) {
	if c.HasTab(tab) {
		return
	}
	c.tabs = append(c.tabs, tab)
	c.updatedAt = time.Now()
}

// RemoveTab removes a tab from the role
func (c *CMSRole) RemoveTab(tab CMSTab) {
	for i, t := range c.tabs {
		if t == tab {
			c.tabs = append(c.tabs[:i], c.tabs[i+1:]...)
			c.updatedAt = time.Now()
			return
		}
	}
}

// Validate validates the CMS role entity
func (c *CMSRole) Validate() error {
	if c.name == "" {
		return ErrInvalidCMSRole
	}
	if len(c.tabs) == 0 {
		return ErrEmptyTabs
	}
	return nil
}
