package middleware

import "errors"

var (
	errInvalidToken = errors.New("invalid token")

	errorAuthorizationHeaderEmpty   = errors.New("authorization header is empty")
	errorInvalidAuthorizationHeader = errors.New("invalid authorization header format")
	errorEmptyToken                 = errors.New("token is empty")
)
