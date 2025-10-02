package entities

import "time"

// SharedUser represents the core user entity used across all modules
type SharedUser struct {
	// Core identification
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// Profile information
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`

	// Authentication
	Password string `json:"-"` // Never serialize password

	// Status fields
	IsActive        bool `json:"is_active"`
	IsEmailVerified bool `json:"is_email_verified"`
	IsPhoneVerified bool `json:"is_phone_verified"`

	// Timestamps
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewUser creates a new user with default values
func NewUser() *SharedUser {
	return &SharedUser{
		IsActive:        true,
		IsEmailVerified: false,
		IsPhoneVerified: false,
	}
}

// IsDeleted checks if the user is soft deleted
func (u *SharedUser) IsDeleted() bool {
	return u.DeletedAt != nil
}

// GetFullName returns the user's full name
func (u *SharedUser) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}
