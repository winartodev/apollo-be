package dto

// RequestResetRequest represents the request for reset password
// swagger:model RequestResetRequest
type RequestResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}
