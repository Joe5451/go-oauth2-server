package domains

import "time"

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
