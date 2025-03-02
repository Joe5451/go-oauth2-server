package domains

import (
	"database/sql"
	"time"

	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
)

type User struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	PhoneNumber string    `json:"phone_number"`
	Username    string    `json:"username"`
	Gender      string    `json:"gender"`
	Avatar      string    `json:"avatar"`
	Birthday    time.Time `json:"birthday"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SocialAccount struct {
	ID             int64         `json:"id"`
	UserID         sql.NullInt64 `json:"user_id"`
	Provider       string        `json:"provider"`
	ProviderUserID string        `json:"provider_user_id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	User           User          `gorm:"references:ID"`
}

type UserUsecase interface {
	Register(user User) error
	LoginWithEmail(email, password string) (User, error)
	GenerateSocialProviderAuthUrl(provider socialproviders.SocialProvider, state, redirectUri string) (string, error)
	LoginWithSocialAccount(provider socialproviders.SocialProvider, authorizationCode, redirectUri string) (User, error)
	GetUser(userID int64) (User, error)
	UpdateUser(userID int64, user User) error
}

type UserRepository interface {
	Create(user User) (User, error)
	GetUser(userID int64) (User, error)
	GetUserByEmail(email string) (User, error)
	FirstOrCreateSocialAccount(provider, providerUserID string) (SocialAccount, error)
	CreateSocialAccount(account SocialAccount) SocialAccount
	UpdateSocialAccount(account SocialAccount) error
	UpdateUser(usreID int64, user User) error
}
