package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"gorm.io/gorm"
)

type ClassSyncRequest struct {
	AcademicClassID string   `json:"academic_class_id" binding:"required"`
	LecturerIDs     []string `json:"lecturer_ids"`
}

type SessionCreateRequest struct {
	SessionNumber int    `json:"session_number" binding:"required"`
	Title         string `json:"title" binding:"required"`
	SessionDate   string `json:"session_date"` // YYYY-MM-DD
	StartTime     string `json:"start_time"`   // HH:MM
	EndTime       string `json:"end_time"`     // HH:MM
}

type SessionHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

func NewSessionHandler(db *gorm.DB) *SessionHandler {
	return &SessionHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// SyncClassFromAcademic - Sync kelas dari Academic Service
func (h *SessionHandler) SyncClassFromAcademic(c *gin.Context) {
	var req ClassSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var lecID *string
	if len(req.LecturerIDs) > 0 {
		lecID = &req.LecturerIDs[0]
	}

	class, err := h.repo.SyncClass(req.AcademicClassID, lecID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal melakukan sync kelas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(class).WithContext(c))
}

// CreateSession - Membuat sesi perkuliahan baru
func (h *SessionHandler) CreateSession(c *gin.Context) {
	classID := c.Param("id")
	var req SessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var sessionDate *time.Time
	if req.SessionDate != "" {
		parsed, err := time.Parse("2006-01-02", req.SessionDate)
		if err == nil {
			sessionDate = &parsed
		}
	}

	var startTime, endTime *string
	if req.StartTime != "" {
		startTime = &req.StartTime
	}
	if req.EndTime != "" {
		endTime = &req.EndTime
	}

	session := domain.Session{
		LmsClassID:    classID,
		SessionNumber: req.SessionNumber,
		Title:         req.Title,
		SessionDate:   sessionDate,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        "active",
	}

	if err := h.repo.CreateSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat sesi perkuliahan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(session).WithContext(c))
}
