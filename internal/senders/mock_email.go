package senders

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// MockEmailSender simulates email sending for demo
type MockEmailSender struct {
	log *logrus.Logger
}

func NewMockEmailSender(log *logrus.Logger) *MockEmailSender {
	return &MockEmailSender{log: log}
}

func (s *MockEmailSender) Send(ctx context.Context, payload NotificationPayload) error {
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	// Log as if email was sent
	s.log.WithFields(logrus.Fields{
		"type":    "EMAIL",
		"to":      payload.To,
		"subject": payload.Subject,
	}).Info("[MOCK] Email sent successfully")

	// In real implementation:
	// return smtp.SendEmail(payload.To, payload.Subject, payload.Body)

	return nil
}

func (s *MockEmailSender) GetType() string {
	return "email"
}
