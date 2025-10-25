package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidPermission = errors.New("invalid permission")
	ErrEmptyResource     = errors.New("resource cannot be empty")
	ErrEmptyAction       = errors.New("action cannot be empty")
)

// Permission represents a permission entity in the domain
type Permission struct {
	id          string
	name        string
	resource    string
	action      string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewPermission creates a new Permission entity
func NewPermission(id, name, resource, action, description string) *Permission {
	now := time.Now()
	return &Permission{
		id:          id,
		name:        name,
		resource:    resource,
		action:      action,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}
}

// ReconstructPermission reconstructs a Permission from persistence
func ReconstructPermission(id, name, resource, action, description string, createdAt, updatedAt time.Time) *Permission {
	return &Permission{
		id:          id,
		name:        name,
		resource:    resource,
		action:      action,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters
func (p *Permission) ID() string           { return p.id }
func (p *Permission) Name() string         { return p.name }
func (p *Permission) Resource() string     { return p.resource }
func (p *Permission) Action() string       { return p.action }
func (p *Permission) Description() string  { return p.description }
func (p *Permission) CreatedAt() time.Time { return p.createdAt }
func (p *Permission) UpdatedAt() time.Time { return p.updatedAt }

// Matches checks if this permission matches the given resource and action
func (p *Permission) Matches(resource, action string) bool {
	return p.resource == resource && p.action == action
}

// Validate validates the permission entity
func (p *Permission) Validate() error {
	if p.resource == "" {
		return ErrEmptyResource
	}
	if p.action == "" {
		return ErrEmptyAction
	}
	return nil
}
