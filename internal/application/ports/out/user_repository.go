package out

import (
	"github.com/Joe5451/go-oauth2-server/internal/domain"
)

type UserRepository interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUser(userID int64) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	UpdateOrCreateSocialAccount(socialAccount domain.SocialAccount) (domain.SocialAccount, error)
	GetSocialAccountByProviderUserID(providerUserID string) (domain.SocialAccount, error)
	UpdateSocialAccount(account domain.SocialAccount) error
	UpdateUser(usreID int64, user domain.User) error
}
