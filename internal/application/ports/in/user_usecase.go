package in

import (
	"github.com/Joe5451/go-oauth2-server/internal/domain"
	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
	"github.com/golang-jwt/jwt"
)

type RegisterUserRequest struct {
	Email    string
	Password string
	Name     string
}

type AuthSocialUserStatus string

const (
	AuthSuccess      AuthSocialUserStatus = "success"       // Successful authentication
	AuthLinkRequired AuthSocialUserStatus = "link_required" // Link required for the account
)

type LinkTokenClaims struct {
	UserID               int64                  `json:"user_id"`
	SocialAccountID      int64                  `json:"social_account_id"`
	LinkedSocialAccounts []domain.SocialAccount `json:"linked_social_accounts"`
	jwt.StandardClaims
}

type AuthSocialUserResult struct {
	Status    AuthSocialUserStatus // Authentication result status
	User      domain.User          // User details after successful authentication
	LinkToken string               // Token for linking social accounts
}

type UserUsecase interface {
	Register(req RegisterUserRequest) error
	AuthenticateUser(email, password string) (domain.User, error)
	SocialAuthUrl(provider socialproviders.SocialProvider, state, redirectUri string) (string, error)
	AuthenticateSocialUser(provider socialproviders.SocialProvider, authorizationCode, redirectUri string) (AuthSocialUserResult, error)
	LinkUserWithSocialAccount(provider socialproviders.SocialProvider, authCode string, linkToken string, redirectUri string) (domain.User, error)
	ValidateLinkToken(linkToken string) (LinkTokenClaims, error)
	GetUser(userID int64) (domain.User, error)
	UpdateUser(userID int64, user domain.User) error
}
