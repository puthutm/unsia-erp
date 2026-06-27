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

type MenuResponse struct {
	MenuCode string         `json:"menu_code"`
	Label    string         `json:"label"`
	Path     string         `json:"path"`
	Icon     string         `json:"icon"`
	Children []MenuResponse `json:"children,omitempty"`
}

type Menu struct {
	Code               string  `gorm:"primaryKey;column:code"`
	Label              string  `gorm:"column:label;not null"`
	Path               string  `gorm:"column:path;not null"`
	Icon               string  `gorm:"column:icon"`
	ParentCode         *string `gorm:"column:parent_code"`
	SortOrder          int     `gorm:"column:sort_order;default:0;not null"`
	RequiredPermission string  `gorm:"column:required_permission"`
}

func (Menu) TableName() string {
	return "menus"
}

type RoleMenu struct {
	RoleID   string `gorm:"primaryKey;column:role_id"`
	MenuCode string `gorm:"primaryKey;column:menu_code"`
}

func (RoleMenu) TableName() string {
	return "role_menus"
}

type News struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Title       string    `gorm:"column:title;not null"`
	Content     string    `gorm:"column:content;not null"`
	Author      string    `gorm:"column:author"`
	PublishedAt time.Time `gorm:"column:published_at"`
}

func (News) TableName() string {
	return "news"
}

type Announcement struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Title      string    `gorm:"column:title;not null"`
	Message    string    `gorm:"column:message;not null"`
	TargetRole string    `gorm:"column:target_role;default:'all';not null"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (Announcement) TableName() string {
	return "announcements"
}

type Event struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Title       string    `gorm:"column:title;not null"`
	Description string    `gorm:"column:description"`
	EventDate   time.Time `gorm:"column:event_date;not null"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (Event) TableName() string {
	return "events"
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
