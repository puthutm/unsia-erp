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
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
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

func (h *LMSHandler) ListCourses(c *gin.Context) {
	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)
	ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

	url := fmt.Sprintf("/api/v1/academic/courses?%s", c.Request.URL.RawQuery)
	resp, err := h.academicClient.Get(ctx, url)
	if err != nil {
		c.JSON(http.StatusBadGateway, sharederr.Error("ACADEMIC_SERVICE_UNAVAILABLE", fmt.Sprintf("Gagal menghubungi layanan Akademik: %v", err)).WithContext(c))
		return
	}
	defer resp.Body.Close()

	var academicResp interface{}
	if err := json.NewDecoder(resp.Body).Decode(&academicResp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("PARSE_ERROR", "Gagal membaca response Akademik").WithContext(c))
		return
	}

	c.JSON(resp.StatusCode, academicResp)
}

func (h *LMSHandler) ListClasses(c *gin.Context) {
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

	list, total, err := h.repo.ListClasses(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kelas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"classes": list,
		"total":   total,
	}).WithContext(c))
}

func (h *LMSHandler) ListEnrollments(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("VALIDATION_ERROR", "class_id query parameter is required").WithContext(c))
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

	list, total, err := h.repo.ListEnrollmentsByClassID(classID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pendaftaran").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"enrollments": list,
		"total":       total,
	}).WithContext(c))
}

func (h *LMSHandler) ListSessions(c *gin.Context) {
	classID := c.Param("id")
	list, err := h.repo.ListSessionsByClassID(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data sesi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(list).WithContext(c))
}

func (h *LMSHandler) ListMaterials(c *gin.Context) {
	sessionID := c.Param("id")
	list, err := h.repo.ListMaterialsBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data materi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(list).WithContext(c))
}

func (h *LMSHandler) ListAssignments(c *gin.Context) {
	sessionID := c.Param("id")
	list, err := h.repo.ListAssignmentsBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data tugas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(list).WithContext(c))
}

func (h *LMSHandler) ListSubmissions(c *gin.Context) {
	assignmentID := c.Param("id")
	list, err := h.repo.ListSubmissionsByAssignmentID(assignmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data submission").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(list).WithContext(c))
}

func (h *LMSHandler) ListAttendance(c *gin.Context) {
	sessionID := c.Param("id")
	list, err := h.repo.ListAttendanceBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kehadiran").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(list).WithContext(c))
}

func (h *LMSHandler) ListAllSessions(c *gin.Context) {
	classID := c.Query("course_id")
	
	var list []domain.Session
	var err error
	
	if classID != "" {
		var class domain.Class
		if err := h.db.Where("id = ? OR academic_class_id = ?", classID, classID).First(&class).Error; err == nil {
			list, err = h.repo.ListSessionsByClassID(class.ID)
		} else {
			list = []domain.Session{}
		}
	} else {
		err = h.db.Order("session_number asc").Find(&list).Error
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data sesi").WithContext(c))
		return
	}
	
	// Map to frontend expected model structure
	type SessionResponse struct {
		ID              string `json:"id"`
		CourseID        string `json:"courseId"`
		CourseName      string `json:"courseName"`
		Title           string `json:"title"`
		Description     string `json:"description"`
		ScheduledAt     string `json:"scheduledAt"`
		Duration        int    `json:"duration"`
		Status          string `json:"status"`
		MaterialCount   int    `json:"materialCount"`
		AssignmentCount int    `json:"assignmentCount"`
	}
	
	responses := make([]SessionResponse, 0, len(list))
	for _, session := range list {
		// Get counts
		var matCount, assignCount int64
		h.db.Model(&domain.Material{}).Where("session_id = ?", session.ID).Count(&matCount)
		h.db.Model(&domain.Assignment{}).Where("session_id = ?", session.ID).Count(&assignCount)
		
		scheduledAt := ""
		if session.SessionDate != nil {
			scheduledAt = session.SessionDate.Format("2006-01-02T15:04:05Z")
		} else {
			scheduledAt = time.Now().Format("2006-01-02T15:04:05Z")
		}
		
		status := "upcoming"
		if session.Status == "active" {
			status = "ongoing"
		}
		
		responses = append(responses, SessionResponse{
			ID:              session.ID,
			CourseID:        session.LmsClassID,
			CourseName:      "Kelas " + session.LmsClassID,
			Title:           session.Title,
			Description:     session.Title,
			ScheduledAt:     scheduledAt,
			Duration:        90, // Default duration in minutes
			Status:          status,
			MaterialCount:   int(matCount),
			AssignmentCount: int(assignCount),
		})
	}
	
	c.JSON(http.StatusOK, sharederr.Success(responses).WithContext(c))
}

