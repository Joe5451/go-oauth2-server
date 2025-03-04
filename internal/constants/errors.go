package constants

import "errors"

var (
	ErrInvalidProvider    = errors.New("invalid provider")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
