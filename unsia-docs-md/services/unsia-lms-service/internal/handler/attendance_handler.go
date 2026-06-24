package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"gorm.io/gorm"
)

type AttendanceCreateRequest struct {
	StudentID        string `json:"student_id" binding:"required"`
	AttendanceStatus string `json:"attendance_status" binding:"required"`
}

type AttendanceHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// CreateAttendance - Mengisi kehadiran mahasiswa pada sesi
func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	sessionID := c.Param("id")
	var req AttendanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	attendance := domain.Attendance{
		SessionID:        sessionID,
		StudentID:        req.StudentID,
		AttendanceStatus: req.AttendanceStatus,
		SubmittedAt:      time.Now(),
	}

	if err := h.repo.CreateAttendance(&attendance); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal填写hadiran").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(attendance).WithContext(c))
}
