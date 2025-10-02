package dto

// OtpResponse represents the response for OTP refresh operations
// swagger:model OtpResponse
type OtpResponse struct {
	// Number of retry attempts remaining
	// example: 3
	RetryAttemptsLeft int64 `json:"retry_attempts_left"`

	// Time in seconds until the OTP expires
	// example: 300
	ExpiresIn int64 `json:"expires_in"`

	// Time in seconds to wait before retrying (optional)
	// example: 60
	RetryAfterIn int64 `json:"retry_after_seconds,omitempty"`

	// Indicates if the OTP is valid
	// example: true
	IsValid bool `json:"is_valid"`
}

// OtpValidationResponse represents the response for OTP validation
// swagger:model OtpValidationResponse
type OtpValidationResponse struct {
	// Indicates if the OTP validation was successful
	// example: true
	IsValid bool `json:"is_valid"`

	// Message describing the validation result
	// example: OTP validated successfully
	Message string `json:"message"`

	// Redirection link after successful validation
	// example: https://example.com/dashboard
	RedirectionLink string `json:"redirection_link"`
}
