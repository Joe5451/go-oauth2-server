package socialproviders

import (
	"golang.org/x/oauth2"
)

type SocialProvider interface {
	ProviderName() string
	NewOauth2Config(redirectUri string) *oauth2.Config
	GetUserInformationByAuthorizationCode(code, redirectUri string) (SocialProviderUser, error)
}

type SocialProviderUser struct {
	ProviderUserID string
	Email          string
	Username       string
	Avatar         string
}
