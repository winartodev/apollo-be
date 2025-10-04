package dto

// RequestResetResponse represents the response for reset password
// swagger:model RequestResetResponse
type RequestResetResponse struct {
	// Redirection link after authentication
	// example: /dashboard
	RedirectionLink string `json:"redirection_link"`

	// OTP information (if applicable)
	Otp *OtpResponse `json:"otp,omitempty"`
}
