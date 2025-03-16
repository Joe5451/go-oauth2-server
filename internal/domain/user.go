package domain

import (
	"time"
)

type User struct {
	ID             int64
	Email          string          `json:"email"`
	Password       string          `json:"-"`
	Name           string          `json:"name"`
	Avatar         *string         `json:"avatar"`
	SocialAccounts []SocialAccount `json:"social_accounts"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
