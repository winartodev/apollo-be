package dto

import "github.com/winartodev/apollo-be/modules/auth/usecase/dto"

// SignUpRequest represents user registration request
// swagger:model SignUpRequest
type SignUpRequest struct {
	// Username (required)
	// required: true
	// min length: 3
	// max length: 30
	// pattern: ^[a-zA-Z0-9_]+$
	// example: JohnDoe
	Username string `json:"username"  validate:"required,min=3,max=30"`

	// Password (required)
	// required: true
	// min length: 8
	// max length: 100
	// example: SecurePass123!
	Password string `json:"password"`

	// Email address (required)
	// required: true
	// format: email
	// example: john.doe@example.com
	Email string `json:"email"`

	// Phone number (optional)
	// pattern: ^\+?[1-9]\d{1,14}$
	// example: +1234567890
	PhoneNumber string `json:"phone_number"`
}

func (r SignUpRequest) ToUseCaseData() dto.SignUpDto {
	return dto.SignUpDto{
		Username:    r.Username,
		Password:    r.Password,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
	}
}
