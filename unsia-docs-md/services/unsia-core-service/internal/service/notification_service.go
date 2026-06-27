package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

type Notification struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	UserID     string    `gorm:"column:user_id"`
	Title      string    `gorm:"column:title"`
	Message   string    `gorm:"column:message"`
	Type      string    `gorm:"column:type"` // info, warning, error, success
	IsRead    bool      `gorm:"column:is_read"`
	ActionURL  *string  `gorm:"column:action_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

type CreateNotificationInput struct {
	UserID    string  `json:"user_id" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Message  string  `json:"message" binding:"required"`
	Type     string  `json:"type"`
	ActionURL *string `json:"action_url"`
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(input CreateNotificationInput) (*Notification, error) {
	notifType := input.Type
	if notifType == "" {
		notifType = "info"
	}

	notif := Notification{
		ID:         uuid.New().String(),
		UserID:     input.UserID,
		Title:      input.Title,
		Message:   input.Message,
		Type:      notifType,
		IsRead:    false,
		ActionURL:  input.ActionURL,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&notif).Error; err != nil {
		return nil, err
	}

	return &notif, nil
}

// GetNotifications gets notifications for a user
func (s *NotificationService) GetNotifications(userID string, limit, offset int) ([]Notification, error) {
	var notifs []Notification
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifs).Error
	return notifs, err
}

// MarkAsRead marks notification as read
func (s *NotificationService) MarkAsRead(id string) error {
	return s.db.Model(&Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllAsRead marks all notifications as read for user
func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.db.Model(&Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

// GetUnreadCount gets unread notification count
func (s *NotificationService) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(id string) error {
	return s.db.Where("id = ?", id).Delete(&Notification{}).Error
}

// CleanOldNotifications deletes notifications older than days
func (s *NotificationService) CleanOldNotifications(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ? AND is_read = ?", cutoff, true).Delete(&Notification{})
	return result.RowsAffected, result.Error
}

// NotifyUser is a helper to notify a user
func (s *NotificationService) NotifyUser(userID, title, message, notifType string) error {
	input := CreateNotificationInput{
		UserID:   userID,
		Title:    title,
		Message: message,
		Type:    notifType,
	}
	_, err := s.CreateNotification(input)
	return err
}

// BroadcastNotification creates notification for all users in a list
func (s *NotificationService) BroadcastNotification(userIDs []string, title, message, notifType string) error {
	for _, userID := range userIDs {
		input := CreateNotificationInput{
			UserID:  userID,
			Title:   title,
			Message: message,
			Type:   notifType,
		}
		if _, err := s.CreateNotification(input); err != nil {
			continue // Continue with other users
		}
	}
	return nil
}

// EmailNotification creates notification and sends email
func (s *NotificationService) EmailNotification(userID, email, title, message, notifType string) error {
	// Create in-app notification
	if err := s.NotifyUser(userID, title, message, notifType); err != nil {
		return err
	}

	// TODO: Send email (need email service)
	// emailService := NewEmailService()
	// return emailService.SendEmail(email, title, message)

	return nil
}
