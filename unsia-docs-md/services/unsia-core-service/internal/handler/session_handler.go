package handler

import (
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
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Token       string `json:"token,omitempty"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	Fingerprint string `json:"fingerprint"`
	ExpiresAt   string `json:"expires_at"`
	RevokedAt   *string `json:"revoked_at,omitempty"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	LastActivity string `json:"last_activity"`
}

// GetSession handles GET /api/v1/sessions/:id
func (h *SessionHandler) GetSession(c *gin.Context) {
	id := c.Param("id")

	var session Session
	if err := h.db.Where("id = ?", id).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Session not found").WithContext(c))
		return
	}

	session.Token = ""
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

	// Clear tokens
	for i := range sessions {
		sessions[i].Token = ""
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

	now := getTimestamp()
	if err := h.db.Model(&session).Update("revoked_at", now).Error; err != nil {
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

	now := getTimestamp()
	if err := h.db.Model(&Session{}).Where("user_id = ?", userID).Update("revoked_at", now).Error; err != nil {
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

	// Find session by refresh token
	var session Session
	if err := h.db.Where("token = ? AND is_active = ?", req.RefreshToken, true).First(&session).Error; err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_TOKEN", "Invalid refresh token").WithContext(c))
		return
	}

	// Check if expired
	if session.ExpiresAt < getTimestamp() {
		c.JSON(http.StatusUnauthorized, sharederr.Error("TOKEN_EXPIRED", "Refresh token expired").WithContext(c))
		return
	}

	// Generate new tokens
	newToken := uuid.New().String()
	newRefreshToken := uuid.New().String()

	// Update session
	updates := map[string]interface{}{
		"token":          newToken,
		"expires_at":    getExpireTime(24 * 60), // 24 hours
		"last_activity": getTimestamp(),
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
	var session Session
	err := h.db.Where("token = ? AND is_active = ?", token, true).First(&session).Error
	return &session, err
}

// CleanExpiredSessions cleans expired sessions
func (h *SessionHandler) CleanExpiredSessions() error {
	now := getTimestamp()
	return h.db.Where("expires_at < ?", now).Delete(&Session{}).Error
}
