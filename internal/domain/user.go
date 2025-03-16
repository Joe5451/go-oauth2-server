package domain

import (
	"time"
)

type User struct {
	ID             int64           `json:"-"`
	Email          string          `json:"email"`
	Password       string          `json:"-"`
	Name           string          `json:"name"`
	Avatar         *string         `json:"avatar"`
	SocialAccounts []SocialAccount `json:"social_accounts"`
	CreatedAt      time.Time       `json:"-"`
	UpdatedAt      time.Time       `json:"-"`
}
