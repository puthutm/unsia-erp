package repository

import (
	"errors"

	"github.com/unsia-erp/unsia-portal-service/internal/domain"
	"gorm.io/gorm"
)

type PortalRepository struct {
	db *gorm.DB
}

func NewPortalRepository(db *gorm.DB) *PortalRepository {
	return &PortalRepository{db: db}
}

func (r *PortalRepository) CreateNotification(n *domain.Notification) error {
	return r.db.Create(n).Error
}

func (r *PortalRepository) ListNotificationsByUserID(userID string) ([]domain.Notification, error) {
	var list []domain.Notification
	err := r.db.Where("user_id = ?", userID).Order("sent_at desc").Find(&list).Error
	return list, err
}

func (r *PortalRepository) MarkNotificationAsRead(notificationID string, userID string) error {
	read := domain.NotificationRead{
		NotificationID: notificationID,
		UserID:         userID,
	}
	err := r.db.Create(&read).Error
	if err != nil {
		// If already read, ignore
		return nil
	}
	return nil
}

func (r *PortalRepository) SaveUserPreference(up *domain.UserPreference) error {
	var existing domain.UserPreference
	err := r.db.Where("user_id = ? AND preference_key = ?", up.UserID, up.PreferenceKey).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(up).Error
		}
		return err
	}

	existing.PreferenceValue = up.PreferenceValue
	existing.UpdatedAt = up.UpdatedAt
	return r.db.Save(&existing).Error
}

func (r *PortalRepository) SaveMenuShortcut(ms *domain.MenuShortcut) error {
	var existing domain.MenuShortcut
	err := r.db.Where("user_id = ? AND menu_code = ?", ms.UserID, ms.MenuCode).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(ms).Error
		}
		return err
	}

	existing.MenuLabel = ms.MenuLabel
	existing.TargetUrl = ms.TargetUrl
	existing.SortOrder = ms.SortOrder
	return r.db.Save(&existing).Error
}

func (r *PortalRepository) DeleteMenuShortcut(userID string, menuCode string) error {
	return r.db.Where("user_id = ? AND menu_code = ?", userID, menuCode).Delete(&domain.MenuShortcut{}).Error
}
