package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/gorm"
)

type SessionHandler struct {
	db *gorm.DB
}

func NewSessionHandler(db *gorm.DB) *SessionHandler {
	return &SessionHandler{db: db}
}

type Session struct {
	ID               string     `gorm:"column:id;primaryKey" json:"id"`
	UserID           string     `gorm:"column:user_id" json:"user_id"`
	TokenHash        string     `gorm:"column:token_hash" json:"-"`
	RefreshTokenHash string     `gorm:"column:refresh_token_hash" json:"-"`
	ExpiredAt        time.Time  `gorm:"column:expired_at" json:"expires_at"`
	RevokedAt        *time.Time `gorm:"column:revoked_at" json:"revoked_at,omitempty"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	IsActive         bool       `gorm:"-" json:"is_active"`
}

// GetSession handles GET /api/v1/sessions/:id
func (h *SessionHandler) GetSession(c *gin.Context) {
	id := c.Param("id")

	var session Session
	if err := h.db.Where("id = ?", id).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Session not found").WithContext(c))
		return
	}

	session.IsActive = (session.RevokedAt == nil) && session.ExpiredAt.After(time.Now())
	c.JSON(http.StatusOK, sharederr.Success(session).WithContext(c))
}

// ListSessions handles GET /api/v1/sessions
func (h *SessionHandler) ListSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var sessions []Session
	query := h.db.Where("user_id = ?", userID)
	
	if err := query.Order("created_at DESC").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	for i := range sessions {
		sessions[i].IsActive = (sessions[i].RevokedAt == nil) && sessions[i].ExpiredAt.After(time.Now())
	}

	c.JSON(http.StatusOK, sharederr.Success(sessions).WithContext(c))
}

// DeleteSession handles DELETE /api/v1/sessions/:id
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	// Verify session belongs to user
	var session Session
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Session not found").WithContext(c))
		return
	}

	now := time.Now()
	if err := h.db.Model(&session).Update("revoked_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Session revoked",
	}).WithContext(c))
}

// RevokeAllSessions handles DELETE /api/v1/sessions
func (h *SessionHandler) RevokeAllSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	now := time.Now()
	if err := h.db.Model(&Session{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "All sessions revoked",
	}).WithContext(c))
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *SessionHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	hashRefresh := sha256.Sum256([]byte(req.RefreshToken))
	refreshTokenHash := hex.EncodeToString(hashRefresh[:])

	// Find session by refresh token hash
	var session Session
	if err := h.db.Where("refresh_token_hash = ? AND revoked_at IS NULL", refreshTokenHash).First(&session).Error; err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_TOKEN", "Invalid refresh token").WithContext(c))
		return
	}

	// Check if expired
	if session.ExpiredAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, sharederr.Error("TOKEN_EXPIRED", "Refresh token expired").WithContext(c))
		return
	}

	// Generate new tokens
	newToken := uuid.New().String()
	hashAccess := sha256.Sum256([]byte(newToken))
	accessTokenHash := hex.EncodeToString(hashAccess[:])

	// Update session
	updates := map[string]interface{}{
		"token_hash": accessTokenHash,
		"expired_at": time.Now().Add(24 * time.Hour), // 24 hours
	}

	if err := h.db.Model(&session).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"access_token":  newToken,
		"token_type":  "Bearer",
		"expires_in": 86400,
	}).WithContext(c))
}

// SessionFromToken gets session from token
func (h *SessionHandler) SessionFromToken(token string) (*Session, error) {
	hashAccess := sha256.Sum256([]byte(token))
	accessTokenHash := hex.EncodeToString(hashAccess[:])

	var session Session
	err := h.db.Where("token_hash = ? AND revoked_at IS NULL", accessTokenHash).First(&session).Error
	return &session, err
}

// CleanExpiredSessions cleans expired sessions
func (h *SessionHandler) CleanExpiredSessions() error {
	return h.db.Where("expired_at < ?", time.Now()).Delete(&Session{}).Error
}
