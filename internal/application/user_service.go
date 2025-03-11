package application

import (
	"errors"
	"fmt"

	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
	"github.com/Joe5451/go-oauth2-server/internal/application/ports/out"
	"github.com/Joe5451/go-oauth2-server/internal/constants"
	"github.com/Joe5451/go-oauth2-server/internal/domain"
	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type UserService struct {
	userRepo out.UserRepository
}

func NewUserService(userRepo out.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Register(req in.RegisterUserRequest) error {
	password, err := u.hashPassword(req.Password)
	if err != nil {
		return err
	}

	_, err = u.userRepo.CreateUser(domain.User{
		Email:    req.Email,
		Password: password,
		Name:     req.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserService) AuthenticateUser(email, password string) (domain.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return domain.User{}, constants.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domain.User{}, constants.ErrInvalidCredentials
	}

	return user, nil
}

func (u *UserService) GenerateSocialProviderAuthUrl(
	provider socialproviders.SocialProvider,
	state, redirectUri string,
) (string, error) {
	if provider == nil {
		return "", constants.ErrInvalidProvider
	}

	config := provider.NewOauth2Config(redirectUri)
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// This function needs refactoring.
func (u *UserService) LoginWithSocialAccount(
	provider socialproviders.SocialProvider,
	authorizationCode, redirectUri string,
) (domain.User, error) {
	if provider == nil {
		return domain.User{}, constants.ErrInvalidProvider
	}

	socialUser, err := provider.GetUserInformationByAuthorizationCode(authorizationCode, redirectUri)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user information from social provider: %w", err)
	}

	socialAccount, err := u.userRepo.FirstOrCreateSocialAccount(provider.ProviderName(), socialUser.ProviderUserID)

	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to retrieve or create social account (provider: %s, providerUserID: %s): %w",
			provider.ProviderName(),
			socialUser.ProviderUserID,
			err,
		)
	}

	if socialAccount.UserID.Valid {
		user, err := u.userRepo.GetUser(socialAccount.UserID.Int64)
		if err != nil {
			return domain.User{}, constants.ErrUserNotFound
		}
		return user, nil
	}

	user, err := u.userRepo.GetUserByEmail(socialUser.Email)

	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			// A transaction may be needed here to ensure atomicity
			// between user creation and social account update.
			user, err = u.userRepo.Create(domain.User{
				Email:    socialUser.Email,
				Username: socialUser.Username,
			})

			if err != nil {
				return domain.User{}, err
			}
		} else {
			return domain.User{}, err
		}
	}

	socialAccount.UserID.Int64 = user.ID
	socialAccount.UserID.Valid = true
	err = u.userRepo.UpdateSocialAccount(socialAccount)

	return user, nil
}

func (u *UserService) GetUser(userID int64) (domain.User, error) {
	user, err := u.userRepo.GetUser(userID)
	if err != nil {
		return domain.User{}, constants.ErrUserNotFound
	}
	return user, nil
}

func (u *UserService) UpdateUser(userID int64, user domain.User) error {
	err := u.userRepo.UpdateUser(userID, user)
	return err
}
