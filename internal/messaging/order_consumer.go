package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	messaging "github.com/RehanAthallahAzhar/tokohobby-messaging-go"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/services"
	"github.com/sirupsen/logrus"
)

type OrderEventConsumer struct {
	rmq          *messaging.RabbitMQ
	notifService *services.NotificationService
	log          *logrus.Logger
}

func NewOrderEventConsumer(rmq *messaging.RabbitMQ, notifService *services.NotificationService, log *logrus.Logger) *OrderEventConsumer {
	return &OrderEventConsumer{
		rmq:          rmq,
		notifService: notifService,
		log:          log,
	}
}

func (c *OrderEventConsumer) Start(ctx context.Context) error {
	c.log.Info("Starting Order Event Consumer...")

	// Message handler
	handler := func(ctx context.Context, body []byte) error {
		// Try to determine event type from message
		var eventType struct {
			Type string `json:"type"`
		}

		if err := json.Unmarshal(body, &eventType); err != nil {
			c.log.WithError(err).Error("Failed to parse event type")
			return err
		}

		// Route to appropriate handler based on type
		switch eventType.Type {
		case "order.created":
			return c.handleOrderCreated(ctx, body)
		case "order.status.changed":
			return c.handleOrderStatusChanged(ctx, body)
		case "order.shipped":
			return c.handleOrderShipped(ctx, body)
		default:
			c.log.Warnf("Unknown event type: %s", eventType.Type)
			return nil
		}
	}

	// Create consumer
	opts := messaging.ConsumerOptions{
		QueueName:   "notifications.order.events",
		WorkerCount: 5,
		AutoAck:     false,
	}

	consumer := messaging.NewConsumer(c.rmq, opts, handler)

	// Declare queue
	if err := consumer.DeclareQueue(true, false); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind to order exchange with all order routing keys
	if err := consumer.BindQueue("order.events", "order.#"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	c.log.Info("Order event consumer configured, starting to consume...")

	// Start consuming
	return consumer.Start(ctx)
}

func (c *OrderEventConsumer) handleOrderCreated(ctx context.Context, body []byte) error {
	var event OrderCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal OrderCreatedEvent: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"order_id": event.OrderID,
		"user_id":  event.UserID,
		"amount":   event.TotalAmount,
	}).Info("Processing OrderCreatedEvent")

	// Create notification
	return c.notifService.CreateAndSendNotification(ctx, &services.CreateNotificationRequest{
		UserID:   event.UserID,
		Type:     "order",
		Category: "created",
		Title:    "Pesanan Dikonfirmasi",
		Message:  fmt.Sprintf("Pesanan #%s telah dikonfirmasi dengan total Rp %.0f", event.OrderID, event.TotalAmount),
		Channels: []string{"email", "push", "in_app"},
		Metadata: map[string]interface{}{
			"order_id":     event.OrderID,
			"total_amount": event.TotalAmount,
			"item_count":   event.ItemCount,
		},
	})
}

func (c *OrderEventConsumer) handleOrderStatusChanged(ctx context.Context, body []byte) error {
	var event OrderStatusChangedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal OrderStatusChangedEvent: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"order_id": event.OrderID,
		"user_id":  event.UserID,
		"status":   event.Status,
	}).Info("Processing OrderStatusChangedEvent")

	// Map status to notification message
	title, message := c.getStatusMessage(event.Status, event.OrderID)

	return c.notifService.CreateAndSendNotification(ctx, &services.CreateNotificationRequest{
		UserID:   event.UserID,
		Type:     "order",
		Category: "status_changed",
		Title:    title,
		Message:  message,
		Channels: []string{"email", "push", "in_app"},
		Metadata: map[string]interface{}{
			"order_id": event.OrderID,
			"status":   event.Status,
		},
	})
}

func (c *OrderEventConsumer) handleOrderShipped(ctx context.Context, body []byte) error {
	var event OrderShippedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal OrderShippedEvent: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"order_id":        event.OrderID,
		"user_id":         event.UserID,
		"tracking_number": event.TrackingNumber,
	}).Info("Processing OrderShippedEvent")

	return c.notifService.CreateAndSendNotification(ctx, &services.CreateNotificationRequest{
		UserID:   event.UserID,
		Type:     "order",
		Category: "shipped",
		Title:    "Pesanan Dikirim",
		Message:  fmt.Sprintf("Pesanan #%s telah dikirim via %s. Nomor resi: %s", event.OrderID, event.Courier, event.TrackingNumber),
		Channels: []string{"email", "push", "in_app"},
		Metadata: map[string]interface{}{
			"order_id":          event.OrderID,
			"tracking_number":   event.TrackingNumber,
			"courier":           event.Courier,
			"estimated_arrival": event.EstimatedArrival,
		},
	})
}

func (c *OrderEventConsumer) getStatusMessage(status, orderID string) (string, string) {
	switch status {
	case "pending":
		return "Menunggu Pembayaran", fmt.Sprintf("Pesanan #%s menunggu pembayaran", orderID)
	case "paid":
		return "Pembayaran Diterima", fmt.Sprintf("Pembayaran pesanan #%s telah diterima", orderID)
	case "processing":
		return "Pesanan Diproses", fmt.Sprintf("Pesanan #%s sedang diproses", orderID)
	case "shipped":
		return "Pesanan Dikirim", fmt.Sprintf("Pesanan #%s sedang dalam pengiriman", orderID)
	case "delivered":
		return "Pesanan Sampai", fmt.Sprintf("Pesanan #%s telah sampai", orderID)
	case "cancelled":
		return "Pesanan Dibatalkan", fmt.Sprintf("Pesanan #%s telah dibatalkan", orderID)
	default:
		return "Status Pesanan Diubah", fmt.Sprintf("Status pesanan #%s: %s", orderID, status)
	}
}
