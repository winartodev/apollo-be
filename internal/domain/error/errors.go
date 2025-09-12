package error

import (
	"errors"
	"net/http"
)

var (
	ErrFailedCreateUser          = errors.New("failed to create user")
	ErrFailedUpdateRefreshToken  = errors.New("failed to update refresh token")
	ErrFailedGetUserData         = errors.New("failed to get user data")
	ErrUserNotFound              = errors.New("user_not_found")
	ErrUserAlreadyExists         = errors.New("user_already_exists")
	ErrUsernameAlreadyExists     = errors.New("username_already_exists")
	ErrEmailAlreadyExists        = errors.New("email_already_exists")
	ErrInvalidUsernameOrPassword = errors.New("invalid_credentials")
	ErrOtpTooManyRequest         = errors.New("otp_too_many_request")
	ErrInvalidOTPNumber          = errors.New("otp_invalid_number")
	ErrInvalidEmail              = errors.New("invalid_email")
)

// ErrorCodeMapping pairs an error with an HTTP status code
type ErrorCodeMapping struct {
	Target error
	Status int
}

// errorMappings contains all known error to status mappings
var errorMappings = []ErrorCodeMapping{
	{ErrUserNotFound, http.StatusNotFound},
	{ErrUserAlreadyExists, http.StatusConflict},
	{ErrUsernameAlreadyExists, http.StatusConflict},
	{ErrEmailAlreadyExists, http.StatusConflict},
	{ErrInvalidUsernameOrPassword, http.StatusUnauthorized},
	{ErrOtpTooManyRequest, http.StatusTooManyRequests},
	{ErrInvalidOTPNumber, http.StatusBadRequest},

	// For generic internal failures, map to 500
	{ErrFailedCreateUser, http.StatusInternalServerError},
	{ErrFailedUpdateRefreshToken, http.StatusInternalServerError},
	{ErrFailedGetUserData, http.StatusInternalServerError},
}

// GetHTTPStatusFromError returns the HTTP status code for a given error.
// Returns 500 if no mapping found.
func GetHTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	for _, mapping := range errorMappings {
		if errors.Is(err, mapping.Target) {
			return mapping.Status
		}
	}

	return http.StatusInternalServerError
}
