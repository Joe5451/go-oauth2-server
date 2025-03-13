package socialproviders

import "errors"

var (
	ErrOAuth2RetrieveError = errors.New("OAuth2 retrieve error")
	ErrInvalidProvider     = errors.New("invalid social provider")
)
