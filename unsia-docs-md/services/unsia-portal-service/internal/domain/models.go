package domain

import (
	"time"
)

type Notification struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID       string    `gorm:"column:user_id;not null"` // external_ref: core.users.id
	Title        string    `gorm:"column:title;not null"`
	Message      string    `gorm:"column:message"`
	ModuleSource string    `gorm:"column:module_source"`
	TargetUrl    string    `gorm:"column:target_url"`
	SentAt       time.Time `gorm:"column:sent_at;default:now()"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationRead struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	NotificationID string    `gorm:"column:notification_id;not null"`
	UserID         string    `gorm:"column:user_id;not null"` // external_ref: core.users.id
	ReadAt         time.Time `gorm:"column:read_at;default:now()"`
}

func (NotificationRead) TableName() string {
	return "notification_reads"
}

type UserPreference struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID          string    `gorm:"column:user_id;not null"`
	PreferenceKey   string    `gorm:"column:preference_key;not null"`
	PreferenceValue string    `gorm:"type:jsonb;column:preference_value;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}

type MenuShortcut struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID    string `gorm:"column:user_id;not null"`
	MenuCode  string `gorm:"column:menu_code;not null"`
	MenuLabel string `gorm:"column:menu_label;not null"`
	TargetUrl string `gorm:"column:target_url;not null"`
	SortOrder int    `gorm:"column:sort_order;default:0;not null"`
}

func (MenuShortcut) TableName() string {
	return "menu_shortcuts"
}
