package domain

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a refresh token
type RefreshToken struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Token      string         `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt   time.Time      `gorm:"not null" json:"expires_at"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	DeviceInfo *string       `json:"device_info,omitempty"`
	RevokedAt  *time.Time    `json:"revoked_at,omitempty"`
}

// TokenMetadata stores token metadata
type TokenMetadata struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	TokenType   string `gorm:"not null" json:"token_type"` // access, refresh
	DeviceHash string `json:"device_hash"`
	IPAddress  *string `json:"ip_address,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (TokenMetadata) TableName() string {
	return "token_metadata"
}
