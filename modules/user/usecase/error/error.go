package error

import "errors"

var (
	UsernameAlreadyExists = errors.New("username_already_exists")
	EmailAlreadyExists    = errors.New("email_already_exists")
)
