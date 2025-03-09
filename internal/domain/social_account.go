package domain

import (
	"time"
)

type SocialAccount struct {
	ID             int64     `json:"id"`
	UserID         *int64    `json:"user_id"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	Email          *string   `json:"email"`
	Name           *string   `json:"name"`
	Avatar         *string   `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
