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
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"gorm.io/gorm"
)

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

type EnrollmentHandler struct {
	repo           *repository.LMSRepository
	db             *gorm.DB
	academicClient *sharedhttpclient.Client
}

func NewEnrollmentHandler(db *gorm.DB) *EnrollmentHandler {
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

	return &EnrollmentHandler{
		repo:           repository.NewLMSRepository(db),
		db:             db,
		academicClient: academicClient,
	}
}

// SyncEnrollmentFromKrs - Sync enrollment dari KRS Academic
func (h *EnrollmentHandler) SyncEnrollmentFromKrs(c *gin.Context) {
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

// SyncLmsGrade - Sinkronisasi nilai ke Academic Service
func (h *EnrollmentHandler) SyncLmsGrade(c *gin.Context) {
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
