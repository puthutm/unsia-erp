package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unsia-erp/unsia-sso-token-service/internal/domain"
	"github.com/unsia-erp/unsia-sso-token-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type TokenService struct {
	db  *gorm.DB
	rep *repository.TokenRepository
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{
		db:  db,
		rep: repository.NewTokenRepository(db),
	}
}

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *TokenService) GenerateAccessToken(userID uint, email string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "unsia-erp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func (s *TokenService) GenerateRefreshToken(ctx context.Context, userID uint, deviceInfo string) (*domain.RefreshToken, error) {
	token := "refresh_" + time.Now().Format("20060102150405") + randomString(32)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	refreshToken := &domain.RefreshToken{
		UserID:    userID,
		Token:    token,
		ExpiresAt: expiresAt,
		IsActive:  true,
		DeviceInfo: &deviceInfo,
	}

	err := s.rep.Create(refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	token, err := s.rep.GetRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	if !token.IsActive || token.ExpiresAt.Before(time.Now()) {
		return "", gorm.ErrRecordNotFound
	}

	// In real code, get user info from user service
	accessToken, err := s.GenerateAccessToken(token.UserID, "user@example.com")
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *TokenService) RevokeToken(ctx context.Context, refreshToken string) error {
	return s.rep.RevokeRefreshToken(refreshToken)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
