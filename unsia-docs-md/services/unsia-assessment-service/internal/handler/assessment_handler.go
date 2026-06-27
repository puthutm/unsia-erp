package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedhttpclient "github.com/unsia-erp/shared-httpclient"
	"github.com/unsia-erp/unsia-assessment-service/internal/domain"
	"github.com/unsia-erp/unsia-assessment-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type SessionCreateRequest struct {
	SessionType   string  `json:"session_type" binding:"required,oneof=cbt mid_exam final_exam"`
	Title         string  `json:"title" binding:"required"`
	ContextModule *string `json:"context_module"`
	ContextID     *string `json:"context_id"`
}

type AttemptCreateRequest struct {
	AssessmentSessionID string `json:"assessment_session_id" binding:"required"`
	ParticipantID       string `json:"participant_id" binding:"required"`
}

type ResultPublishRequest struct {
	AttemptID string  `json:"attempt_id" binding:"required"`
	Score     float64 `json:"score" binding:"required"`
}

type ParticipantRegisterRequest struct {
	AssessmentSessionID string  `json:"assessment_session_id" binding:"required"`
	ParticipantType     string  `json:"participant_type" binding:"required,oneof=applicant student"`
	ApplicantID         *string `json:"applicant_id"`
	StudentID           *string `json:"student_id"`
	UserID              *string `json:"user_id"`
}

type QuestionBankCreateRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	ModuleScope string `json:"module_scope"`
}

type QuestionOptionRequest struct {
	OptionLabel string `json:"option_label" binding:"required"`
	OptionText  string `json:"option_text" binding:"required"`
	IsCorrect   bool   `json:"is_correct"`
	SortOrder   int    `json:"sort_order"`
}

type QuestionCreateRequest struct {
	QuestionBankID    string                  `json:"question_bank_id" binding:"required"`
	QuestionType      string                  `json:"question_type" binding:"required"`
	Difficulty        string                  `json:"difficulty"`
	QuestionText      string                  `json:"question_text" binding:"required"`
	AnswerExplanation string                  `json:"answer_explanation"`
	Options           []QuestionOptionRequest `json:"options"`
}

type AnswerSaveRequest struct {
	QuestionID       string  `json:"question_id" binding:"required"`
	SelectedOptionID *string `json:"selected_option_id"`
	AnswerText       string  `json:"answer_text"`
}

type AssessmentHandler struct {
	repo           *repository.AssessmentRepository
	db             *gorm.DB
	pmbClient      *sharedhttpclient.Client
	academicClient *sharedhttpclient.Client
}

func NewAssessmentHandler(db *gorm.DB) *AssessmentHandler {
	pmbURL := os.Getenv("PMB_SERVICE_URL")
	if pmbURL == "" {
		pmbURL = "http://localhost:8004"
	}
	academicURL := os.Getenv("ACADEMIC_SERVICE_URL")
	if academicURL == "" {
		academicURL = "http://localhost:8006"
	}
	srvToken := os.Getenv("ASSESSMENT_SERVICE_TOKEN")
	if srvToken == "" {
		srvToken = "assessment_service_secret_token"
	}

	pmbClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      pmbURL,
		ServiceToken: srvToken,
		SourceName:   "assessment-service",
		Timeout:      10 * time.Second,
	})

	academicClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      academicURL,
		ServiceToken: srvToken,
		SourceName:   "assessment-service",
		Timeout:      10 * time.Second,
	})

	return &AssessmentHandler{
		repo:           repository.NewAssessmentRepository(db),
		db:             db,
		pmbClient:      pmbClient,
		academicClient: academicClient,
	}
}

func (h *AssessmentHandler) CreateSession(c *gin.Context) {
	var req SessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	session := domain.AssessmentSession{
		SessionType:   req.SessionType,
		Title:         req.Title,
		ContextModule: req.ContextModule,
		ContextID:     req.ContextID,
		Status:        "active",
	}

	if err := h.repo.CreateSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan sesi ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(session).WithContext(c))
}

func (h *AssessmentHandler) CreateAttempt(c *gin.Context) {
	var req AttemptCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	attempt := domain.AssessmentAttempt{
		AssessmentSessionID: req.AssessmentSessionID,
		ParticipantID:       req.ParticipantID,
		Status:              "started",
		StartedAt:           time.Now(),
	}

	if err := h.repo.CreateAttempt(&attempt); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat attempt ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(attempt).WithContext(c))
}

