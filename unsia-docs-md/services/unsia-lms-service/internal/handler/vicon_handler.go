package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// VideoConferenceHandler - Handler for video conference (Vicon) features
type VideoConferenceHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

// NewVideoConferenceHandler - Create new video conference handler
func NewVideoConferenceHandler(db *gorm.DB) *VideoConferenceHandler {
	return &VideoConferenceHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// StartViconRequest - Request body for starting vicon
type StartViconRequest struct {
	SessionID   string  `json:"session_id" binding:"required"`
	MeetingURL  *string `json:"meeting_url"`  // External meeting URL (Zoom, Meet, etc)
}

// StartVideoConference - Start a video conference session
func (h *VideoConferenceHandler) StartVideoConference(c *gin.Context) {
	var req StartViconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if session exists
	var session domain.Session
	if err := h.db.First(&session, "id = ?", req.SessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi perkuliahan tidak ditemukan").WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	now := time.Now()
	videoSession := domain.VideoSession{
		SessionID:   req.SessionID,
		HostID:     userStr,
		Status:     "live",
		StartedAt:  &now,
		MeetingURL: "",
	}

	if req.MeetingURL != nil {
		videoSession.MeetingURL = *req.MeetingURL
	}

	if err := h.db.Create(&videoSession).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memulai video conference").WithContext(c))
		return
	}

	// Add host as participant
	participant := domain.VideoParticipant{
		VideoSessionID: videoSession.ID,
		UserID:        userStr,
		Role:         "host",
		JoinedAt:     now,
		IsMuted:      false,
		IsVideoOn:    true,
	}
	h.db.Create(&participant)

	c.JSON(http.StatusCreated, sharederr.Success(videoSession).WithContext(c))
}

// EndVideoConference - End a video conference session
func (h *VideoConferenceHandler) EndVideoConference(c *gin.Context) {
	videoSessionID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Check if video session exists and user is host
	var videoSession domain.VideoSession
	if err := h.db.First(&videoSession, "id = ?", videoSessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi video tidak ditemukan").WithContext(c))
		return
	}

	if videoSession.HostID != userStr {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya host yang dapat mengakhiri sesi").WithContext(c))
		return
	}

	// Calculate duration
	duration := int(time.Since(*videoSession.StartedAt).Seconds())

	now := time.Now()
	if err := h.db.Model(&videoSession).Updates(map[string]interface{}{
		"status":    "ended",
		"ended_at":  now,
		"duration": duration,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengakhiri video conference").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(videoSession, "Video conference berhasil diakhiri").WithContext(c))
}

// GetVideoSession - Get video conference session details
func (h *VideoConferenceHandler) GetVideoSession(c *gin.Context) {
	videoSessionID := c.Param("id")

	var videoSession domain.VideoSession
	if err := h.db.First(&videoSession, "id = ?", videoSessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi video tidak ditemukan").WithContext(c))
		return
	}

	// Get participants
	var participants []domain.VideoParticipant
	h.db.Where("video_session_id = ?", videoSessionID).Find(&participants)

	// Get linked session info
	var session domain.Session
	h.db.First(&session, "id = ?", videoSession.SessionID)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"video_session": videoSession,
		"session":      session,
		"participants": participants,
	}).WithContext(c))
}

// JoinVideoConference - Join a video conference
func (h *VideoConferenceHandler) JoinVideoConference(c *gin.Context) {
	videoSessionID := c.Param("id")

	// Check if video session exists
	var videoSession domain.VideoSession
	if err := h.db.First(&videoSession, "id = ?", videoSessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi video tidak ditemukan").WithContext(c))
		return
	}

	if videoSession.Status != "live" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATE", "Sesi video tidak sedang berlangsung").WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Check if already joined
	var existing domain.VideoParticipant
	if err := h.db.Where("video_session_id = ? AND user_id = ?", videoSessionID, userStr).First(&existing).Error; err == nil {
		// Update join time if re-joining
		h.db.Model(&existing).Update("joined_at", time.Now())
		c.JSON(http.StatusOK, sharederr.Success(existing).WithContext(c))
		return
	}

	now := time.Now()
	participant := domain.VideoParticipant{
		VideoSessionID: videoSessionID,
		UserID:        userStr,
		Role:         "participant",
		JoinedAt:     now,
		IsMuted:      false,
		IsVideoOn:    true,
	}

	if err := h.db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal bergabung ke video conference").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(participant).WithContext(c))
}

// LeaveVideoConference - Leave a video conference
func (h *VideoConferenceHandler) LeaveVideoConference(c *gin.Context) {
	videoSessionID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	now := time.Now()
	if err := h.db.Model(&domain.VideoParticipant{}).
		Where("video_session_id = ? AND user_id = ?", videoSessionID, userStr).
		Updates(map[string]interface{}{
			"left_at": now,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal keluar dari video conference").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Berhasil keluar dari video conference").WithContext(c))
}

// GetVideoSessionsBySession - Get all video sessions for a class session
func (h *VideoConferenceHandler) GetVideoSessionsBySession(c *gin.Context) {
	sessionID := c.Query("session_id")

	var videoSessions []domain.VideoSession
	query := h.db.Model(&domain.VideoSession{})

	if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	if err := query.Order("created_at DESC").Find(&videoSessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data video conference").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(videoSessions).WithContext(c))
}

// ToggleMute - Toggle mute status
func (h *VideoConferenceHandler) ToggleMute(c *gin.Context) {
	videoSessionID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	var participant domain.VideoParticipant
	if err := h.db.Where("video_session_id = ? AND user_id = ?", videoSessionID, userStr).First(&participant).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Peserta tidak ditemukan").WithContext(c))
		return
	}

	newMuteStatus := !participant.IsMuted
	if err := h.db.Model(&participant).Update("is_muted", newMuteStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah status mute").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"is_muted": newMuteStatus}).WithContext(c))
}

// ToggleVideo - Toggle video status
func (h *VideoConferenceHandler) ToggleVideo(c *gin.Context) {
	videoSessionID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	var participant domain.VideoParticipant
	if err := h.db.Where("video_session_id = ? AND user_id = ?", videoSessionID, userStr).First(&participant).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Peserta tidak ditemukan").WithContext(c))
		return
	}

	newVideoStatus := !participant.IsVideoOn
	if err := h.db.Model(&participant).Update("is_video_on", newVideoStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah status video").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"is_video_on": newVideoStatus}).WithContext(c))
}
