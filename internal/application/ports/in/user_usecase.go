package in

import (
	"github.com/Joe5451/go-oauth2-server/internal/domain"
	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
)

type RegisterUserRequest struct {
	Email    string
	Password string
	Username string
}

type UserUsecase interface {
	Register(req RegisterUserRequest) error
	LoginWithEmail(email, password string) (domain.User, error)
	GenerateSocialProviderAuthUrl(provider socialproviders.SocialProvider, state, redirectUri string) (string, error)
	LoginWithSocialAccount(provider socialproviders.SocialProvider, authorizationCode, redirectUri string) (domain.User, error)
	GetUser(userID int64) (domain.User, error)
	UpdateUser(userID int64, user domain.User) error
}
