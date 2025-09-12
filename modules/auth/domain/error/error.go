package error

import "errors"

const (
	ErrFailedCreateUser         = "failed to create user %v"
	ErrFailedUpdateRefreshToken = "failed to update refresh token %v"
	ErrFailedGetUserData        = "failed to get user data %v"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("invalid_credentials")
	ErrUserNotFound              = errors.New("user_not_found")
	ErrOtpToManyRequest          = errors.New("otp_to_many_request")
	ErrInvalidOTPNumber          = errors.New("otp_invalid_number")
	ErrUserAlreadyExists         = errors.New("user_already_exists")
)
