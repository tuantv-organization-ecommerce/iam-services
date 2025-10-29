package model

import (
	"errors"
	"time"
)

var (
	// ErrInvalidUsername indicates username validation failed
	ErrInvalidUsername = errors.New("invalid username")
	// ErrInvalidEmail indicates email validation failed
	ErrInvalidEmail = errors.New("invalid email")
	// ErrUserAlreadyActive indicates user is already in active state
	ErrUserAlreadyActive = errors.New("user is already active")
	// ErrUserAlreadyInactive indicates user is already in inactive state
	ErrUserAlreadyInactive = errors.New("user is already inactive")
	// ErrPasswordTooShort indicates password does not meet minimum length
	ErrPasswordTooShort = errors.New("password is too short")
	// ErrEmptyPassword indicates password field is empty
	ErrEmptyPassword = errors.New("password cannot be empty")
)

// User represents a user aggregate root in the domain
type User struct {
	id           string
	username     string
	email        string
	passwordHash string
	fullName     string
	isActive     bool
	roles        []string
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUser creates a new User entity
func NewUser(id, username, email, fullName string) *User {
	now := time.Now()
	return &User{
		id:        id,
		username:  username,
		email:     email,
		fullName:  fullName,
		isActive:  true,
		roles:     make([]string, 0),
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructUser reconstructs a User from persistence
func ReconstructUser(id, username, email, passwordHash, fullName string, isActive bool, createdAt, updatedAt time.Time) *User {
	return &User{
		id:           id,
		username:     username,
		email:        email,
		passwordHash: passwordHash,
		fullName:     fullName,
		isActive:     isActive,
		roles:        make([]string, 0),
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters

// ID returns the user's unique identifier
func (u *User) ID() string { return u.id }

// Username returns the user's username
func (u *User) Username() string     { return u.username }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) FullName() string     { return u.fullName }
func (u *User) IsActive() bool       { return u.isActive }
func (u *User) Roles() []string      { return u.roles }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// SetPasswordHash sets the hashed password
func (u *User) SetPasswordHash(hash string) error {
	if hash == "" {
		return ErrEmptyPassword
	}
	u.passwordHash = hash
	u.updatedAt = time.Now()
	return nil
}

// Activate activates the user account
func (u *User) Activate() error {
	if u.isActive {
		return ErrUserAlreadyActive
	}
	u.isActive = true
	u.updatedAt = time.Now()
	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() error {
	if !u.isActive {
		return ErrUserAlreadyInactive
	}
	u.isActive = false
	u.updatedAt = time.Now()
	return nil
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(fullName string) {
	u.fullName = fullName
	u.updatedAt = time.Now()
}

// AddRole adds a role to the user
func (u *User) AddRole(role string) {
	for _, r := range u.roles {
		if r == role {
			return // Already has role
		}
	}
	u.roles = append(u.roles, role)
}

// RemoveRole removes a role from the user
func (u *User) RemoveRole(role string) {
	for i, r := range u.roles {
		if r == role {
			u.roles = append(u.roles[:i], u.roles[i+1:]...)
			return
		}
	}
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role string) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

// Validate validates the user entity
func (u *User) Validate() error {
	if u.username == "" || len(u.username) < 3 {
		return ErrInvalidUsername
	}
	if u.email == "" || !isValidEmail(u.email) {
		return ErrInvalidEmail
	}
	return nil
}

// Helper function to validate email (basic validation)
func isValidEmail(email string) bool {
	// Basic email validation
	return len(email) > 3 && len(email) < 255
}
