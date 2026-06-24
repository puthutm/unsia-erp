package repository

import (
	"time"

	"github.com/unsia-erp/unsia-sso-session-service/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) GetSessionByToken(token string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.Where("token = ? AND is_active = true AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetSessionsByUserID(userID uint) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.Where("user_id = ? AND is_active = true", userID).Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) DeleteSession(sessionID uint) error {
	return r.db.Model(&domain.Session{}).Where("id = ?", sessionID).Update("is_active", false).Error
}

func (r *SessionRepository) DeleteAllUserSessions(userID uint) error {
	return r.db.Model(&domain.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

func (r *SessionRepository) UpdateLastActive(sessionID uint, lastActive time.Time) error {
	return r.db.Model(&domain.Session{}).Where("id = ?", sessionID).Update("last_active", lastActive).Error
}

func (r *SessionRepository) UpsertActiveSession(userID uint, device, location string) error {
	return r.db.Where(domain.ActiveSession{
		UserID: userID,
	}).Assign(domain.ActiveSession{
		Device:   device,
		Location: &location,
		Active:   true,
	}).FirstOrCreate(&domain.ActiveSession{
		UserID: userID,
		Device: device,
	}).Error
}
