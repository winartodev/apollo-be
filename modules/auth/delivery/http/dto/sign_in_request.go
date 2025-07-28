package dto

import "github.com/winartodev/apollo-be/modules/auth/usecase/dto"

// SignInRequest represents user sign-in request
// swagger:model SignInRequest
type SignInRequest struct {
	// Username or email address (required)
	// required: true
	// min length: 3
	// max length: 50
	// example: john.doe@example.com
	Username string `json:"username" validate:"required,min=3,max=30"`

	// Password (required)
	// required: true
	// min length: 6
	// max length: 100
	// example: secretPassword123
	Password string `json:"password" validate:"required,min=6,max=30"`
}

func (r SignInRequest) ToUseCaseData() dto.SignInDto {
	return dto.SignInDto{
		Username: r.Username,
		Password: r.Password,
	}
}
