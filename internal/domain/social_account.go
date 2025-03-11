package domain

import (
	"time"
)

type SocialAccount struct {
	ID             int64     `json:"-"`
	UserID         *int64    `json:"-"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"-"`
	Email          *string   `json:"email"`
	Name           *string   `json:"name"`
	Avatar         *string   `json:"avatar"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
