package domains

import "time"

type SocialAccount struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           User      `gorm:"references:ID"`
}
