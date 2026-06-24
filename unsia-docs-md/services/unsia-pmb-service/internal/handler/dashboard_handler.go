package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-pmb-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ CAMABA Dashboard Handler ============
// For DASHBOARD PENDAFTARAN.html - Applicant Dashboard

type DashboardHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

// GetDashboard handles GET /api/v1/camaba/dashboard
// Get applicant's dashboard data
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Get applicant ID from token claims
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	// Get applicant
	applicant, err := h.repo.GetApplicantByID(applicantID)
	if err != nil || applicant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	// Get biodata
	biodata, _ := h.repo.GetBiodata(applicantID)

	// Get documents
	docs, _ := h.repo.GetApplicantDocuments(applicantID, "")

	// Get selection schedule
	schedules, _ := h.repo.GetSchedules(applicantID, "")

	// Get payment status
	payments, _ := h.repo.GetPayments(applicantID)

	// Count documents by status
	docVerified := 0
	docPending := 0
	for _, d := range docs {
		if d.VerificationStatus == "approved" {
			docVerified++
		} else {
			docPending++
		}
	}

	// Check registration status progress
	step := 1 // Registration started
	if biodata != nil && biodata.FullName != "" {
		step = 2 // Biodata completed
	}
	if docVerified > 0 {
		step = 3 // Documents uploaded
	}
	if len(schedules) > 0 {
		step = 4 // Selection scheduled
	}

	// Find next payment
	var nextPayment map[string]interface{}
	for _, p := range payments {
		if p.PaymentStatus == "pending" {
			nextPayment = map[string]interface{}{
				"amount":    p.Amount,
				"due_date": p.DueDate,
				"type":    p.PaymentType,
			}
			break
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"registration_number": applicant.RegistrationNumber,
		"status":             applicant.Status,
		"step":              step,
		"biodata":           biodata,
		"documents": map[string]interface{}{
			"total":      len(docs),
			"verified":   docVerified,
			"pending":    docPending,
		},
		"schedules":       schedules,
		"next_payment":    nextPayment,
		"selection_date": applicant.SelectionDate,
		"selection_venue": applicant.SelectionVenue,
	}).WithContext(c))
}

// GetMyDocuments handles GET /api/v1/camaba/documents
func (h *DashboardHandler) GetMyDocuments(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	docs, err := h.repo.GetApplicantDocuments(applicantID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(docs).WithContext(c))
}

