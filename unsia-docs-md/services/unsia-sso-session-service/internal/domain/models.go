package domain

import (
	"time"

	"gorm.io/gorm"
)

// Session represents a user session
type Session struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Token       string         `gorm:"uniqueIndex;not null" json:"token"`
	IPAddress   *string        `json:"ip_address,omitempty"`
	UserAgent   *string        `json:"user_agent,omitempty"`
	ExpiresAt   time.Time      `gorm:"not null" json:"expires_at"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	LastActive time.Time     `json:"last_active"`
	DeviceInfo  *string        `json:"device_info,omitempty"`
}

// ActiveSession represents active sessions for a user
type ActiveSession struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	Device   string `json:"device"`
	Location *string `json:"location,omitempty"`
	Active   bool   `json:"active"`
}

func (Session) TableName() string {
	return "sessions"
}

func (ActiveSession) TableName() string {
	return "active_sessions"
}
