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

// ============ Advisor/PA Request Types ============

type AssignAdvisorRequest struct {
	AdvisorID     string `json:"advisor_id" binding:"required"`
	StudentIDs    []string `json:"student_ids" binding:"required,gt=0"`
	AcademicYear string `json:"academic_year" binding:"required"`
	Semester      int    `json:"semester" binding:"required"`
}

type UnassignAdvisorRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

// ============ AdvisorHandler ============
// Handler untuk Dosen PA (Pembimbing Akademik)

type AdvisorHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewAdvisorHandler creates a new AdvisorHandler
func NewAdvisorHandler(db *gorm.DB) *AdvisorHandler {
	return &AdvisorHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ PA Assignment APIs ============

// AssignStudentsToAdvisor assigns students to an advisor (PA)
// POST /api/v1/academic/advisors/assign
func (h *AdvisorHandler) AssignStudentsToAdvisor(c *gin.Context) {
	var req AssignAdvisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify advisor exists
	advisor, err := h.repo.GetLecturerByID(req.AdvisorID)
	if err != nil || advisor == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Dosen tidak ditemukan").WithContext(c))
		return
	}

	var assigned []domain.StudentAdvisor
	var failed []string

	for _, studentID := range req.StudentIDs {
		// Verify student exists
		student, err := h.repo.GetStudentByID(studentID)
		if err != nil || student == nil {
			failed = append(failed, studentID+" (mahasiswa tidak ditemukan)")
			continue
		}

		// Check if already assigned
		existing, _ := h.repo.GetStudentAdvisor(studentID, req.AcademicYear, req.Semester)
		if existing != nil {
			failed = append(failed, studentID+" (sudah memiliki PA)")
			continue
		}

		sa := domain.StudentAdvisor{
			StudentID:    studentID,
			AdvisorID:    req.AdvisorID,
			AcademicYear: req.AcademicYear,
			Semester:     req.Semester,
			AssignedAt:   time.Now(),
		}

		if err := h.repo.CreateStudentAdvisor(&sa); err != nil {
			failed = append(failed, studentID+" (gagal assign)")
			continue
		}

		assigned = append(assigned, sa)
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.advisor.assign",
		Module:       "academic",
		ResourceType: "student_advisor",
		ResourceID:   req.AdvisorID,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"assigned": assigned,
		"failed":   failed,
		"total":    len(req.StudentIDs),
	}).WithContext(c))
}

// GetMyAdvisees gets students assigned to current advisor
// GET /api/v1/academic/advisors/my-students
func (h *AdvisorHandler) GetMyAdvisees(c *gin.Context) {
	// Get advisor ID from token/context
	advisorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Unauthorized").WithContext(c))
		return
	}

	advisorIDStr, _ := advisorID.(string)
	academicYear := c.Query("academic_year")
	semester := c.Query("semester")

	var advisees []domain.StudentAdvisor
	query := h.db.Model(&domain.StudentAdvisor{}).Where("advisor_id = ?", advisorIDStr)

	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}

	if err := query.Find(&advisees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar mahasiswa bimbingan").WithContext(c))
		return
	}

	// Get student details
	var studentIDs []string
	for _, adv := range advisees {
		studentIDs = append(studentIDs, adv.StudentID)
	}

	var students []domain.Student
	if len(studentIDs) > 0 {
		h.db.Where("id IN ?", studentIDs).Find(&students)
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"advisees": advisees,
		"students": students,
	}).WithContext(c))
}

