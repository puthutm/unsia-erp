package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/unsia-erp/unsia-core-service/internal/domain"
	"gorm.io/gorm"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// SessionResponse represents session info for API response
type SessionResponse struct {
	ID        string    `json:"id"`
	UserID   string    `json:"user_id"`
	Active   bool      `json:"active"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateSession creates a new session
func (s *SessionService) CreateSession(userID, refreshToken string, expiresAt time.Time) (*Session, error) {
	hashRefresh := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hashRefresh[:])

	session := Session{
		ID:               uuid.New().String(),
		UserID:           userID,
		TokenHash:        "dummy-" + uuid.New().String(),
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        expiresAt,
		CreatedAt:        time.Now(),
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal membuat session")
	}

	return &session, nil
}

// GetSessionByID retrieves session by ID
func (s *SessionService) GetSessionByID(id string) (*Session, error) {
	var session Session
	if err := s.db.First(&session, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("SESSION_NOT_FOUND: Session tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil session")
	}
	return &session, nil
}

// GetSessionByRefreshToken retrieves session by refresh token
func (s *SessionService) GetSessionByRefreshToken(refreshToken string) (*Session, error) {
	hashRefresh := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hashRefresh[:])

	var session Session
	if err := s.db.Where("refresh_token_hash = ?", refreshTokenHash).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("SESSION_NOT_FOUND: Session tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil session")
	}
	return &session, nil
}

// ValidateSession checks if session is valid
func (s *SessionService) ValidateSession(refreshToken string) (*Session, error) {
	session, err := s.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if session.IsRevoked {
		return nil, errors.New("SESSION_REVOKED: Session telah dicabut")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("SESSION_EXPIRED: Session telah kedaluwarsa")
	}

	return session, nil
}

// RevokeSession revokes a session
func (s *SessionService) RevokeSession(id string) error {
	session, err := s.GetSessionByID(id)
	if err != nil {
		return err
	}

	session.IsRevoked = true

	if err := s.db.Save(session).Error; err != nil {
		return errors.New("DB_ERROR: Gagal mencabut session")
	}

	return nil
}

// RevokeUserSessions revokes all sessions for a user
func (s *SessionService) RevokeUserSessions(userID string) error {
	result := s.db.Model(&Session{}).Where("user_id = ? AND is_revoked = false", userID).
		Update("is_revoked", true)
	if result.Error != nil {
		return errors.New("DB_ERROR: Gagal mencabut semua session")
	}

	return nil
}

// ListUserSessions returns all active sessions for a user
func (s *SessionService) ListUserSessions(userID string) ([]Session, error) {
	var sessions []Session
	if err := s.db.Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengambil daftar session")
	}

	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions
func (s *SessionService) CleanupExpiredSessions() (int64, error) {
	result := s.db.Where("expires_at < ? AND is_revoked = false", time.Now()).
		Model(&Session{}).
		Update("is_revoked", true)
	if result.Error != nil {
		return 0, errors.New("DB_ERROR: Gagal membersihkan session expired")
	}

	return result.RowsAffected, nil
}

// =============================================================================
// ACTIVE ROLE SESSION SERVICE
// =============================================================================

type ActiveRoleSessionService struct {
	db *gorm.DB
}

func NewActiveRoleSessionService(db *gorm.DB) *ActiveRoleSessionService {
	return &ActiveRoleSessionService{db: db}
}

// CreateActiveRoleSession creates a new active role session
func (s *ActiveRoleSessionService) CreateActiveRoleSession(userID, roleID, sessionID string, applicationID *string, studyProgramID *string) (*ActiveRoleSession, error) {
	appID := ""
	if applicationID != nil {
		appID = *applicationID
	}
	activeSession := ActiveRoleSession{
		ID:             uuid.New().String(),
		UserID:         userID,
		RoleID:         roleID,
		SessionID:      sessionID,
		ApplicationID:  appID,
		StudyProgramID: studyProgramID,
		CreatedAt:      time.Now(),
	}

	if err := s.db.Create(&activeSession).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal membuat active role session")
	}

	return &activeSession, nil
}

// GetActiveRoleSessionBySessionID retrieves active role session by session ID
func (s *ActiveRoleSessionService) GetActiveRoleSessionBySessionID(sessionID string) (*ActiveRoleSession, error) {
	var activeSession ActiveRoleSession
	if err := s.db.Preload("Role").First(&activeSession, "session_id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ACTIVE_ROLE_SESSION_NOT_FOUND: Active role session tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil active role session")
	}
	return &activeSession, nil
}

// UpdateActiveRoleSession updates active role session
func (s *ActiveRoleSessionService) UpdateActiveRoleSession(sessionID, roleID string, applicationID *string, studyProgramID *string) (*ActiveRoleSession, error) {
	activeSession, err := s.GetActiveRoleSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	activeSession.RoleID = roleID
	if applicationID != nil {
		activeSession.ApplicationID = *applicationID
	} else {
		activeSession.ApplicationID = ""
	}
	activeSession.StudyProgramID = studyProgramID

	if err := s.db.Save(&activeSession).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengupdate active role session")
	}

	return activeSession, nil
}

// DeleteActiveRoleSession deletes active role session
func (s *ActiveRoleSessionService) DeleteActiveRoleSession(sessionID string) error {
	result := s.db.Where("session_id = ?", sessionID).Delete(&ActiveRoleSession{})
	if result.Error != nil {
		return errors.New("DB_ERROR: Gagal menghapus active role session")
	}
	if result.RowsAffected == 0 {
		return errors.New("ACTIVE_ROLE_SESSION_NOT_FOUND: Active role session tidak ditemukan")
	}
	return nil
}

// =============================================================================
// SERVICE TOKEN SERVICE
// =============================================================================

type ServiceTokenService struct {
	db *gorm.DB
}

func NewServiceTokenService(db *gorm.DB) *ServiceTokenService {
	return &ServiceTokenService{db: db}
}

type CreateServiceTokenInput struct {
	ApplicationID string   `json:"application_id" binding:"required"`
	Scopes       []string `json:"scopes"`
	ExpiresAt    *time.Time `json:"expires_at"`
}

type ServiceTokenResponse struct {
	ID           string    `json:"id"`
	ApplicationID string  `json:"application_id"`
	Token        string    `json:"token"`
	Scopes       []string `json:"scopes"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateServiceToken creates a new service token
func (s *ServiceTokenService) CreateServiceToken(input CreateServiceTokenInput) (*ServiceTokenResponse, error) {
	// Generate token
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes) // Use crypto/rand
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Hash token for storage
	tokenHash := hashToken(token)

	expiresAt := time.Now().AddDate(1, 0, 0) // Default 1 year
	if input.ExpiresAt != nil {
		expiresAt = *input.ExpiresAt
	}

	serviceToken := ServiceToken{
		ID:            uuid.New().String(),
		ApplicationID: input.ApplicationID,
		TokenHash:     tokenHash,
		Scopes:       toJSONB(input.Scopes),
		ExpiresAt:     expiresAt,
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(&serviceToken).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal membuat service token")
	}

	return &ServiceTokenResponse{
		ID:            serviceToken.ID,
		ApplicationID: serviceToken.ApplicationID,
		Token:        token,
		Scopes:       input.Scopes,
		ExpiresAt:    serviceToken.ExpiresAt,
		CreatedAt:    serviceToken.CreatedAt,
	}, nil
}

