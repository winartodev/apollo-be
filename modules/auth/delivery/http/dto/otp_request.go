package dto

// OtpRequest represents OTP validation request
// swagger:model OtpRequest
type OtpRequest struct {
	// OTP number (required)
	// required: true
	// minimum: 0
	// maximum: 6
	// example: 123456
	OTPNumber int64 `json:"otp" validate:"required"`
}

// OtpRequest represents OTP resend request
// swagger:model OtpResendRequest
type OtpResendRequest struct {
	Username string `json:"username" validate"required"`
}
