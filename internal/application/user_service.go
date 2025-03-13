package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
	"github.com/Joe5451/go-oauth2-server/internal/application/ports/out"
	"github.com/Joe5451/go-oauth2-server/internal/config"
	"github.com/Joe5451/go-oauth2-server/internal/domain"
	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
	"github.com/golang-jwt/jwt"
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (u *UserService) AuthenticateUser(email, password string) (domain.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, domain.ErrInvalidCredentials
		}
		return domain.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domain.User{}, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (u *UserService) SocialAuthUrl(provider socialproviders.SocialProvider, state, redirectUri string) (string, error) {
	if provider == nil {
		return "", domain.ErrInvalidProvider
	}

	config := provider.NewOauth2Config(redirectUri)
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (u *UserService) AuthenticateSocialUser(provider socialproviders.SocialProvider, authorizationCode, redirectUri string) (in.AuthSocialUserResult, error) {
	if provider == nil {
		return in.AuthSocialUserResult{}, domain.ErrInvalidProvider
	}

	socialUser, err := provider.GetUserInformationByAuthorizationCode(authorizationCode, redirectUri)
	if err != nil {
		return in.AuthSocialUserResult{}, err
	}

	socialAccount, err := u.updateOrCreateSocialAccount(provider, socialUser)
	if err != nil {
		return in.AuthSocialUserResult{}, err
	}

	if socialAccount.UserID != nil {
		return u.authenticateLinkedUser(*socialAccount.UserID)
	}

	return u.handleUnlinkedSocialAccount(&socialAccount)
}

func (u *UserService) updateOrCreateSocialAccount(
	provider socialproviders.SocialProvider,
	socialUser socialproviders.SocialProviderUser,
) (domain.SocialAccount, error) {
	socialAccount, err := u.userRepo.UpdateOrCreateSocialAccount(domain.SocialAccount{
		Provider:       provider.ProviderName(),
		ProviderUserID: socialUser.ProviderUserID,
		Email:          &socialUser.Email,
		Name:           &socialUser.Name,
		Avatar:         &socialUser.Avatar,
	})
	if err != nil {
		return domain.SocialAccount{}, fmt.Errorf("failed to update or create social account: %w", err)
	}
	return socialAccount, nil
}

func (u *UserService) authenticateLinkedUser(userID int64) (in.AuthSocialUserResult, error) {
	user, err := u.userRepo.GetUser(userID)
	if err != nil {
		return in.AuthSocialUserResult{}, fmt.Errorf("unable to retrieve user associated with social account: %w", err)
	}
	return in.AuthSocialUserResult{Status: in.AuthSuccess, User: user}, nil
}

func (u *UserService) handleUnlinkedSocialAccount(socialAccount *domain.SocialAccount) (in.AuthSocialUserResult, error) {
	user, err := u.userRepo.GetUserByEmail(*socialAccount.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			user, err := u.createUserBySocialAccount(*socialAccount)
			if err != nil {
				return in.AuthSocialUserResult{}, fmt.Errorf("failed to create new user or link social account: %w", err)
			}
			return in.AuthSocialUserResult{Status: in.AuthSuccess, User: user}, nil
		}
		return in.AuthSocialUserResult{}, err
	}

	linkToken, err := u.generateLinkToken(user, socialAccount.ID)
	if err != nil {
		return in.AuthSocialUserResult{}, fmt.Errorf("failed to generate social account link token: %w", err)
	}

	return in.AuthSocialUserResult{
		Status:    in.AuthLinkRequired,
		User:      user,
		LinkToken: linkToken,
	}, nil
}

func (u *UserService) createUserBySocialAccount(account domain.SocialAccount) (domain.User, error) {
	user, err := u.userRepo.CreateUser(domain.User{
		Email:  *account.Email,
		Name:   *account.Name,
		Avatar: account.Avatar,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user for social account: %w", err)
	}

	err = u.userRepo.UpdateSocialAccountUserID(account.ID, user.ID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserService) generateLinkToken(user domain.User, socialAccountID int64) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := in.LinkTokenClaims{
		UserID:               user.ID,
		SocialAccountID:      socialAccountID,
		LinkedSocialAccounts: user.SocialAccounts,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.AppConfig.JwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserService) LinkUserWithSocialAccount(
	provider socialproviders.SocialProvider,
	authCode string,
	linkToken string,
	redirectUri string,
) (domain.User, error) {
	if provider == nil {
		return domain.User{}, domain.ErrInvalidProvider
	}

	claims, err := u.ValidateLinkToken(linkToken)
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %v", domain.ErrInvalidLinkToken, err.Error())
	}
	userID, socialAccountID := claims.UserID, claims.SocialAccountID

	socialUser, err := provider.GetUserInformationByAuthorizationCode(authCode, redirectUri)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user information from social provider: %w", err)
	}

	socialAccount, err := u.userRepo.GetSocialAccountByProviderUserID(socialUser.ProviderUserID)
	if err != nil {
		return domain.User{}, err
	}

	if *socialAccount.UserID != userID {
		return domain.User{}, domain.ErrMismatchedLinkedUser
	}

	err = u.userRepo.UpdateSocialAccountUserID(socialAccountID, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to link account: %w", err)
	}

	return u.GetUser(userID)
}

func (u *UserService) ValidateLinkToken(linkToken string) (in.LinkTokenClaims, error) {
	secretKey := []byte(config.AppConfig.JwtSecret)

	var claims in.LinkTokenClaims

	_, err := jwt.ParseWithClaims(linkToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return secretKey, nil
	})

	if err != nil {
		return in.LinkTokenClaims{}, fmt.Errorf("invalid token: %v", err)
	}

	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return in.LinkTokenClaims{}, fmt.Errorf("token has expired")
	}

	return claims, nil
}

func (u *UserService) GetUser(userID int64) (domain.User, error) {
	user, err := u.userRepo.GetUser(userID)
	if err != nil {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (u *UserService) UpdateUser(userID int64, user domain.User) error {
	err := u.userRepo.UpdateUser(userID, user)
	return err
}
