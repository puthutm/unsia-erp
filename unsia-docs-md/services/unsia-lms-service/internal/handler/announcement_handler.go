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

// AnnouncementHandler - Handler for pengumuman/announcement features
type AnnouncementHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

// NewAnnouncementHandler - Create new announcement handler
func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
	return &AnnouncementHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// CreateAnnouncementRequest - Request body for creating announcement
type CreateAnnouncementRequest struct {
	LmsClassID *string `json:"lms_class_id"`
	SessionID *string `json:"session_id"`
	Title     string  `json:"title" binding:"required"`
	Content   string  `json:"content" binding:"required"`
	IsPinned  *bool   `json:"is_pinned"`
	Priority string  `json:"priority"` // urgent, high, normal, low
	StartDate *string `json:"start_date"` // RFC3339
	EndDate   *string `json:"end_date"`   // RFC3339
}

// CreateAnnouncement - Create new announcement
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Parse dates
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.StartDate)
		if err == nil {
			startDate = &parsed
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.EndDate)
		if err == nil {
			endDate = &parsed
		}
	}

	priority := "normal"
	if req.Priority != "" {
		priority = req.Priority
	}

	isPinned := false
	if req.IsPinned != nil {
		isPinned = *req.IsPinned
	}

	announcement := domain.Announcement{
		LmsClassID: req.LmsClassID,
		SessionID: req.SessionID,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userStr,
		IsPinned:  isPinned,
		IsActive:  true,
		Priority: priority,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := h.db.Create(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(announcement).WithContext(c))
}

// GetAnnouncements - Get announcements (with filters)
func (h *AnnouncementHandler) GetAnnouncements(c *gin.Context) {
	classID := c.Query("class_id")
	sessionID := c.Query("session_id")
	priority := c.Query("priority")
	activeOnly := c.Query("active") == "true"

	query := h.db.Model(&domain.Announcement{})

	if classID != "" {
		query = query.Where("lms_class_id = ?", classID)
	}
	if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if activeOnly {
		now := time.Now()
		query = query.Where("is_active = ?", true).
			Where("(start_date IS NULL OR start_date <= ?)", now).
			Where("(end_date IS NULL OR end_date >= ?)", now)
	}

	var announcements []domain.Announcement
	if err := query.Order("is_pinned DESC, created_at DESC").Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(announcements).WithContext(c))
}

// GetAnnouncement - Get single announcement
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement domain.Announcement
	if err := h.db.First(&announcement, "id = ?", announcementID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengumuman tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(announcement).WithContext(c))
}

// UpdateAnnouncement - Update announcement
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	// Check if announcement exists
	var announcement domain.Announcement
	if err := h.db.First(&announcement, "id = ?", announcementID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengumuman tidak ditemukan").WithContext(c))
		return
	}

	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := map[string]interface{}{
		"title":    req.Title,
		"content":  req.Content,
		"priority": req.Priority,
	}

	if req.IsPinned != nil {
		updates["is_pinned"] = *req.IsPinned
	}
	if req.StartDate != nil {
		parsed, _ := time.Parse(time.RFC3339, *req.StartDate)
		updates["start_date"] = parsed
	}
	if req.EndDate != nil {
		parsed, _ := time.Parse(time.RFC3339, *req.EndDate)
		updates["end_date"] = parsed
	}

	if err := h.db.Model(&announcement).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(announcement).WithContext(c))
}

// DeleteAnnouncement - Delete announcement (soft delete)
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Check if announcement exists
	var announcement domain.Announcement
	if err := h.db.First(&announcement, "id = ?", announcementID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengumuman tidak ditemukan").WithContext(c))
		return
	}

	// Check if user is author or admin
	if announcement.AuthorID != userStr {
		// In real implementation, check if user has admin role
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Tidak memiliki akses untuk menghapus pengumuman ini").WithContext(c))
		return
	}

	if err := h.db.Model(&announcement).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Pengumuman berhasil dihapus").WithContext(c))
}

// PinAnnouncement - Pin/unpin announcement
func (h *AnnouncementHandler) PinAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement domain.Announcement
	if err := h.db.First(&announcement, "id = ?", announcementID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengumuman tidak ditemukan").WithContext(c))
		return
	}

	newPinnedStatus := !announcement.IsPinned
	if err := h.db.Model(&announcement).Update("is_pinned", newPinnedStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah status pin").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"is_pinned": newPinnedStatus}).WithContext(c))
}
