package messaging

import "time"

// UserRegisteredEvent represents user registration event
type UserRegisteredEvent struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	RegisteredAt time.Time `json:"registered_at"`
}