func (h *AssessmentHandler) RegisterParticipant(c *gin.Context) {
	var req ParticipantRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	session, err := h.repo.GetSessionByID(req.AssessmentSessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi ujian tidak ditemukan").WithContext(c))
		return
	}

	part := domain.AssessmentParticipant{
		AssessmentSessionID: req.AssessmentSessionID,
		ParticipantType:     req.ParticipantType,
		ApplicantID:         req.ApplicantID,
		StudentID:           req.StudentID,
		UserID:              req.UserID,
		Status:              "registered",
	}

	if err := h.repo.RegisterParticipant(&part); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mendaftarkan peserta").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(part).WithContext(c))
}

func (h *AssessmentHandler) ListSessions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	offset := 0

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}
	if val, err := strconv.Atoi(offsetStr); err == nil {
		offset = val
	}

	list, total, err := h.repo.ListSessions(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data sesi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"sessions": list,
		"total":    total,
	}).WithContext(c))
}

func (h *AssessmentHandler) ListParticipants(c *gin.Context) {
	sessionID := c.Query("assessment_session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("VALIDATION_ERROR", "assessment_session_id query parameter is required").WithContext(c))
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	offset := 0

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}
	if val, err := strconv.Atoi(offsetStr); err == nil {
		offset = val
	}

	list, total, err := h.repo.ListParticipantsBySessionID(sessionID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data peserta").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"participants": list,
		"total":        total,
	}).WithContext(c))
}

func (h *AssessmentHandler) ListAttempts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	offset := 0

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}
	if val, err := strconv.Atoi(offsetStr); err == nil {
		offset = val
	}

	list, total, err := h.repo.ListAttempts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data attempts").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"attempts": list,
		"total":    total,
	}).WithContext(c))
}


