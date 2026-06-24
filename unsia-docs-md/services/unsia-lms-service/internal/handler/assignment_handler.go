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

type AssignmentCreateRequest struct {
	SessionID   string `json:"session_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Instruction string `json:"instruction"`
	DueAt       string `json:"due_at"` // RFC3339 format
}

type SubmissionCreateRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	FileURL   string `json:"file_url" binding:"required"`
}

type AssignmentHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

func NewAssignmentHandler(db *gorm.DB) *AssignmentHandler {
	return &AssignmentHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// CreateAssignment - Membuat tugas/quiz baru
func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var req AssignmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var dueAt *time.Time
	if req.DueAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueAt)
		if err == nil {
			dueAt = &parsed
		}
	}

	assignment := domain.Assignment{
		SessionID:   req.SessionID,
		Title:       req.Title,
		Instruction: req.Instruction,
		DueAt:       dueAt,
		Status:      "active",
	}

	if err := h.repo.CreateAssignment(&assignment); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat tugas").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(assignment).WithContext(c))
}

// CreateSubmission - Mengupload jawaban tugas
func (h *AssignmentHandler) CreateSubmission(c *gin.Context) {
	assignmentID := c.Param("id")
	var req SubmissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	submission := domain.AssignmentSubmission{
		AssignmentID: assignmentID,
		StudentID:    req.StudentID,
		FileURL:      req.FileURL,
		SubmittedAt:  time.Now(),
	}

	if err := h.repo.CreateAssignmentSubmission(&submission); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupload jawaban tugas").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(submission).WithContext(c))
}
