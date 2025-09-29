package domain

import "time"

type SharedUser struct {
	ID              int64
	Username        string
	Email           string
	PhoneNumber     string
	FirstName       string
	LastName        string
	IsActive        bool
	IsEmailVerified bool
	IsPhoneVerified bool
	LastLogin       *time.Time
}