// GetPendingKrsForAdvisor gets KRS pending approval for advisor's advisees
// GET /api/v1/academic/advisors/my-students/pending-krs
func (h *AdvisorHandler) GetPendingKrsForAdvisor(c *gin.Context) {
	advisorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Unauthorized").WithContext(c))
		return
	}

	advisorIDStr, _ := advisorID.(string)

	// Get advisee student IDs
	var advisees []domain.StudentAdvisor
	if err := h.db.Where("advisor_id = ?", advisorIDStr).Find(&advisees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data").WithContext(c))
		return
	}

	var studentIDs []string
	for _, adv := range advisees {
		studentIDs = append(studentIDs, adv.StudentID)
	}

	if len(studentIDs) == 0 {
		c.JSON(http.StatusOK, sharederr.Success([]interface{}{}).WithContext(c))
		return
	}

	// Get KRS with status submitted
	var krsList []domain.KRS
	h.db.Where("student_id IN ? AND status = ?", studentIDs, "submitted").Find(&krsList)

	// Enrich with student info
	var result []gin.H
	for _, krs := range krsList {
		student, _ := h.repo.GetStudentByID(krs.StudentID)
		result = append(result, gin.H{
			"krs_id":              krs.ID,
			"student_id":         krs.StudentID,
			"student_name":       student.Name,
			"student_nim":        student.NIM,
			"academic_period_id": krs.AcademicPeriodID,
			"status":             krs.Status,
			"submitted_at":       krs.SubmittedAt,
		})
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// GetAdvisorList lists all advisors
// GET /api/v1/academic/advisors
func (h *AdvisorHandler) GetAdvisorList(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")

	var advisors []domain.Lecturer
	query := h.db.Model(&domain.Lecturer{})

	if studyProgramID != "" {
		query = query.Where("study_program_id = ?", studyProgramID)
	}

	if err := query.Find(&advisors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar dosen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(advisors).WithContext(c))
}

// GetAdvisorStudents gets students for a specific advisor
// GET /api/v1/academic/advisors/:id/students
func (h *AdvisorHandler) GetAdvisorStudents(c *gin.Context) {
	advisorID := c.Param("id")
	academicYear := c.Query("academic_year")
	semester := c.Query("semester")

	var studentAdvisors []domain.StudentAdvisor
	query := h.db.Model(&domain.StudentAdvisor{}).Where("advisor_id = ?", advisorID)

	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}

	if err := query.Find(&studentAdvisors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data").WithContext(c))
		return
	}

	// Get student details
	var studentIDs []string
	for _, sa := range studentAdvisors {
		studentIDs = append(studentIDs, sa.StudentID)
	}

	var students []domain.Student
	if len(studentIDs) > 0 {
		h.db.Where("id IN ?", studentIDs).Find(&students)
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"advisor_id": advisorID,
		"students":  students,
	}).WithContext(c))
}

// UnassignStudentFromAdvisor removes student from advisor
// POST /api/v1/academic/advisors/unassign
func (h *AdvisorHandler) UnassignStudentFromAdvisor(c *gin.Context) {
	var req UnassignAdvisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if err := h.db.Where("student_id = ?", req.StudentID).Delete(&domain.StudentAdvisor{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus hubungan PA").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.advisor.unassign",
		Module:       "academic",
		ResourceType: "student_advisor",
		ResourceID:   req.StudentID,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Mahasiswa berhasil dilepas dari PA").WithContext(c))
}

// GetAdvisorStats gets advisor statistics
// GET /api/v1/academic/advisors/:id/stats
func (h *AdvisorHandler) GetAdvisorStats(c *gin.Context) {
	advisorID := c.Param("id")
	academicYear := c.Query("academic_year")
	semester := c.Query("semester")

	// Count total advisees
	var totalCount int64
	query := h.db.Model(&domain.StudentAdvisor{}).Where("advisor_id = ?", advisorID)
	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	query.Count(&totalCount)

	// Count pending KRS
	var pendingKrsCount int64
	h.db.Model(&domain.KRS{}).
		Joins("JOIN student_advisors ON student_advisors.student_id = krs.student_id").
		Where("student_advisors.advisor_id = ?", advisorID).
		Where("krs.status = ?", "submitted").
		Count(&pendingKrsCount)

	// Count approved KRS
	var approvedKrsCount int64
	h.db.Model(&domain.KRS{}).
		Joins("JOIN student_advisors ON student_advisors.student_id = krs.student_id").
		Where("student_advisors.advisor_id = ?", advisorID).
		Where("krs.status = ?", "approved").
		Count(&approvedKrsCount)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"advisor_id":         advisorID,
		"total_advisees":    totalCount,
		"pending_krs":       pendingKrsCount,
		"approved_krs":      approvedKrsCount,
		"academic_year":    academicYear,
		"semester":         semester,
	}).WithContext(c))
}
