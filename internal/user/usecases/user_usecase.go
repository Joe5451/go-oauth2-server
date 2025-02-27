package usecases

import (
	"github.com/Joe5451/go-oauth2-server/internal/domains"
	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type UserUsecase struct {
	userRepo domains.UserRepository
}

func (u *UserUsecase) Register(user domains.User) error {
	user, err := u.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) LoginWithEmail(email, password string) (domains.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return domains.User{}, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domains.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (u *UserUsecase) GenerateSocialProviderAuthUrl(
	provider socialproviders.SocialProvider,
	state, redirectUri string,
) (string, error) {
	if provider == nil {
		return "", ErrInvalidProvider
	}

	config := provider.NewOauth2Config(redirectUri)
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (u *UserUsecase) GetUserByID(userID int64) (domains.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return domains.User{}, ErrUserNotFound
	}

	return user, nil
}
