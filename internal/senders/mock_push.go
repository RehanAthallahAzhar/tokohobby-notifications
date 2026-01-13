package senders

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// simulates push notification sending for demo
type MockPushSender struct {
	log *logrus.Logger
}

func NewMockPushSender(log *logrus.Logger) *MockPushSender {
	return &MockPushSender{log: log}
}

func (s *MockPushSender) Send(ctx context.Context, payload NotificationPayload) error {
	// Simulate push notification delay
	time.Sleep(50 * time.Millisecond)

	// Log as if push was sent
	s.log.WithFields(logrus.Fields{
		"type":  "PUSH",
		"to":    payload.To,
		"title": payload.Subject,
		"body":  payload.Body,
	}).Info("[MOCK] Push notification sent successfully")

	// In real implementation:
	// return fcm.SendPush(payload.To, payload.Subject, payload.Body)

	return nil
}

func (s *MockPushSender) GetType() string {
	return "push"
}
