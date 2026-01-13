package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// NotificationRepository handles notification persistence
type NotificationRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewNotificationRepository(db *pgxpool.Pool, log *logrus.Logger) *NotificationRepository {
	return &NotificationRepository{
		db:  db,
		log: log,
	}
}

// Create inserts a new notification
func (r *NotificationRepository) Create(ctx context.Context, notif *entities.Notification) error {
	metadataJSON, err := json.Marshal(notif.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO notifications (
			id, user_id, type, category, title, message, metadata, 
			channels, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.db.Exec(ctx, query,
		notif.ID,
		notif.UserID,
		notif.Type,
		notif.Category,
		notif.Title,
		notif.Message,
		metadataJSON,
		notif.Channels,
		notif.Status,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}

	r.log.WithField("notification_id", notif.ID).Debug("Notification created")
	return nil
}

// GetUserNotifications retrieves notifications for a user
func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit int) ([]entities.Notification, error) {
	query := `
		SELECT id, user_id, type, category, title, message, metadata, channels,
		       status, is_read, read_at, email_sent_at, push_sent_at,
		       retry_count, last_error, created_at, updated_at, expires_at
		FROM notifications
		WHERE user_id = $1
		AND ($2::boolean IS FALSE OR is_read = FALSE)
		ORDER BY created_at DESC
		LIMIT $3
	`

	rows, err := r.db.Query(ctx, query, userID, unreadOnly, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []entities.Notification
	for rows.Next() {
		var notif entities.Notification
		var metadataJSON []byte

		err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.Type,
			&notif.Category,
			&notif.Title,
			&notif.Message,
			&metadataJSON,
			&notif.Channels,
			&notif.Status,
			&notif.IsRead,
			&notif.ReadAt,
			&notif.EmailSentAt,
			&notif.PushSentAt,
			&notif.RetryCount,
			&notif.LastError,
			&notif.CreatedAt,
			&notif.UpdatedAt,
			&notif.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		// Unmarshal metadata
		if err := json.Unmarshal(metadataJSON, &notif.Metadata); err != nil {
			r.log.WithError(err).Warn("Failed to unmarshal metadata")
			notif.Metadata = make(map[string]interface{})
		}

		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notifID, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, notifID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// UpdateStatus updates notification status
func (r *NotificationRepository) UpdateStatus(ctx context.Context, notifID uuid.UUID, status string) error {
	query := `
		UPDATE notifications
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, notifID, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// UpdateEmailSentAt updates email sent timestamp
func (r *NotificationRepository) UpdateEmailSentAt(ctx context.Context, notifID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET email_sent_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, notifID)
	return err
}

// UpdatePushSentAt updates push sent timestamp
func (r *NotificationRepository) UpdatePushSentAt(ctx context.Context, notifID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET push_sent_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, notifID)
	return err
}

// GetUnreadCount counts unread notifications
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread: %w", err)
	}

	return count, nil
}
