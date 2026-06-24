package repository

import (
	"time"

	"github.com/unsia-erp/unsia-sso-token-service/internal/domain"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(refreshToken *domain.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

func (r *TokenRepository) GetRefreshToken(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.Where("token = ? AND is_active = true AND expires_at > ?", token, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *TokenRepository) RevokeRefreshToken(token string) error {
	now := time.Now()
	return r.db.Model(&domain.RefreshToken{}).Where("token = ?", token).Updates(map[string]interface{}{
		"is_active": false,
		"revoked_at": now,
	}).Error
}

func (r *TokenRepository) GetTokensByUserID(userID uint) ([]domain.RefreshToken, error) {
	var tokens []domain.RefreshToken
	err := r.db.Where("user_id = ? AND is_active = true", userID).Find(&tokens).Error
	return tokens, err
}
