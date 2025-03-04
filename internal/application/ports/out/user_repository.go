package out

import (
	"github.com/Joe5451/go-oauth2-server/internal/domain"
)

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	GetUser(userID int64) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	FirstOrCreateSocialAccount(provider, providerUserID string) (domain.SocialAccount, error)
	CreateSocialAccount(account domain.SocialAccount) domain.SocialAccount
	UpdateSocialAccount(account domain.SocialAccount) error
	UpdateUser(usreID int64, user domain.User) error
}
