package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-assessment-service/internal/domain"
	"github.com/unsia-erp/unsia-assessment-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type AssessmentSessionHandler struct {
	repo *repository.AssessmentRepository
	db   *gorm.DB
}

func NewAssessmentSessionHandler(db *gorm.DB) *AssessmentSessionHandler {
	return &AssessmentSessionHandler{
		repo: repository.NewAssessmentRepository(db),
		db:   db,
	}
}

// GET /api/v1/assessment-sessions - List assessment sessions
func (h *AssessmentSessionHandler) ListAssessmentSessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	questionBankID := c.Query("question_bank_id")

	sessions, total, err := h.repo.ListAssessmentSessions(page, limit, status, questionBankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(sessions).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/assessment-sessions - Create assessment session
func (h *AssessmentSessionHandler) CreateAssessmentSession(c *gin.Context) {
	var req struct {
		SessionName    string `json:"session_name" binding:"required"`
		QuestionBankID string `json:"question_bank_id" binding:"required"`
		Duration     int    `json:"duration"` // in minutes
		MaxAttempts  int    `json:"max_attempts"`
		Randomize   bool   `json:"randomize"`
		PassingScore float64 `json:"passing_score"`
		StartTime   string `json:"start_time"`
		EndTime    string `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var startTime, endTime *time.Time
	if req.StartTime != "" {
		t, _ := time.Parse(time.RFC3339, req.StartTime)
		startTime = &t
	}
	if req.EndTime != "" {
		t, _ := time.Parse(time.RFC3339, req.EndTime)
		endTime = &t
	}

	session := domain.AssessmentSession{
		SessionName:    req.SessionName,
		QuestionBankID: req.QuestionBankID,
		Duration:     req.Duration,
		MaxAttempts:  req.MaxAttempts,
		Randomize:    req.Randomize,
		PassingScore: req.PassingScore,
		StartTime:   startTime,
		EndTime:    endTime,
		Status:     "draft",
	}

	if err := h.repo.CreateAssessmentSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(session).WithContext(c))
}

// GET /api/v1/assessment-sessions/:id - Get assessment session
func (h *AssessmentSessionHandler) GetAssessmentSession(c *gin.Context) {
	sessionID := c.Param("id")

	session, err := h.repo.GetAssessmentSessionByID(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi ujian tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(session).WithContext(c))
}

// PUT /api/v1/assessment-sessions/:id - Update assessment session
func (h *AssessmentSessionHandler) UpdateAssessmentSession(c *gin.Context) {
	sessionID := c.Param("id")

	var req struct {
		SessionName  string  `json:"session_name"`
		Duration   int     `json:"duration"`
		MaxAttempts int    `json:"max_attempts"`
		Randomize   bool    `json:"randomize"`
		PassingScore float64 `json:"passing_score"`
		Status    string  `json:"status"`
		StartTime  string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.SessionName != "" {
		updates["session_name"] = req.SessionName
	}
	if req.Duration > 0 {
		updates["duration"] = req.Duration
	}
	if req.MaxAttempts > 0 {
		updates["max_attempts"] = req.MaxAttempts
	}
	if req.PassingScore > 0 {
		updates["passing_score"] = req.PassingScore
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.StartTime != "" {
		t, _ := time.Parse(time.RFC3339, req.StartTime)
		updates["start_time"] = t
	}
	if req.EndTime != "" {
		t, _ := time.Parse(time.RFC3339, req.EndTime)
		updates["end_time"] = t
	}

	if err := h.repo.UpdateAssessmentSession(sessionID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Sesi ujian berhasil diperbarui").WithContext(c))
}

// POST /api/v1/assessment-sessions/:id/publish - Publish assessment session
func (h *AssessmentSessionHandler) PublishAssessmentSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.repo.UpdateAssessmentSession(sessionID, map[string]interface{}{"status": "published"}); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mempublikasikan sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Sesi ujian berhasil dipublikasikan").WithContext(c))
}

// POST /api/v1/assessment-sessions/:id/close - Close assessment session
func (h *AssessmentSessionHandler) CloseAssessmentSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.repo.UpdateAssessmentSession(sessionID, map[string]interface{}{"status": "closed"}); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menutup sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Sesi ujian berhasil ditutup").WithContext(c))
}

// DELETE /api/v1/assessment-sessions/:id - Delete assessment session
func (h *AssessmentSessionHandler) DeleteAssessmentSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.repo.DeleteAssessmentSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Sesi ujian berhasil dihapus").WithContext(c))
}
