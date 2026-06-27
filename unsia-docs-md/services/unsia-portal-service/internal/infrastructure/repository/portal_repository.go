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

func (r *PortalRepository) ListNews(limit, offset int) ([]domain.News, int64, error) {
	var list []domain.News
	var total int64
	r.db.Model(&domain.News{}).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Order("published_at desc").Find(&list).Error
	return list, total, err
}

func (r *PortalRepository) CreateNews(n *domain.News) error {
	return r.db.Create(n).Error
}

func (r *PortalRepository) ListAnnouncements(role string, limit, offset int) ([]domain.Announcement, int64, error) {
	var list []domain.Announcement
	var total int64
	query := r.db.Model(&domain.Announcement{})
	if role != "" {
		query = query.Where("target_role = ? OR target_role = 'all'", role)
	}
	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *PortalRepository) CreateAnnouncement(a *domain.Announcement) error {
	return r.db.Create(a).Error
}

func (r *PortalRepository) ListEvents(limit, offset int) ([]domain.Event, int64, error) {
	var list []domain.Event
	var total int64
	r.db.Model(&domain.Event{}).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Order("event_date asc").Find(&list).Error
	return list, total, err
}

func (r *PortalRepository) CreateEvent(e *domain.Event) error {
	return r.db.Create(e).Error
}

func (r *PortalRepository) GetMenusByRole(role string) ([]domain.Menu, error) {
	var list []domain.Menu
	err := r.db.Table("menus").
		Select("menus.*").
		Joins("JOIN role_menus ON role_menus.menu_code = menus.code").
		Where("role_menus.role_id = ?", role).
		Order("menus.sort_order asc").
		Find(&list).Error
	return list, err
}

