package dto

import (
	"time"

	"github.com/winartodev/apollo-be/modules/user/delivery/http/dto"
)

type UserDto struct {
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
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

func (u *UserDto) ToResponse() *dto.UserResponse {
	return &dto.UserResponse{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		PhoneNumber:     u.PhoneNumber,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		IsActive:        u.IsActive,
		IsEmailVerified: u.IsEmailVerified,
		IsPhoneVerified: u.IsPhoneVerified,
		LastLogin:       u.LastLogin,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		DeletedAt:       u.DeletedAt,
	}
}
