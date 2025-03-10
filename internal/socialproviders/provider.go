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

func NewSocialProvider(provider string) (SocialProvider, error) {
	switch provider {
	case "google":
		return NewGoogleProvider(), nil
	case "facebook":
		return NewFacebookProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
