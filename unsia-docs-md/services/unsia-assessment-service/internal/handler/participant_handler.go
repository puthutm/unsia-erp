package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-assessment-service/internal/domain"
	"github.com/unsia-erp/unsia-assessment-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type ParticipantHandler struct {
	repo *repository.AssessmentRepository
	db   *gorm.DB
}

func NewParticipantHandler(db *gorm.DB) *ParticipantHandler {
	return &ParticipantHandler{
		repo: repository.NewAssessmentRepository(db),
		db:   db,
	}
}

// GET /api/v1/participants - List participants
func (h *ParticipantHandler) ListParticipants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sessionID := c.Query("session_id")
	studentID := c.Query("student_id")
	status := c.Query("status")

	participants, total, err := h.repo.ListParticipants(page, limit, sessionID, studentID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data partisipan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(participants).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/participants - Register participant
func (h *ParticipantHandler) RegisterParticipant(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id" binding:"required"`
		StudentID string `json:"student_id" binding:"required"` // academic.students.id
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if session exists
	session, err := h.repo.GetAssessmentSessionByID(req.SessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi ujian tidak ditemukan").WithContext(c))
		return
	}

	// Check if already registered
	existing, _ := h.repo.GetParticipantBySessionAndStudent(req.SessionID, req.StudentID)
	if existing != nil {
		c.JSON(http.StatusConflict, sharederr.Error("ALREADY_REGISTERED", "Mahasiswa sudah terdaftar").WithContext(c))
		return
	}

	participant := domain.Participant{
		SessionID: req.SessionID,
		StudentID:  req.StudentID,
		Status:    "registered",
	}

	if err := h.repo.CreateParticipant(&participant); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mendaftarkan partisipan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(participant).WithContext(c))
}

// GET /api/v1/participants/:id - Get participant
func (h *ParticipantHandler) GetParticipant(c *gin.Context) {
	participantID := c.Param("id")

	participant, err := h.repo.GetParticipantByID(participantID)
	if err != nil || participant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Partisipan tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(participant).WithContext(c))
}

// POST /api/v1/participants/bulk - Bulk register participants
func (h *ParticipantHandler) BulkRegisterParticipants(c *gin.Context) {
	var req struct {
		SessionID string   `json:"session_id" binding:"required"`
		StudentIDs []string `json:"student_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if session exists
	session, err := h.repo.GetAssessmentSessionByID(req.SessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi ujian tidak ditemukan").WithContext(c))
		return
	}

	registered := 0
	skipped := 0

	for _, studentID := range req.StudentIDs {
		existing, _ := h.repo.GetParticipantBySessionAndStudent(req.SessionID, studentID)
		if existing != nil {
			skipped++
			continue
		}

		participant := domain.Participant{
			SessionID: req.SessionID,
			StudentID: studentID,
			Status:   "registered",
		}

		if err := h.repo.CreateParticipant(&participant); err == nil {
			registered++
		}
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(map[string]interface{}{
		"registered": registered,
		"skipped":    skipped,
	}, "Pendaftaran massal selesai").WithContext(c))
}

// DELETE /api/v1/participants/:id - Remove participant
func (h *ParticipantHandler) RemoveParticipant(c *gin.Context) {
	participantID := c.Param("id")

	if err := h.repo.DeleteParticipant(participantID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus partisipan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Partisipan berhasil dihapus").WithContext(c))
}

// GET /api/v1/participants/check/:session_id/:student_id - Check registration
func (h *ParticipantHandler) CheckRegistration(c *gin.Context) {
	sessionID := c.Param("session_id")
	studentID := c.Param("student_id")

	participant, err := h.repo.GetParticipantBySessionAndStudent(sessionID, studentID)
	if err != nil || participant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Partisipan tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(participant).WithContext(c))
}
