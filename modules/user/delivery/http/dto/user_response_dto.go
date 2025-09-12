package dto

import "time"

type UserResponse struct {
	ID              int64      `json:"id"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	PhoneNumber     string     `json:"phone_number"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	IsActive        bool       `json:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified"`
	IsPhoneVerified bool       `json:"is_phone_verified"`
	LastLogin       *time.Time `json:"last_login"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
