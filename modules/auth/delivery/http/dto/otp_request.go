package dto

// OtpRequest represents the payload to validate an OTP code.
//
// swagger:model OtpRequest
type OtpRequest struct {
	// OTP Number (6 digits)
	// required: true
	// minimum: 100000
	// maximum: 999999
	// example: 123456
	OTPNumber string `json:"otp" validate:"required,min=0,max=999999"`

	// Email user email address
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`

	// Type of OTP request, e.g. signup, reset_password
	// required: true
	// enum: signup,request_reset
	// example: request_reset
	Type string `json:"type" validate:"required,oneof=signup request_reset"`
}

// OtpResendRequest represents OTP resend request
//
// swagger:model OtpResendRequest
type OtpResendRequest struct {
	// Email email address
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`

	// Type of OTP request, e.g. signup, reset_password, login
	// required: true
	// enum: signup,reset_password,login
	// example: reset_password
	Type string `json:"type" validate:"required,oneof=signup request_reset"`
}