func (h *AssessmentHandler) PublishResult(c *gin.Context) {
	var req ResultPublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	attempt, err := h.repo.GetAttemptByID(req.AttemptID)
	if err != nil || attempt == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Attempt tidak ditemukan").WithContext(c))
		return
	}

	session, _ := h.repo.GetSessionByID(attempt.AssessmentSessionID)
	part, _ := h.repo.GetParticipantByID(attempt.ParticipantID)

	err = h.db.Transaction(func(tx *gorm.DB) error {
		attempt.TotalScore = &req.Score
		attempt.Status = "evaluated"
		tx.Model(attempt).Updates(map[string]interface{}{
			"total_score": &req.Score,
			"status":      "evaluated",
			"submitted_at": time.Now(),
		})

		// Notify modules depending on the context
		if session != nil && session.ContextModule != nil {
			ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

			if *session.ContextModule == "pmb" && part != nil && part.ApplicantID != nil {
				// CBT Seleksi PMB
				payload := map[string]interface{}{
					"score":         req.Score,
					"result_status": "pass", // simplify
				}
				url := fmt.Sprintf("/api/v1/pmb/applicants/%s/selection-results", *part.ApplicantID)
				resp, err := h.pmbClient.Post(ctx, url, payload)
				if err == nil {
					resp.Body.Close()
				}
			} else if *session.ContextModule == "lms" && part != nil && part.StudentID != nil && session.ContextID != nil {
				// LMS quiz
				payload := map[string]interface{}{
					"source_module":     "assessment",
					"academic_class_id": *session.ContextID,
					"student_id":        *part.StudentID,
					"score":             req.Score,
				}
				resp, err := h.academicClient.Post(ctx, "/api/v1/academic/grades/source-imports", payload)
				if err == nil {
					resp.Body.Close()
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mempublikasikan hasil ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Hasil ujian berhasil dipublikasikan").WithContext(c))
}

func (h *AssessmentHandler) CreateQuestionBank(c *gin.Context) {
	var req QuestionBankCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	qb := domain.QuestionBank{
		Code:        req.Code,
		Name:        req.Name,
		ModuleScope: req.ModuleScope,
		Status:      "active",
	}

	if err := h.repo.CreateQuestionBank(&qb); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan question bank").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(qb).WithContext(c))
}

func (h *AssessmentHandler) CreateQuestion(c *gin.Context) {
	var req QuestionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	diff := req.Difficulty
	if diff == "" {
		diff = "MEDIUM"
	}

	question := domain.Question{
		QuestionBankID:    req.QuestionBankID,
		QuestionType:      req.QuestionType,
		Difficulty:        diff,
		QuestionText:      req.QuestionText,
		AnswerExplanation: req.AnswerExplanation,
		Status:            "active",
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&question).Error; err != nil {
			return err
		}

		for _, optReq := range req.Options {
			opt := domain.QuestionOption{
				QuestionID:  question.ID,
				OptionLabel: optReq.OptionLabel,
				OptionText:  optReq.OptionText,
				IsCorrect:   optReq.IsCorrect,
				SortOrder:   optReq.SortOrder,
			}
			if err := tx.Create(&opt).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan soal").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(question).WithContext(c))
}

func (h *AssessmentHandler) CreateQuestionVersion(c *gin.Context) {
	questionID := c.Param("id")

	// Fetch Question and its options
	question, err := h.repo.GetQuestionByID(questionID)
	if err != nil || question == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Soal tidak ditemukan").WithContext(c))
		return
	}

	// Increment version number
	vNum, err := h.repo.GetLatestQuestionVersion(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data versi terbaru").WithContext(c))
		return
	}
	vNum++

	// JSON snapshot of options
	snapshotBytes, _ := json.Marshal(question.Options)
	snapshotStr := string(snapshotBytes)

	qv := domain.QuestionVersion{
		QuestionID:        questionID,
		VersionNumber:     vNum,
		QuestionType:      question.QuestionType,
		Difficulty:        question.Difficulty,
		QuestionText:      question.QuestionText,
		AnswerExplanation: question.AnswerExplanation,
		OptionsSnapshot:   snapshotStr,
		Status:            "active",
	}

	if err := h.repo.CreateQuestionVersion(&qv); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat versi baru soal").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(qv).WithContext(c))
}

func (h *AssessmentHandler) SaveAnswer(c *gin.Context) {
	attemptID := c.Param("id")
	var req AnswerSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Fetch Question Options to see which is correct
	var score float64 = 0.0
	if req.SelectedOptionID != nil {
		var option domain.QuestionOption
		err := h.db.Where("id = ? AND question_id = ?", *req.SelectedOptionID, req.QuestionID).First(&option).Error
		if err == nil && option.IsCorrect {
			score = 1.0 // Correct answer gets 1.0 point
		}
	}

	ans := domain.AssessmentAnswer{
		AttemptID:        attemptID,
		QuestionID:       req.QuestionID,
		SelectedOptionID: req.SelectedOptionID,
		AnswerText:       req.AnswerText,
		Score:            &score,
	}

	if err := h.repo.SaveAnswer(&ans); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jawaban").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(ans).WithContext(c))
}

func (h *AssessmentHandler) SubmitAttempt(c *gin.Context) {
	attemptID := c.Param("id")

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	attempt, err := h.repo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Attempt tidak ditemukan").WithContext(c))
		return
	}

	if attempt.Status == "submitted" || attempt.Status == "evaluated" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Attempt ini sudah disubmit").WithContext(c))
		return
	}

	// Fetch answers
	answers, err := h.repo.GetAnswersByAttemptID(attemptID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jawaban attempt").WithContext(c))
		return
	}

	// Calculate score
	var totalQuestions int64
	session, _ := h.repo.GetSessionByID(attempt.AssessmentSessionID)
	if session != nil && session.QuestionSetID != nil {
		h.db.Table("question_set_items").Where("question_set_id = ?", *session.QuestionSetID).Count(&totalQuestions)
	}

	if totalQuestions == 0 {
		totalQuestions = int64(len(answers))
	}

	var sumScores float64
	for _, ans := range answers {
		if ans.Score != nil {
			sumScores += *ans.Score
		}
	}

	var calculatedScore float64
	if totalQuestions > 0 {
		calculatedScore = (sumScores / float64(totalQuestions)) * 100
		if calculatedScore > 100 {
			calculatedScore = 100
		}
	}

	part, _ := h.repo.GetParticipantByID(attempt.ParticipantID)

	err = h.db.Transaction(func(tx *gorm.DB) error {
		attempt.TotalScore = &calculatedScore
		attempt.Status = "evaluated"
		now := time.Now()
		attempt.SubmittedAt = &now
		if err := tx.Model(attempt).Updates(map[string]interface{}{
			"total_score":  &calculatedScore,
			"status":       "evaluated",
			"submitted_at": &now,
		}).Error; err != nil {
			return err
		}

		// Notify modules depending on the context
		if session != nil && session.ContextModule != nil {
			ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

			if *session.ContextModule == "pmb" && part != nil && part.ApplicantID != nil {
				payload := map[string]interface{}{
					"score":         calculatedScore,
					"result_status": "pass",
				}
				url := fmt.Sprintf("/api/v1/pmb/applicants/%s/selection-results", *part.ApplicantID)
				resp, err := h.pmbClient.Post(ctx, url, payload)
				if err == nil {
					resp.Body.Close()
				}
			} else if *session.ContextModule == "lms" && part != nil && part.StudentID != nil && session.ContextID != nil {
				payload := map[string]interface{}{
					"source_module":     "assessment",
					"academic_class_id": *session.ContextID,
					"student_id":        *part.StudentID,
					"score":             calculatedScore,
				}
				resp, err := h.academicClient.Post(ctx, "/api/v1/academic/grades/source-imports", payload)
				if err == nil {
					resp.Body.Close()
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses submit ujian").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"attempt_id":  attempt.ID,
		"total_score": calculatedScore,
		"status":      "evaluated",
	}).WithContext(c))
}
