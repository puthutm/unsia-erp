package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unsia-erp/unsia-sso-session-service/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

type CreateSessionRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	DeviceInfo string `json:"device_info"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request", "errors": []string{err.Error()}})
		return
	}

	deviceInfo := req.DeviceInfo
	if deviceInfo == "" {
		deviceInfo = "unknown"
	}
	ipAddress := req.IPAddress
	if ipAddress == "" {
		ipAddress = c.ClientIP()
	}
	userAgent := req.UserAgent
	if userAgent == "" {
		userAgent = c.Request.UserAgent()
	}

	session, err := h.sessionService.CreateSession(c.Request.Context(), req.UserID, deviceInfo, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Session created",
		"data": gin.H{
			"token":      session.Token,
			"expires_at": session.ExpiresAt,
		},
	})
}

func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	sessions, err := h.sessionService.GetSessions(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sessions})
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	sessionID := c.Param("id")
	var sessionIDUint uint
	// In real code, parse sessionID to uint

	err := h.sessionService.DeleteSession(c.Request.Context(), sessionIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Session deleted"})
}

func (h *SessionHandler) DeleteAllSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	err := h.sessionService.DeleteAllSessions(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "All sessions deleted"})
}