// UploadDocument handles POST /api/v1/camaba/documents
func (h *DashboardHandler) UploadDocument(c *gin.Context) {
	var req struct {
		DocumentType string `json:"document_type" binding:"required"`
		FileURL     string `json:"file_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	// Check if document type already uploaded
	existing, _ := h.repo.GetDocumentByType(applicantID, req.DocumentType)
	if existing != nil {
		// Update existing
		h.db.Model(existing).Updates(map[string]interface{}{
			"file_url":            req.FileURL,
			"verification_status": "pending",
			"rejection_reason":    nil,
		})
		c.JSON(http.StatusOK, sharederr.Success(existing).WithContext(c))
		return
	}

	// Get document type info
	docType, _ := h.repo.GetDocumentType(req.DocumentType)

	// Create new document
	doc := map[string]interface{}{
		"applicant_id":          applicantID,
		"document_type":        req.DocumentType,
		"document_type_name":   docType.Name,
		"file_url":            req.FileURL,
		"verification_status": "pending",
	}

	result := h.db.Create(&doc)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal upload dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(doc).WithContext(c))
}

// GetSelectionSchedule handles GET /api/v1/camaba/schedule
func (h *DashboardHandler) GetSelectionSchedule(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	schedules, err := h.repo.GetSchedules(applicantID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// GetPayments handles GET /api/v1/camaba/payments
func (h *DashboardHandler) GetPayments(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	payments, err := h.repo.GetPayments(applicantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil tagihan").WithContext(c))
		return
	}

	// Calculate total paid and pending
	var totalPaid, totalPending float64
	for _, p := range payments {
		if p.PaymentStatus == "paid" || p.PaymentStatus == "success" {
			totalPaid += p.Amount
		} else if p.PaymentStatus == "pending" || p.PaymentStatus == "unpaid" {
			totalPending += p.Amount
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"payments":   payments,
		"total_paid":   totalPaid,
		"total_pending": totalPending,
	}).WithContext(c))
}

// GetProfile handles GET /api/v1/camaba/profile
func (h *DashboardHandler) GetProfile(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	applicant, err := h.repo.GetApplicantByID(applicantID)
	if err != nil || applicant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	biodata, _ := h.repo.GetBiodata(applicantID)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"applicant": applicant,
		"biodata":  biodata,
	}).WithContext(c))
}

// UpdateProfile handles PUT /api/v1/camaba/profile
func (h *DashboardHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		FullName        string `json:"full_name"`
		Gender         string `json:"gender"`
		BirthPlace     string `json:"birth_place"`
		BirthDate      string `json:"birth_date"`
		Address        string `json:"address"`
		ProvinceID     string `json:"province_id"`
		RegencyID      string `json:"regency_id"`
		DistrictID     string `json:"district_id"`
		PostalCode     string `json:"postal_code"`
		Phone          string `json:"phone"`
		ParentName     string `json:"parent_name"`
		ParentPhone    string `json:"parent_phone"`
		SchoolName     string `json:"school_name"`
		SchoolGradYear int    `json:"school_grad_year"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	// Parse birth date
	var birthDate *time.Time
	if req.BirthDate != "" {
		parsed, _ := time.Parse("2006-01-02", req.BirthDate)
		birthDate = &parsed
	}

	// Update biodata
	updates := map[string]interface{}{
		"full_name":      req.FullName,
		"gender":        req.Gender,
		"birth_place":   req.BirthPlace,
		"birth_date":    birthDate,
		"address":       req.Address,
		"province_id":   req.ProvinceID,
		"regency_id":    req.RegencyID,
		"district_id":   req.DistrictID,
		"postal_code":   req.PostalCode,
		"phone":        req.Phone,
		"parent_name":  req.ParentName,
		"parent_phone": req.ParentPhone,
		"school_name":  req.SchoolName,
		"school_grad_year": req.SchoolGradYear,
	}

	biodata, _ := h.repo.GetBiodata(applicantID)
	if biodata != nil {
		h.db.Model(biodata).Updates(updates)
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Biodata berhasil diperbarui",
	}).WithContext(c))
}

// GetAnnouncements handles GET /api/v1/camaba/announcements
func (h *DashboardHandler) GetAnnouncements(c *gin.Context) {
	announcements, err := h.repo.GetActiveAnnouncements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(announcements).WithContext(c))
}

// SubmitRegistration handles POST /api/v1/camaba/submit
func (h *DashboardHandler) SubmitRegistration(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	applicantID := claims.Subject

	// Get applicant
	applicant, err := h.repo.GetApplicantByID(applicantID)
	if err != nil || applicant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	// Check if biodata is complete
	biodata, _ := h.repo.GetBiodata(applicantID)
	if biodata == nil || biodata.FullName == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Silakan lengkapi biodata terlebih dahulu").WithContext(c))
		return
	}

	// Check if required documents are uploaded
	docs, _ := h.repo.GetApplicantDocuments(applicantID, "required")
	if len(docs) == 0 {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Silakan upload dokumen persyaratan terlebih dahulu").WithContext(c))
		return
	}

	// Check if all documents are verified
	allVerified := true
	for _, d := range docs {
		if d.VerificationStatus != "approved" {
			allVerified = false
			break
		}
	}

	// Update status
	newStatus := "submitted"
	if allVerified {
		newStatus = "verified"
	}

	h.db.Model(applicant).Updates(map[string]interface{}{
		"status": newStatus,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message":        "Pendaftaran berhasil disubmit",
		"status":       newStatus,
		"all_verified": allVerified,
	}).WithContext(c))
}
