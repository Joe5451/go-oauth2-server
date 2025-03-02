package usecases

import (
	"errors"
	"fmt"

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

// This function needs refactoring.
func (u *UserUsecase) LoginWithSocialAccount(
	provider socialproviders.SocialProvider,
	authorizationCode, redirectUri string,
) (domains.User, error) {
	if provider == nil {
		return domains.User{}, ErrInvalidProvider
	}

	socialUser, err := provider.GetUserInformationByAuthorizationCode(authorizationCode, redirectUri)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to get user information from social provider: %w", err)
	}

	socialAccount, err := u.userRepo.FirstOrCreateSocialAccount(provider.ProviderName(), socialUser.ProviderUserID)

	if err != nil {
		return domains.User{}, fmt.Errorf(
			"failed to retrieve or create social account (provider: %s, providerUserID: %s): %w",
			provider.ProviderName(),
			socialUser.ProviderUserID,
			err,
		)
	}

	if socialAccount.UserID.Valid {
		user, err := u.userRepo.GetUser(socialAccount.UserID.Int64)
		if err != nil {
			return domains.User{}, ErrUserNotFound
		}
		return user, nil
	}

	user, err := u.userRepo.GetUserByEmail(socialUser.Email)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// A transaction may be needed here to ensure atomicity
			// between user creation and social account update.
			user, err = u.userRepo.Create(domains.User{
				Email:    socialUser.Email,
				Username: socialUser.Username,
			})

			if err != nil {
				return domains.User{}, err
			}
		} else {
			return domains.User{}, err
		}
	}

	socialAccount.UserID.Int64 = user.ID
	socialAccount.UserID.Valid = true
	err = u.userRepo.UpdateSocialAccount(socialAccount)

	return user, nil
}

func (u *UserUsecase) GetUser(userID int64) (domains.User, error) {
	user, err := u.userRepo.GetUser(userID)
	if err != nil {
		return domains.User{}, ErrUserNotFound
	}
	return user, nil
}
