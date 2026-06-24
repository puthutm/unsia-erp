package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedhttpclient "github.com/unsia-erp/shared-httpclient"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type ClassSyncRequest struct {
	AcademicClassID string   `json:"academic_class_id" binding:"required"`
	LecturerIDs     []string `json:"lecturer_ids"`
}

type EnrollmentSyncRequest struct {
	AcademicClassID  string `json:"academic_class_id" binding:"required"`
	StudentID        string `json:"student_id" binding:"required"`
	EnrollmentStatus string `json:"enrollment_status" binding:"required"` // active, inactive
}

type GradeSyncRequest struct {
	AcademicClassID string  `json:"academic_class_id" binding:"required"`
	StudentID       string  `json:"student_id" binding:"required"`
	Score           float64 `json:"score" binding:"required"`
}

type SessionCreateRequest struct {
	SessionNumber int    `json:"session_number" binding:"required"`
	Title         string `json:"title" binding:"required"`
	SessionDate   string `json:"session_date"` // YYYY-MM-DD
	StartTime     string `json:"start_time"`   // HH:MM
	EndTime       string `json:"end_time"`     // HH:MM
}

type MaterialCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	FileURL     string `json:"file_url" binding:"required"`
}

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

type AttendanceCreateRequest struct {
	StudentID        string `json:"student_id" binding:"required"`
	AttendanceStatus string `json:"attendance_status" binding:"required"`
}

type LMSHandler struct {
	repo           *repository.LMSRepository
	db             *gorm.DB
	academicClient *sharedhttpclient.Client
}

func NewLMSHandler(db *gorm.DB) *LMSHandler {
	academicURL := os.Getenv("ACADEMIC_SERVICE_URL")
	if academicURL == "" {
		academicURL = "http://localhost:8006"
	}
	srvToken := os.Getenv("LMS_SERVICE_TOKEN")
	if srvToken == "" {
		srvToken = "lms_service_secret_token"
	}

	academicClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      academicURL,
		ServiceToken: srvToken,
		SourceName:   "lms-service",
		Timeout:      10 * time.Second,
	})

	return &LMSHandler{
		repo:           repository.NewLMSRepository(db),
		db:             db,
		academicClient: academicClient,
	}
}

func (h *LMSHandler) SyncClassFromAcademic(c *gin.Context) {
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

func (h *LMSHandler) SyncEnrollmentFromKrs(c *gin.Context) {
	var req EnrollmentSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	class, err := h.repo.GetClassByAcademicID(req.AcademicClassID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengecek data kelas").WithContext(c))
		return
	}
	if class == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kelas LMS terkait tidak ditemukan, lakukan sync kelas terlebih dahulu").WithContext(c))
		return
	}

	enrollment, err := h.repo.SyncEnrollment(class.ID, req.StudentID, req.EnrollmentStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal melakukan sync enrollment").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(enrollment).WithContext(c))
}

func (h *LMSHandler) SyncLmsGrade(c *gin.Context) {
	var req GradeSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

	// Post grade source import to academic service
	payload := map[string]interface{}{
		"source_module":     "lms",
		"academic_class_id": req.AcademicClassID,
		"student_id":        req.StudentID,
		"score":             req.Score,
	}

	resp, err := h.academicClient.Post(ctx, "/api/v1/academic/grades/source-imports", payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, sharederr.Error("ACADEMIC_SERVICE_UNAVAILABLE", fmt.Sprintf("Gagal menghubungi layanan Akademik: %v", err)).WithContext(c))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusBadGateway, sharederr.Error("ACADEMIC_SERVICE_ERROR", fmt.Sprintf("Layanan Akademik mengembalikan error code: %d", resp.StatusCode)).WithContext(c))
		return
	}

	var academicResp struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&academicResp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("PARSE_ERROR", "Gagal membaca response Akademik").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(academicResp.Data, "Nilai berhasil disinkronisasikan ke Akademik").WithContext(c))
}

func (h *LMSHandler) CreateSession(c *gin.Context) {
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

func (h *LMSHandler) CreateMaterial(c *gin.Context) {
	sessionID := c.Param("id")
	var req MaterialCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	now := time.Now()
	material := domain.Material{
		SessionID:   sessionID,
		Title:       req.Title,
		ContentType: req.ContentType,
		FileURL:     req.FileURL,
		PublishedAt: &now,
	}

	if err := h.repo.CreateMaterial(&material); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat materi perkuliahan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(material).WithContext(c))
}

func (h *LMSHandler) CreateAssignment(c *gin.Context) {
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

func (h *LMSHandler) CreateSubmission(c *gin.Context) {
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

func (h *LMSHandler) CreateAttendance(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengisi kehadiran").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(attendance).WithContext(c))
}
