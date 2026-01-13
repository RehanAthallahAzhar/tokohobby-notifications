package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/entities"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/repositories"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/senders"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type NotificationService struct {
	repo        *repositories.NotificationRepository
	emailSender senders.Sender
	pushSender  senders.Sender
	log         *logrus.Logger
}

func NewNotificationService(repo *repositories.NotificationRepository, emailSender, pushSender senders.Sender, log *logrus.Logger) *NotificationService {
	return &NotificationService{
		repo:        repo,
		emailSender: emailSender,
		pushSender:  pushSender,
		log:         log,
	}
}

// CreateNotificationRequest represents notification creation request
type CreateNotificationRequest struct {
	UserID   string
	Type     string
	Category string
	Title    string
	Message  string
	Channels []string
	Metadata map[string]interface{}
}

// CreateAndSendNotification creates notification and sends via configured channels
func (s *NotificationService) CreateAndSendNotification(ctx context.Context, req *CreateNotificationRequest) error {
	s.log.WithFields(logrus.Fields{
		"user_id":  req.UserID,
		"type":     req.Type,
		"category": req.Category,
		"title":    req.Title,
	}).Info("Creating notification")

	// Create notification entity
	notification := &entities.Notification{
		ID:       uuid.New(),
		UserID:   uuid.MustParse(req.UserID),
		Type:     req.Type,
		Category: req.Category,
		Title:    req.Title,
		Message:  req.Message,
		Metadata: req.Metadata,
		Channels: req.Channels,
		Status:   "processing",
	}

	// Save to database
	if s.repo != nil {
		if err := s.repo.Create(ctx, notification); err != nil {
			s.log.WithError(err).Error("Failed to save notification to database")
			// Continue even if DB save fails (notification still sent)
		}
	}

	// Send via channels
	for _, channel := range req.Channels {
		switch channel {
		case "email":
			if err := s.sendEmail(ctx, req, notification); err != nil {
				s.log.WithError(err).Error("Failed to send email")
				notification.Status = "failed"
			}
		case "push":
			if err := s.sendPush(ctx, req, notification); err != nil {
				s.log.WithError(err).Error("Failed to send push")
				notification.Status = "failed"
			}
		case "in_app":
			// In-app already saved to DB
			s.log.Info("In-app notification saved")
		}
	}

	if notification.Status != "failed" {
		notification.Status = "sent"
	}

	// Update status in database
	if s.repo != nil {
		if err := s.repo.UpdateStatus(ctx, notification.ID, notification.Status); err != nil {
			s.log.WithError(err).Warn("Failed to update notification status")
		}
	}

	return nil
}

func (s *NotificationService) sendEmail(ctx context.Context, req *CreateNotificationRequest, notif *entities.Notification) error {
	payload := senders.NotificationPayload{
		To:      req.UserID, // In real: get user email from DB
		Subject: req.Title,
		Body:    s.formatEmailBody(req),
		Data:    req.Metadata,
	}

	if err := s.emailSender.Send(ctx, payload); err != nil {
		return fmt.Errorf("email send failed: %w", err)
	}

	// Update email_sent_at in DB
	if s.repo != nil {
		if err := s.repo.UpdateEmailSentAt(ctx, notif.ID); err != nil {
			s.log.WithError(err).Warn("Failed to update email sent timestamp")
		}
	}

	return nil
}

func (s *NotificationService) sendPush(ctx context.Context, req *CreateNotificationRequest, notif *entities.Notification) error {
	payload := senders.NotificationPayload{
		To:      req.UserID, // In real: get FCM token from DB
		Subject: req.Title,
		Body:    req.Message,
		Data:    req.Metadata,
	}

	if err := s.pushSender.Send(ctx, payload); err != nil {
		return fmt.Errorf("push send failed: %w", err)
	}

	// Update push_sent_at in DB
	if s.repo != nil {
		if err := s.repo.UpdatePushSentAt(ctx, notif.ID); err != nil {
			s.log.WithError(err).Warn("Failed to update push sent timestamp")
		}
	}

	return nil
}

func (s *NotificationService) formatEmailBody(req *CreateNotificationRequest) string {
	// Simple template replacement
	body := req.Message

	// Replace variables from metadata
	for key, value := range req.Metadata {
		placeholder := fmt.Sprintf("{{%s}}", key)
		body = strings.ReplaceAll(body, placeholder, fmt.Sprintf("%v", value))
	}

	return body
}

// Helper to interpolate template variables
func (s *NotificationService) interpolateTemplate(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// Get metadata as JSON string for logging
func metadataJSON(metadata map[string]interface{}) string {
	if metadata == nil {
		return "{}"
	}
	b, _ := json.MarshalIndent(metadata, "", "  ")
	return string(b)
}
