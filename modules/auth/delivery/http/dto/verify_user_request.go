package dto

type VerifyUserRequest struct {
	Username string `json:"username" validate:"required"`
}
