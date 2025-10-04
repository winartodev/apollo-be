package dto

import "github.com/winartodev/apollo-be/modules/auth/usecase/dto"

// ResetPasswordRequest represents the request payload for resetting a user's password.
// It includes the new password, its confirmation, and the user's email.
//
// swagger:model ResetPasswordRequest
type ResetPasswordRequest struct {
	// Email of the user requesting the password reset.
	// required: true
	// format: email
	Email string `json:"email" validate:"required,email"`

	// Password is the new password the user wants to set.
	// required: true
	// min length: 3
	Password string `json:"password" validate:"required,min=3"`

	// PasswordConfirmation must match the password field.
	// required: true
	// min length: 3
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=3"`
}

func (e *ResetPasswordRequest) ToUseCaseData() dto.ResetPasswordDto {
	return dto.ResetPasswordDto{
		Email:                e.Email,
		Password:             e.Password,
		PasswordConfirmation: e.PasswordConfirmation,
	}
}
