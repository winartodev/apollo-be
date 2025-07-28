package dto

// OtpRequest represents OTP validation request
// swagger:model OtpRequest
type OtpRequest struct {
	// OTP number (required)
	// required: true
	// minimum: 6
	// maximum: 6
	// example: 123456
	OTPNumber int64 `json:"otp" validate:"required,min=0,max=6"`
}
