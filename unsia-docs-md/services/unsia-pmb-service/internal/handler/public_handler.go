package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-pmb-service/internal/domain"
	"github.com/unsia-erp/unsia-pmb-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Public Registration Handlers ============
// For PMB PUBLIK.html - Public registration page

type PublicRegistrationRequest struct {
	FullName           string `json:"full_name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	StudyProgramID     string `json:"study_program_id" binding:"required"`
	PmbWaveID         string `json:"pmb_wave_id" binding:"required"`
	AdmissionPathID   string `json:"admission_path_id"`
	TargetPeriodID    string `json:"target_period_id"`
}

type PublicVerifyRequest struct {
	OTP string `json:"otp" binding:"required"`
}

type PublicResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PublicHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewPublicHandler(db *gorm.DB) *PublicHandler {
	return &PublicHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

func generateOTP() string {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", nBig.Int64())
}

func generateRegNumber() string {
	now := time.Now().Format("2006")
	nBig, _ := rand.Int(rand.Reader, big.NewInt(90000))
	num := nBig.Int64() + 10000
	return fmt.Sprintf("PMB%s%d", now, num)
}

// CheckAvailability handles GET /api/v1/public/pmb/availability
// Check open waves and available study programs
func (h *PublicHandler) CheckAvailability(c *gin.Context) {
	waves, err := h.repo.GetActiveWaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data").WithContext(c))
		return
	}

	programs, err := h.repo.GetAvailableStudyPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil program studi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"waves":         waves,
		"study_programs": programs,
	}).WithContext(c))
}

// GetWaveDetails handles GET /api/v1/public/pmb/waves/:id
func (h *PublicHandler) GetWaveDetails(c *gin.Context) {
	id := c.Param("id")

	wave, err := h.repo.GetWaveByID(id)
	if err != nil || wave == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Gelombang tidak ditemukan").WithContext(c))
		return
	}

	// Get study programs for this wave
	programs, _ := h.repo.GetStudyPrograms(repository.StudyProgramFilter{
		IsActive: true,
		Page:   1,
		Limit:  100,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"wave":          wave,
		"study_programs": programs,
	}).WithContext(c))
}

// InitiateRegistration handles POST /api/v1/public/pmb/register
// Step 1: Public user starts registration
func (h *PublicHandler) InitiateRegistration(c *gin.Context) {
	var req PublicRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if wave is active
	wave, err := h.repo.GetWaveByID(req.PmbWaveID)
	if err != nil || wave == nil || !wave.IsActive {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Gelombang tidak tersedia").WithContext(c))
		return
	}

	// Check if study program exists
	sp, err := h.repo.GetStudyProgramByID(req.StudyProgramID)
	if err != nil || sp == nil || !sp.IsActive {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Program studi tidak tersedia").WithContext(c))
		return
	}

	// Generate OTP
	otp := generateOTP()
	expiresAt := time.Now().Add(10 * 60 * time.Second) // 10 minutes

	// Check if email already registered
	existingOTP, _ := h.repo.GetPendingRegistration(req.Email)
	if existingOTP != nil {
		// Update existing OTP
		h.db.Model(existingOTP).Updates(map[string]interface{}{
			"otp":        otp,
			"expires_at": expiresAt,
		})
	} else {
		// Create new pending registration
		pending := domain.PendingRegistration{
			Email:           req.Email,
			Phone:          req.Phone,
			StudyProgramID:  req.StudyProgramID,
			PmbWaveID:      req.PmbWaveID,
			AdmissionPathID: req.AdmissionPathID,
			TargetPeriodID: req.TargetPeriodID,
			OTP:           otp,
			ExpiresAt:     &expiresAt,
			Status:        "pending",
		}
		h.db.Create(&pending)
	}

	// TODO: Send OTP via email/SMS
	// For now, we'll return it (in production, don't return OTP)
	_ = otp // Remove this in production

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message":    "Kode verifikasi telah dikirim ke email Anda",
		"expires_in": 600, // seconds
	}).WithContext(c))
}

// VerifyOTPAndCreateApplicant handles POST /api/v1/public/pmb/verify
// Step 2: Verify OTP and create applicant
func (h *PublicHandler) VerifyOTPAndCreateApplicant(c *gin.Context) {
	var req PublicVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	email := c.GetHeader("X-Registration-Email")
	if email == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Email diperlukan").WithContext(c))
		return
	}

	// Find pending registration
	pending, err := h.repo.GetPendingRegistration(email)
	if err != nil || pending == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	// Verify OTP
	if pending.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Kode verifikasi tidak valid").WithContext(c))
		return
	}

	if time.Now().After(*pending.ExpiresAt) {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Kode verifikasi telah kadaluarsa").WithContext(c))
		return
	}

	// Generate registration number
	regNumber := generateRegNumber()

	// Create applicant
	applicant := domain.Applicant{
		PersonID:            pending.Email, // Will be linked to person later
		StudyProgramID:      &pending.StudyProgramID,
		PmbWaveID:          &pending.PmbWaveID,
		AdmissionPathID:     &pending.AdmissionPathID,
		TargetEntryPeriodID: &pending.TargetPeriodID,
		RegistrationNumber: regNumber,
		Status:             "draft",
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Create applicant
		if err := tx.Create(&applicant).Error; err != nil {
			return err
		}

		// Create biodata
		biodata := domain.ApplicantBiodata{
			ApplicantID:    applicant.ID,
			FullName:       pending.Email,
			Phone:         pending.Phone,
			CoreSyncStatus: "pending",
		}
		if err := tx.Create(&biodata).Error; err != nil {
			return err
		}

		// Update pending status
		tx.Model(pending).Update("status", "completed")

		// Publish event
		envelope := sharedevent.EventEnvelope{
			EventName:        "pmb.applicant_created",
			EventVersion:    "v1",
			PublisherService: "pmb-service",
			AggregateType:  "applicant",
			AggregateID:    applicant.ID,
			Payload: map[string]interface{}{
				"applicant_id":         applicant.ID,
				"registration_number":  applicant.RegistrationNumber,
				"email":              email,
				"study_program_id":   applicant.StudyProgramID,
				"status":            applicant.Status,
			},
		}

		conn := tx.Statement.ConnPool
		_, err = sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat pendaftaran").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"applicant_id":         applicant.ID,
		"registration_number": applicant.RegistrationNumber,
		"message":            "Pendaftaran berhasil! Silakan lengkapi data Anda.",
	}).WithContext(c))
}

// ResendOTP handles POST /api/v1/public/pmb/resend-otp
func (h *PublicHandler) ResendOTP(c *gin.Context) {
	var req PublicResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	pending, err := h.repo.GetPendingRegistration(req.Email)
	if err != nil || pending == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Email tidak ditemukan").WithContext(c))
		return
	}

	// Generate new OTP
	otp := generateOTP()
	expiresAt := time.Now().Add(10 * 60 * time.Second)

	h.db.Model(pending).Updates(map[string]interface{}{
		"otp":         otp,
		"expires_at":  expiresAt,
		"attempts":   0,
	})

	// TODO: Send OTP
	_ = otp

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Kode verifikasi baru telah dikirim",
	}).WithContext(c))
}

// GetRegistrationStatus handles GET /api/v1/public/pmb/status/:registration_number
func (h *PublicHandler) GetRegistrationStatus(c *gin.Context) {
	regNumber := c.Param("registration_number")

	applicant, err := h.repo.GetApplicantByRegNumber(regNumber)
	if err != nil || applicant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	// Get biodata
	biodata, _ := h.repo.GetBiodata(applicant.ID)

	// Get documents count
	docs, _ := h.repo.GetApplicantDocuments(applicant.ID, "")

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"registration_number": applicant.RegistrationNumber,
		"status":            applicant.Status,
		"biodata":          biodata,
		"documents_count":  len(docs),
	}).WithContext(c))
}

// GetPublicWaveList handles GET /api/v1/public/pmb/waves
func (h *PublicHandler) GetPublicWaveList(c *gin.Context) {
	waves, err := h.repo.GetActiveWaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil gelombang").WithContext(c))
		return
	}

	// Filter only active waves with registration still open
	var openWaves []domain.PmbWave
	now := time.Now()
	for _, w := range waves {
		if w.IsActive && now.After(w.StartDate) && now.Before(w.EndDate) {
			openWaves = append(openWaves, w)
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(openWaves).WithContext(c))
}

// GetPublicStudyPrograms handles GET /api/v1/public/pmb/study-programs
func (h *PublicHandler) GetPublicStudyPrograms(c *gin.Context) {
	waveID := c.Query("wave_id")

	programs, err := h.repo.GetAvailableStudyPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil program studi").WithContext(c))
		return
	}

	_ = waveID // Could filter by wave

	c.JSON(http.StatusOK, sharederr.Success(programs).WithContext(c))
}
