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

type AttemptHandler struct {
	repo *repository.AssessmentRepository
	db   *gorm.DB
}

func NewAttemptHandler(db *gorm.DB) *AttemptHandler {
	return &AttemptHandler{
		repo: repository.NewAssessmentRepository(db),
		db:   db,
	}
}

// GET /api/v1/attempts - List attempts
func (h *AttemptHandler) ListAttempts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sessionID := c.Query("session_id")
	studentID := c.Query("student_id")
	status := c.Query("status")

	attempts, total, err := h.repo.ListAttempts(page, limit, sessionID, studentID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data attempt").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(attempts).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/attempts - Start attempt
func (h *AttemptHandler) StartAttempt(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id" binding:"required"`
		StudentID string `json:"student_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if participant exists
	participant, err := h.repo.GetParticipantBySessionAndStudent(req.SessionID, req.StudentID)
	if err != nil || participant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Partisipan tidak terdaftar").WithContext(c))
		return
	}

	// Check max attempts
	session, _ := h.repo.GetAssessmentSessionByID(req.SessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi tidak ditemukan").WithContext(c))
		return
	}

	count, _ := h.repo.CountAttemptsBySessionAndStudent(req.SessionID, req.StudentID)
	if int(count) >= session.MaxAttempts {
		c.JSON(http.StatusForbidden, sharederr.Error("MAX_ATTEMPTS", "Menggunakan maximal attempt").WithContext(c))
		return
	}

	// Get questions for session
	questions, err := h.repo.GetQuestionsByQuestionBank(session.QuestionBankID)
	if err != nil || len(questions) == 0 {
		c.JSON(http.StatusBadRequest, sharederr.Error("NO_QUESTIONS", "Tidak ada soal untuk sesi ini").WithContext(c))
		return
	}

	// Randomize if needed
	if session.Randomize {
		// Shuffle questions (implementation depends on your needs)
	}

	now := time.Now()
	attempt := domain.Attempt{
		SessionID:  req.SessionID,
		StudentID:  req.StudentID,
		StartTime: &now,
		Status:    "in_progress",
		Score:     0,
	}

	if err := h.repo.CreateAttempt(&attempt); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memulai attempt").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(attempt).WithContext(c))
}

// GET /api/v1/attempts/:id - Get attempt
func (h *AttemptHandler) GetAttempt(c *gin.Context) {
	attemptID := c.Param("id")

	attempt, err := h.repo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Attempt tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(attempt).WithContext(c))
}

// POST /api/v1/attempts/:id/submit - Submit attempt
func (h *AttemptHandler) SubmitAttempt(c *gin.Context) {
	attemptID := c.Param("id")

	var req struct {
		Answers []struct {
			QuestionID string `json:"question_id"`
			Answer    string `json:"answer"`
		} `json:"answers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get attempt
	attempt, err := h.repo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Attempt tidak ditemukan").WithContext(c))
		return
	}

	if attempt.Status != "in_progress" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATUS", "Attempt sudah selesai").WithContext(c))
		return
	}

	// Get session to calculate score
	session, _ := h.repo.GetAssessmentSessionByID(attempt.SessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi tidak ditemukan").WithContext(c))
		return
	}

	// Calculate score
	correctAnswers := 0
	totalQuestions := len(req.Answers)

	for _, answer := range req.Answers {
		// Check if answer is correct (simplified - you need full implementation)
		isCorrect, _ := h.repo.CheckAnswer(answer.QuestionID, answer.Answer)
		if isCorrect {
			correctAnswers++
		}
	}

	score := float64(correctAnswers) / float64(totalQuestions) * 100
	now := time.Now()

	status := "completed"
	if score >= session.PassingScore {
		status = "passed"
	}

	// Update attempt
	updates := map[string]interface{}{
		"score":     score,
		"status":    status,
		"end_time": now,
	}

	if err := h.repo.UpdateAttempt(attemptID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jawaban").WithContext(c))
		return
	}

	// Optionally sync grade to Academic service
	// ...

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(map[string]interface{}{
		"score":            score,
		"correct_answers": correctAnswers,
		"total_questions": totalQuestions,
		"status":         status,
	}, "Attempt berhasil disubmit").WithContext(c))
}

// GET /api/v1/attempts/:id/questions - Get questions for attempt
func (h *AttemptHandler) GetAttemptQuestions(c *gin.Context) {
	attemptID := c.Param("id")

	attempt, err := h.repo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Attempt tidak ditemukan").WithContext(c))
		return
	}

	if attempt.Status != "in_progress" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATUS", "Attempt sudah selesai").WithContext(c))
		return
	}

	// Get session questions
	session, _ := h.repo.GetAssessmentSessionByID(attempt.SessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi tidak ditemukan").WithContext(c))
		return
	}

	questions, err := h.repo.GetQuestionsByQuestionBank(session.QuestionBankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil soal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(questions).WithContext(c))
}

// POST /api/v1/attempts/:id/answers - Save answer (checkpoint)
func (h *AttemptHandler) SaveAnswer(c *gin.Context) {
	attemptID := c.Param("id")

	var req struct {
		QuestionID string `json:"question_id" binding:"required"`
		Answer    string `json:"answer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check attempt status
	attempt, err := h.repo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Attempt tidak ditemukan").WithContext(c))
		return
	}

	if attempt.Status != "in_progress" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATUS", "Attempt sudah selesai").WithContext(c))
		return
	}

	// Save answer
	answer := domain.Answer{
		AttemptID:  attemptID,
		QuestionID: req.QuestionID,
		Answer:    req.Answer,
	}

	if err := h.repo.SaveAnswer(&answer); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jawaban").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Jawaban disimpan").WithContext(c))
}

// GET /api/v1/attempts/:id/answers - Get saved answers
func (h *AttemptHandler) GetAnswers(c *gin.Context) {
	attemptID := c.Param("id")

	answers, err := h.repo.GetAnswersByAttemptID(attemptID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jawaban").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(answers).WithContext(c))
}
