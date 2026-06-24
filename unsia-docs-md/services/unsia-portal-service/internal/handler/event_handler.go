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

type EventHandler struct {
	repo *repository.PortalRepository
	db   *gorm.DB
}

func NewEventHandler(db *gorm.DB) *EventHandler {
	return &EventHandler{
		repo: repository.NewPortalRepository(db),
		db:   db,
	}
}

// GET /api/v1/portal/events - List events
func (h *EventHandler) ListEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	eventType := c.Query("type")
	status := c.Query("status")

	events, total, err := h.repo.ListEvents(page, limit, eventType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data acara").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(events).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/portal/events - Create event
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Location   string  `json:"location"`
		StartDate  string   `json:"start_date"`
		EndDate    string   `json:"end_date"`
		EventType  string   `json:"type"` // seminar, workshop, webinar, gathering
		Organizer  string   `json:"organizer"`
		Max Participants int     `json:"max_participants"`
		Status     string   `json:"status"` // draft, published, cancelled, completed
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	event := domain.Event{
		Name:            req.Name,
		Description:     req.Description,
		Location:       req.Location,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		EventType:      req.EventType,
		Organizer:      req.Organizer,
		MaxParticipants: req.MaxParticipants,
		Status:         req.Status,
	}

	if err := h.repo.CreateEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat acara").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(event).WithContext(c))
}

// GET /api/v1/portal/events/:id - Get event
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")

	event, err := h.repo.GetEventByID(eventID)
	if err != nil || event == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Acara tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(event).WithContext(c))
}

// PUT /api/v1/portal/events/:id - Update event
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")

	var req struct {
		Name            string `json:"name"`
		Description    string `json:"description"`
		Location      string `json:"location"`
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		EventType     string `json:"type"`
		Organizer     string `json:"organizer"`
		MaxParticipants int `json:"max_participants"`
		Status        string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.StartDate != "" {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		updates["end_date"] = req.EndDate
	}
	if req.EventType != "" {
		updates["event_type"] = req.EventType
	}
	if req.Organizer != "" {
		updates["organizer"] = req.Organizer
	}
	if req.MaxParticipants > 0 {
		updates["max_participants"] = req.MaxParticipants
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.repo.UpdateEvent(eventID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui acara").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Acara berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/portal/events/:id - Delete event
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	if err := h.repo.DeleteEvent(eventID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus acara").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Acara berhasil dihapus").WithContext(c))
}
