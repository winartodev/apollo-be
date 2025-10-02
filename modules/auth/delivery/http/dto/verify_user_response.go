package dto

type VerifyUserResponse struct {
	UserExists  bool     `json:"user_exists"`
	Suggestions []string `json:"suggestions,omitempty"`
}
