package socialproviders

import (
	"context"
	"errors"
	"fmt"

	"github.com/Joe5451/go-oauth2-server/internal/config"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
	jwt.StandardClaims
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
}

func (p *GoogleProvider) ProviderName() string {
	return "google"
}

func (p *GoogleProvider) NewOauth2Config(redirectUri string) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.AppConfig.GoogleOauth2ClientID,
		ClientSecret: config.AppConfig.GoogleOauth2ClientSecret,
		RedirectURL:  redirectUri,
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

func (p *GoogleProvider) GetUserInformationByAuthorizationCode(code, redirectUri string) (SocialProviderUser, error) {
	config := p.NewOauth2Config(redirectUri)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		var retrieveError *oauth2.RetrieveError
		if errors.As(err, &retrieveError) {
			return SocialProviderUser{}, fmt.Errorf("%w: %v", ErrOAuth2RetrieveError, retrieveError.ErrorCode)
		}
		return SocialProviderUser{}, err
	}

	rawIdToken := token.Extra("id_token").(string)
	idToken, _, err := new(jwt.Parser).ParseUnverified(rawIdToken, &GoogleClaims{})
	if err != nil {
		return SocialProviderUser{}, fmt.Errorf("failed to parse Google id_token: %w", err)
	}

	claims, ok := idToken.Claims.(*GoogleClaims)
	if !ok {
		return SocialProviderUser{}, fmt.Errorf("failed to extract user claims from Google id_token: invalid claims format")
	}

	return SocialProviderUser{
		ProviderUserID: claims.Sub,
		Email:          claims.Email,
		Name:           claims.Name,
		Avatar:         claims.Picture,
	}, nil
}
