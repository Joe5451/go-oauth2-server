package handlers

import (
	"errors"
	"net/http"

	"github.com/Joe5451/go-oauth2-server/internal/domain"
)

var (
	ErrValidation   = errors.New("Validation error")
	ErrUnauthorized = errors.New("Requires authentication.")
)

var errorMap = map[error]struct {
	httpCode  int
	errorCode string
	message   string
}{
	ErrUnauthorized:              {http.StatusUnauthorized, "UNAUTHORIZED", "Requires authentication."},
	domain.ErrUserNotFound:       {http.StatusNotFound, "USER_NOT_FOUND", "The user does not exist."},
	domain.ErrInvalidCredentials: {http.StatusUnauthorized, "INVALID_CREDENTIALS", "Incorrect email or password."},
	domain.ErrDuplicateEmail:     {http.StatusConflict, "DUPLICATE_EMAIL", "The email is already in use."},
	domain.ErrInvalidProvider:    {http.StatusBadRequest, "INVALID_SOCIAL_PROVIDER", "Invalid social provider."},
	domain.ErrSocialUserFetch:    {http.StatusInternalServerError, "SOCIAL_PROVIDER_ERROR", "Failed to fetch user information from social provider."},
}
