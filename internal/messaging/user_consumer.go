package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	messaging "github.com/RehanAthallahAzhar/tokohobby-messaging-go"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/services"
	"github.com/sirupsen/logrus"
)

type UserEventConsumer struct {
	rmq          *messaging.RabbitMQ
	notifService *services.NotificationService
	log          *logrus.Logger
}

func NewUserEventConsumer(rmq *messaging.RabbitMQ, notifService *services.NotificationService, log *logrus.Logger) *UserEventConsumer {
	return &UserEventConsumer{
		rmq:          rmq,
		notifService: notifService,
		log:          log,
	}
}

func (c *UserEventConsumer) Start(ctx context.Context) error {
	c.log.Info("Starting User Event Consumer...")

	// Message handler
	handler := func(ctx context.Context, body []byte) error {
		// Try to determine event type
		var eventType struct {
			Type string `json:"type"`
		}

		if err := json.Unmarshal(body, &eventType); err != nil {
			c.log.WithError(err).Error("Failed to parse event type")
			return err
		}

		// Route based on type
		switch eventType.Type {
		case "user.registered":
			return c.handleUserRegistered(ctx, body)
		default:
			c.log.Warnf("Unknown user event type: %s", eventType.Type)
			return nil
		}
	}

	// Create consumer
	opts := messaging.ConsumerOptions{
		QueueName:   "notifications.user.events",
		WorkerCount: 3,
		AutoAck:     false,
	}

	consumer := messaging.NewConsumer(c.rmq, opts, handler)

	// Declare queue
	if err := consumer.DeclareQueue(true, false); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind to user exchange with all user routing keys
	if err := consumer.BindQueue("user.events", "user.#"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	c.log.Info("User event consumer configured, starting to consume...")

	// Start consuming
	return consumer.Start(ctx)
}

func (c *UserEventConsumer) handleUserRegistered(ctx context.Context, body []byte) error {
	var event UserRegisteredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal UserRegisteredEvent: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"user_id":  event.UserID,
		"username": event.Username,
		"email":    event.Email,
	}).Info("Processing UserRegisteredEvent")

	// Create welcome notification
	return c.notifService.CreateAndSendNotification(ctx, &services.CreateNotificationRequest{
		UserID:   event.UserID,
		Type:     "account",
		Category: "registered",
		Title:    "Selamat Datang di TokoHobby!",
		Message:  fmt.Sprintf("Hai %s, terima kasih telah bergabung dengan TokoHobby. Akun Anda telah aktif dan siap digunakan!", event.Username),
		Channels: []string{"email", "in_app"},
		Metadata: map[string]interface{}{
			"username":  event.Username,
			"full_name": event.FullName,
			"email":     event.Email,
		},
	})
}
