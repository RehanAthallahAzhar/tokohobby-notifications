package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	messaging "github.com/RehanAthallahAzhar/tokohobby-messaging/rabbitmq"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/services"
	"github.com/sirupsen/logrus"
)

type BlogEventConsumer struct {
	rmq          *messaging.RabbitMQ
	notifService *services.NotificationService
	log          *logrus.Logger
}

func NewBlogEventConsumer(
	rmq *messaging.RabbitMQ,
	notifService *services.NotificationService,
	log *logrus.Logger,
) *BlogEventConsumer {
	return &BlogEventConsumer{
		rmq:          rmq,
		notifService: notifService,
		log:          log,
	}
}

func (c *BlogEventConsumer) Start(ctx context.Context) error {
	handler := func(ctx context.Context, body []byte) error {
		var eventType struct {
			Type string `json:"type"`
		}

		if err := json.Unmarshal(body, &eventType); err != nil {
			return fmt.Errorf("failed to unmarshal event type: %w", err)
		}

		switch eventType.Type {
		case "blog.published":
			return c.handleBlogPublished(ctx, body)
		case "comment.added":
			return c.handleCommentAdded(ctx, body)
		default:
			c.log.Warnf("Unknown event type: %s", eventType.Type)
			return nil
		}
	}

	opts := messaging.ConsumerOptions{
		QueueName:   "blog.notifications",
		WorkerCount: 5,
		AutoAck:     false,
	}

	consumer := messaging.NewConsumer(c.rmq, opts, handler)

	// Declare queue
	if err := consumer.DeclareQueue(true, false); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind to blog.events with wildcard
	if err := consumer.BindQueue("blog.events", "blog.#"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	c.log.Info("Blog event consumer configured, starting to consume...")

	return consumer.Start(ctx)
}

func (c *BlogEventConsumer) handleBlogPublished(ctx context.Context, body []byte) error {
	var event BlogPublishedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal BlogPublishedEvent: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"blog_id": event.BlogID,
		"author":  event.AuthorName,
		"title":   event.Title,
	}).Info("Processing blog published event")

	// TODO: Get followers from database and send notifications
	// For now, just log the event
	message := fmt.Sprintf("%s published a new blog: %s", event.AuthorName, event.Title)
	c.log.Infof("Would notify followers: %s - %s", message, event.Excerpt)

	// When follower system is implemented, use this:
	// return c.notifService.SendBlogPublishedNotification(ctx, &services.BlogNotificationRequest{
	//     FollowerIDs: followerIDs,
	//     Title:       "New Blog Published",
	//     Message:     message,
	//     BlogURL:     fmt.Sprintf("/blogs/%s", event.BlogID),
	// })

	return nil
}

func (c *BlogEventConsumer) handleCommentAdded(ctx context.Context, body []byte) error {
	var event CommentAddedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CommentAddedEvent: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"comment_id": event.CommentID,
		"blog_title": event.BlogTitle,
		"commenter":  event.Commenter,
		"blog_owner": event.BlogOwnerID,
	}).Info("Processing comment added event")

	// Notify blog owner about new comment
	return c.notifService.CreateAndSendNotification(ctx, &services.CreateNotificationRequest{
		UserID:   event.BlogOwnerID,
		Type:     "blog",
		Category: "comment",
		Title:    "New Comment",
		Message:  fmt.Sprintf("%s commented on your blog '%s': %s", event.Commenter, event.BlogTitle, event.Comment),
		Channels: []string{"email", "push", "in_app"},
		Metadata: map[string]interface{}{
			"blog_id":    event.BlogID,
			"comment_id": event.CommentID,
			"commenter":  event.Commenter,
		},
	})
}
