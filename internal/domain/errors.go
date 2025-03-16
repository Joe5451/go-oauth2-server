package domain

import "errors"

var (
	ErrInvalidProvider              = errors.New("invalid provider")
	ErrUserNotFound                 = errors.New("user not found")
	ErrInvalidCredentials           = errors.New("invalid credentials")
	ErrSocialAccountNotFound        = errors.New("social account not found")
	ErrDuplicateEmail               = errors.New("duplicate email found")
	ErrSocialUserFetch              = errors.New("failed to fetch user information from social provider")
	ErrInvalidLinkToken             = errors.New("invalid link token")
	ErrMismatchedLinkedUser         = errors.New("mismatched linked user")
	ErrSocialAccountAlreadyLinked   = errors.New("the social account has already been linked to a user")
	ErrSocialAccountAlreadyUnlinked = errors.New("social account is not linked or has already been unlinked")
)
