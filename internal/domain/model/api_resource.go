package model

import (
	"errors"
	"strings"
	"time"
)

var (
	// ErrInvalidAPIResource indicates an invalid API resource entity
	ErrInvalidAPIResource = errors.New("invalid API resource")
	// ErrEmptyPath indicates path field is empty
	ErrEmptyPath = errors.New("path cannot be empty")
	// ErrEmptyMethod indicates method field is empty
	ErrEmptyMethod = errors.New("method cannot be empty")
	// ErrInvalidMethod indicates an invalid HTTP method
	ErrInvalidMethod = errors.New("invalid HTTP method")
)

// HTTPMethod represents valid HTTP methods
type HTTPMethod string

const (
	// MethodGET represents HTTP GET method
	MethodGET HTTPMethod = "GET"
	// MethodPOST represents HTTP POST method
	MethodPOST HTTPMethod = "POST"
	// MethodPUT represents HTTP PUT method
	MethodPUT HTTPMethod = "PUT"
	// MethodDELETE represents HTTP DELETE method
	MethodDELETE HTTPMethod = "DELETE"
	// MethodPATCH represents HTTP PATCH method
	MethodPATCH HTTPMethod = "PATCH"
)

// APIResource represents an API endpoint resource
type APIResource struct {
	id          string
	path        string
	method      HTTPMethod
	service     string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewAPIResource creates a new APIResource entity
func NewAPIResource(id, path string, method HTTPMethod, service, description string) *APIResource {
	now := time.Now()
	return &APIResource{
		id:          id,
		path:        path,
		method:      method,
		service:     service,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}
}

// ReconstructAPIResource reconstructs an APIResource from persistence
func ReconstructAPIResource(id, path string, method HTTPMethod, service, description string, createdAt, updatedAt time.Time) *APIResource {
	return &APIResource{
		id:          id,
		path:        path,
		method:      method,
		service:     service,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters
func (a *APIResource) ID() string           { return a.id }
func (a *APIResource) Path() string         { return a.path }
func (a *APIResource) Method() HTTPMethod   { return a.method }
func (a *APIResource) Service() string      { return a.service }
func (a *APIResource) Description() string  { return a.description }
func (a *APIResource) CreatedAt() time.Time { return a.createdAt }
func (a *APIResource) UpdatedAt() time.Time { return a.updatedAt }

// UpdateDetails updates API resource details
func (a *APIResource) UpdateDetails(description string) {
	a.description = description
	a.updatedAt = time.Now()
}

// Matches checks if this resource matches the given path and method
func (a *APIResource) Matches(path string, method HTTPMethod) bool {
	return a.path == path && a.method == method
}

// Validate validates the API resource entity
func (a *APIResource) Validate() error {
	if a.path == "" {
		return ErrEmptyPath
	}
	if a.method == "" {
		return ErrEmptyMethod
	}
	if !isValidHTTPMethod(string(a.method)) {
		return ErrInvalidMethod
	}
	return nil
}

// Helper function to validate HTTP method
func isValidHTTPMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	methodUpper := strings.ToUpper(method)
	for _, m := range validMethods {
		if m == methodUpper {
			return true
		}
	}
	return false
}
