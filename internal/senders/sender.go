package senders

import "context"

// NotificationPayload represents data to send
type NotificationPayload struct {
	To      string
	Subject string
	Body    string
	Data    map[string]interface{}
}

// Sender interface for notification senders
type Sender interface {
	Send(ctx context.Context, payload NotificationPayload) error
	GetType() string
}
