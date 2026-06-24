package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-portal-service/internal/domain"
	"github.com/unsia-erp/unsia-portal-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type AnnouncementHandler struct {
	repo *repository.PortalRepository
	db   *gorm.DB
}

func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
	return &AnnouncementHandler{
		repo: repository.NewPortalRepository(db),
		db:   db,
	}
}

// GET /api/v1/portal/announcements - List announcements
func (h *AnnouncementHandler) ListAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	priority := c.Query("priority")
	status := c.Query("status")

	announcements, total, err := h.repo.ListAnnouncements(page, limit, priority, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(announcements).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/portal/announcements - Create announcement
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content  string `json:"content"`
		Priority string `json:"priority"` // normal, high, urgent
		Target   string `json:"target"`   // all, students, lecturers, staff
		StartDate string `json:"start_date"`
		EndDate  string `json:"end_date"`
		Status   string `json:"status"` // draft, published, archived
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	announcement := domain.Announcement{
		Title:     req.Title,
		Content:   req.Content,
		Priority:  req.Priority,
		Target:    req.Target,
		StartDate: req.StartDate,
		EndDate:  req.EndDate,
		Status:   req.Status,
	}

	if err := h.repo.CreateAnnouncement(&announcement); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(announcement).WithContext(c))
}

// GET /api/v1/portal/announcements/:id - Get announcement
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	announcement, err := h.repo.GetAnnouncementByID(announcementID)
	if err != nil || announcement == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengumuman tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(announcement).WithContext(c))
}

// PUT /api/v1/portal/announcements/:id - Update announcement
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	var req struct {
		Title    string `json:"title"`
		Content string `json:"content"`
		Priority string `json:"priority"`
		Target  string `json:"target"`
		Status  string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}
	if req.Target != "" {
		updates["target"] = req.Target
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.repo.UpdateAnnouncement(announcementID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Pengumuman berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/portal/announcements/:id - Delete announcement
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	if err := h.repo.DeleteAnnouncement(announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Pengumuman berhasil dihapus").WithContext(c))
}