// GetServiceTokenByID retrieves service token by ID
func (s *ServiceTokenService) GetServiceTokenByID(id string) (*ServiceToken, error) {
	var st ServiceToken
	if err := s.db.First(&st, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("SERVICE_TOKEN_NOT_FOUND: Service token tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil service token")
	}
	return &st, nil
}

// ValidateServiceToken validates a service token
func (s *ServiceTokenService) ValidateServiceToken(token string) (*ServiceToken, error) {
	tokenHash := hashToken(token)

	var st ServiceToken
	if err := s.db.Where("token_hash = ?", tokenHash).First(&st).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("SERVICE_TOKEN_INVALID: Service token tidak valid")
		}
		return nil, errors.New("DB_ERROR: Gagal memvalidasi service token")
	}

	if st.RevokedAt != nil {
		return nil, errors.New("SERVICE_TOKEN_REVOKED: Service token telah dicabut")
	}

	if st.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("SERVICE_TOKEN_EXPIRED: Service token telah kedaluwarsa")
	}

	return &st, nil
}

// ListServiceTokens returns all service tokens for an application
func (s *ServiceTokenService) ListServiceTokens(applicationID string) ([]ServiceToken, error) {
	var tokens []ServiceToken
	if err := s.db.Where("application_id = ?", applicationID).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengambil daftar service token")
	}
	return tokens, nil
}

// RevokeServiceToken revokes a service token
func (s *ServiceTokenService) RevokeServiceToken(id string) error {
	st, err := s.GetServiceTokenByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	st.RevokedAt = &now

	if err := s.db.Save(st).Error; err != nil {
		return errors.New("DB_ERROR: Gagal mencabut service token")
	}

	return nil
}

// RotateServiceToken rotates a service token (creates new one and revokes old)
func (s *ServiceTokenService) RotateServiceToken(id string) (*ServiceTokenResponse, error) {
	st, err := s.GetServiceTokenByID(id)
	if err != nil {
		return nil, err
	}

	// Revoke old token
	now := time.Now()
	st.RevokedAt = &now
	if err := s.db.Save(st).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mencabut service token lama")
	}

	// Create new token with same application and scopes
	input := CreateServiceTokenInput{
		ApplicationID: st.ApplicationID,
		Scopes:       fromJSONB(st.Scopes),
		ExpiresAt:    &st.ExpiresAt,
	}

	return s.CreateServiceToken(input)
}

func hashToken(token string) string {
	// Simple hash for storage - in production use proper hashing
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func toJSONB(scopes []string) string {
	// GORM handles []string as jsonb automatically when using pq.StringArray or similar
	// For now we'll use jsonb directly
	return "[]" // Simplified
}

func fromJSONB(scopes string) []string {
	return []string{} // Simplified
}
