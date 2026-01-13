package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID       uuid.UUID              `json:"id"`
	UserID   uuid.UUID              `json:"user_id"`
	Type     string                 `json:"type"`
	Category string                 `json:"category"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
	Channels []string               `json:"channels"`
	Status   string                 `json:"status"`

	IsRead bool       `json:"is_read"`
	ReadAt *time.Time `json:"read_at,omitempty"`

	EmailSentAt *time.Time `json:"email_sent_at,omitempty"`
	PushSentAt  *time.Time `json:"push_sent_at,omitempty"`

	RetryCount int    `json:"retry_count"`
	LastError  string `json:"last_error,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type NotificationPreference struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`

	EmailEnabled bool `json:"email_enabled"`
	PushEnabled  bool `json:"push_enabled"`
	InAppEnabled bool `json:"in_app_enabled"`

	OrderNotifications   map[string]bool `json:"order_notifications"`
	AccountNotifications map[string]bool `json:"account_notifications"`
	ProductNotifications map[string]bool `json:"product_notifications"`

	QuietHoursEnabled bool       `json:"quiet_hours_enabled"`
	QuietHoursStart   *time.Time `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd     *time.Time `json:"quiet_hours_end,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
