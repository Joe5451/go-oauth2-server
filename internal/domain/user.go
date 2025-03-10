package domain

import (
	"time"
)

type User struct {
	ID             int64   `json:"id"`
	Email          string  `json:"email"`
	Password       string  `json:"-"`
	Name           string  `json:"name"`
	Avatar         *string `json:"avatar"`
	SocialAccounts []SocialAccount
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
