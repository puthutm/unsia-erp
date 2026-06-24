package service

import (
	"context"
	"time"

	"github.com/unsia-erp/unsia-sso-session-service/internal/domain"
	"github.com/unsia-erp/unsia-sso-session-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type SessionService struct {
	db  *gorm.DB
	rep *repository.SessionRepository
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{
		db:  db,
		rep: repository.NewSessionRepository(db),
	}
}

func (s *SessionService) CreateSession(ctx context.Context, userID uint, deviceInfo, ipAddress, userAgent string) (*domain.Session, error) {
	token := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	session := &domain.Session{
		UserID:    userID,
		Token:    token,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		ExpiresAt: expiresAt,
		IsActive:  true,
		LastActive: time.Now(),
		DeviceInfo: &deviceInfo,
	}

	err := s.rep.Create(session)
	if err != nil {
		return nil, err
	}

	// Update active session
	s.rep.UpsertActiveSession(userID, deviceInfo, "")

	return session, nil
}

func (s *SessionService) GetSessions(ctx context.Context, userID uint) ([]domain.Session, error) {
	return s.rep.GetSessionsByUserID(userID)
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID uint) error {
	return s.rep.DeleteSession(sessionID)
}

func (s *SessionService) DeleteAllSessions(ctx context.Context, userID uint) error {
	return s.rep.DeleteAllUserSessions(userID)
}

func (s *SessionService) ValidateSession(ctx context.Context, token string) (*domain.Session, error) {
	return s.rep.GetSessionByToken(token)
}

func (s *SessionService) UpdateLastActive(ctx context.Context, sessionID uint) error {
	return s.rep.UpdateLastActive(sessionID, time.Now())
}

func generateToken() string {
	return "sess_" + time.Now().Format("20060102150405") + randomString(32)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
