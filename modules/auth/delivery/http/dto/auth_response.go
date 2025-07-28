package dto

// AuthResponse represents authentication response
// swagger:model AuthResponse
type AuthResponse struct {
	// JWT access token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`

	// Refresh token for obtaining new access tokens
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token"`

	// Redirection link after authentication
	// example: /dashboard
	RedirectionLink string `json:"redirection_link"`

	// OTP information (if applicable)
	Otp *OtpRefreshResponse `json:"otp,omitempty"`
}
