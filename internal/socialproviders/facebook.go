package socialproviders

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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
		return SocialProviderUser{}, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/v13.0/me?fields=id,name,email,picture")
	if err != nil {
		return SocialProviderUser{}, fmt.Errorf("failed to fetch user info: %w", err)
	}

	defer resp.Body.Close()

	var user FacebookUser

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SocialProviderUser{}, fmt.Errorf("failed to read response body: %w", err)
	}

	json.Unmarshal(body, &user)

	return SocialProviderUser{
		ProviderUserID: user.ID,
		Email:          user.Email,
		Name:           user.Name,
		Avatar:         user.Picture.Data.URL,
	}, nil
}
