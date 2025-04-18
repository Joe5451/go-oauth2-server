package socialproviders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Joe5451/go-oauth2-server/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookProvider struct {
}

type PictureData struct {
	Height       int    `json:"height"`
	IsSilhouette bool   `json:"is_silhouette"`
	URL          string `json:"url"`
	Width        int    `json:"width"`
}

type Picture struct {
	Data PictureData `json:"data"`
}

type FacebookUser struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Picture Picture `json:"picture"`
}

func NewFacebookProvider() *FacebookProvider {
	return &FacebookProvider{}
}

func (p *FacebookProvider) ProviderName() string {
	return "facebook"
}

func (p *FacebookProvider) NewOauth2Config(redirectUri string) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.AppConfig.FacebookOauth2ClientID,
		ClientSecret: config.AppConfig.FacebookOauth2ClientSecret,
		RedirectURL:  redirectUri,
		Scopes: []string{
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	return conf
}

func (p *FacebookProvider) GetUserInformationByAuthorizationCode(code, redirectUri string) (SocialProviderUser, error) {
	config := p.NewOauth2Config(redirectUri)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		var retrieveError *oauth2.RetrieveError
		if errors.As(err, &retrieveError) {
			return SocialProviderUser{}, fmt.Errorf("%w: %v", ErrOAuth2RetrieveError, retrieveError)
		}
		return SocialProviderUser{}, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/v13.0/me?fields=id,name,email,picture")
	if err != nil {
		return SocialProviderUser{}, fmt.Errorf("failed to fetch user info from Facebook: %w", err)
	}

	defer resp.Body.Close()

	var user FacebookUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return SocialProviderUser{}, fmt.Errorf("failed to decode Facebook user info response: %v", err)
	}

	return SocialProviderUser{
		ProviderUserID: user.ID,
		Email:          user.Email,
		Name:           user.Name,
		Avatar:         user.Picture.Data.URL,
	}, nil
}
