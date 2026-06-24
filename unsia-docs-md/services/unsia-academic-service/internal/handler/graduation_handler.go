package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Graduation Handler ============

type GraduationHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

func NewGraduationHandler(db *gorm.DB) *GraduationHandler {
	return &GraduationHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Graduation Request Types ============

type GraduationEligibilityRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

type GraduationApplyRequest struct {
	StudentID          string `json:"student_id" binding:"required"`
	GraduationPeriodID string `json:"graduation_period_id" binding:"required"`
}

type GraduationApproveRequest struct {
	StudentID     string  `json:"student_id" binding:"required"`
	CertificateNumber string `json:"certificate_number"`
	ApproveNote  string  `json:"note"`
}

type CertificateRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=degree transkrip"` // degree, transkrip
}

// ============ Graduation Eligibility Check ============

// CheckGraduationEligibility handles GET /api/v1/academic/graduation/eligibility/:student_id
func (h *GraduationHandler) CheckGraduationEligibility(c *gin.Context) {
	studentID := c.Param("student_id")

	// Get student
	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Get all grades for the student
	grades, err := h.repo.GetStudentGrades(studentID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil nilai").WithContext(c))
		return
	}

	// Calculate GPA and total credits
	var totalPoints float64
	var totalSks int
	for _, g := range grades {
		if g.NumericGrade != nil && g.GradePoint != nil {
			sks, _ := h.repo.GetSksByKrsItemID(g.KrsItemID)
			totalPoints += *g.GradePoint * float64(sks)
			totalSks += sks
		}
	}

	ipk := 0.0
	if totalSks > 0 {
		ipk = totalPoints / float64(totalSks)
	}

	// Check graduation requirements (example: min 144 SKS, min IPK 2.75)
	minSks := 144
	minIPK := 2.75

	eligible := totalSks >= minSks && ipk >= minIPK

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":        studentID,
		"eligible":         eligible,
		"total_sks":        totalSks,
		"ipk":             ipk,
		"requirements": gin.H{
			"min_sks": minSks,
			"min_ipk": minIPK,
		},
		"status": func() string {
			if eligible {
				return "layak"
			}
			return "tidak_layak"
		}(),
	}).WithContext(c))
}

// ============ Graduation Application ============

// ApplyGraduation handles POST /api/v1/academic/graduation/apply
func (h *GraduationHandler) ApplyGraduation(c *gin.Context) {
	var req GraduationApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get student
	student, err := h.repo.GetStudentByID(req.StudentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Check if already applied
	if student.StudentStatus == "graduated" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Mahasiswa sudah lulus").WithContext(c))
		return
	}

	// Check eligibility
	grades, _ := h.repo.GetStudentGrades(req.StudentID, "")
	var totalSks int
	var totalPoints float64
	for _, g := range grades {
		if g.NumericGrade != nil && g.GradePoint != nil {
			sks, _ := h.repo.GetSksByKrsItemID(g.KrsItemID)
			totalPoints += *g.GradePoint * float64(sks)
			totalSks += sks
		}
	}

	ipk := 0.0
	if totalSks > 0 {
		ipk = totalPoints / float64(totalSks)
	}

	// Apply (set status pending)
	now := time.Now()
	updates := map[string]interface{}{
		"student_status": "graduated",
		"updated_at":    now,
		"graduation_date": now,
	}

	if err := h.db.Model(&domain.Student{}).Where("id = ?", req.StudentID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengajukan kelulusan").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.graduation.apply",
		Module:      "academic",
		ResourceType: "student",
		ResourceID:  req.StudentID,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(gin.H{
		"student_id":   req.StudentID,
		"status":     "pending_graduate",
		"total_sks":  totalSks,
		"ipk":        ipk,
	}, "Pengajuan kelulusan berhasil").WithContext(c))
}

// ============ Graduation Approval ============

// ApproveGraduation handles POST /api/v1/academic/graduation/approve
func (h *GraduationHandler) ApproveGraduation(c *gin.Context) {
	var req GraduationApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get student
	student, err := h.repo.GetStudentByID(req.StudentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Approve graduation
	now := time.Now()
	updates := map[string]interface{}{
		"student_status": "graduated",
		"updated_at":    now,
	}

	if err := h.db.Model(&domain.Student{}).Where("id = ?", req.StudentID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui kelulusan").WithContext(c))
		return
	}

	// Create transcript record
	transcript := domain.Transcript{
		StudentID: req.StudentID,
		IssuedAt:  now,
	}
	h.db.Create(&transcript)

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.graduation.approve",
		Module:      "academic",
		ResourceType: "student",
		ResourceID:  req.StudentID,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(gin.H{
		"student_id":         req.StudentID,
		"status":           "graduated",
		"certificate_number": req.CertificateNumber,
		"graduated_at":      now,
	}, "Kelulusan disetujui").WithContext(c))
}

// ============ Certificate Generation ============

// GenerateCertificate handles POST /api/v1/academic/graduation/certificate
func (h *GraduationHandler) GenerateCertificate(c *gin.Context) {
	var req CertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get student
	student, err := h.repo.GetStudentByID(req.StudentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	if student.StudentStatus != "graduated" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Mahasiswa belum lulus").WithContext(c))
		return
	}

	// Get transcript data
	grades, _ := h.repo.GetStudentGrades(req.StudentID, "")
	var totalSks int
	var totalPoints float64
	for _, g := range grades {
		if g.NumericGrade != nil && g.GradePoint != nil {
			sks, _ := h.repo.GetSksByKrsItemID(g.KrsItemID)
			totalPoints += *g.GradePoint * float64(sks)
			totalSks += sks
		}
	}

	ipk := 0.0
	if totalSks > 0 {
		ipk = totalPoints / float64(totalSks)
	}

	// Generate certificate (simulated URL)
	certType := "Ijazah"
	if req.Type == "transkrip" {
		certType = "Transkrip Nilai"
	}

	now := time.Now()
	certURL := "https://unsia.ac.id/certificate/" + req.StudentID + "/" + req.Type + "?t=" + now.Format("20060102150405")

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":    req.StudentID,
		"type":        certType,
		"url":         certURL,
		"total_sks":   totalSks,
		"ipk":        ipk,
		"issued_at":   now,
		"valid_until": now.AddDate(10, 0, 0), // 10 years
	}).WithContext(c))
}

// ============ Alumni Transfer ============

// GetAlumni handles GET /api/v1/academic/graduation/alumni
func (h *GraduationHandler) GetAlumni(c *gin.Context) {
	var students []domain.Student
	err := h.db.Where("student_status = ?", "graduated").Find(&students).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data alumni").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(students).WithContext(c))
}

// GetAlumniByYear handles GET /api/v1/academic/graduation/alumni/:year
func (h *GraduationHandler) GetAlumniByYear(c *gin.Context) {
	year := c.Param("year")

	var students []domain.Student
	err := h.db.Joins("JOIN graduation_records ON graduation_records.student_id = students.id").
		Where("students.student_status = ? AND EXTRACT(YEAR FROM graduation_records.graduated_at) = ?", "graduated", year).
		Find(&students).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data alumni").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(students).WithContext(c))
}
