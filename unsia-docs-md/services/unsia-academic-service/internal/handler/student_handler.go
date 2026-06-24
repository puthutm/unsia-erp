package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Student Request Types ============

type StudentGenerateRequest struct {
	ApplicantID         string  `json:"applicant_id" binding:"required"`
	CurriculumID        *string `json:"curriculum_id"`
	EntryAcademicYearID *string `json:"entry_academic_year_id"`
	EntryPeriodID       *string `json:"entry_period_id"`
	StudyProgramID      string  `json:"study_program_id" binding:"required"`
	Reason              string  `json:"reason"`
}

type StudentUpdateRequest struct {
	StudyProgramID      *string `json:"study_program_id"`
	StudentStatus     *string `json:"student_status"`
	CurrentSemester   *int    `json:"current_semester"`
	CurriculumID      *string `json:"curriculum_id"`
}

type AssignAdvisorRequest struct {
	LecturerID string `json:"lecturer_id" binding:"required"` // HRIS lecturer ID
}

type StudentPromotionRequest struct {
	StudyProgramID      string `json:"study_program_id" binding:"required"`
	EntryAcademicYearID string `json:"entry_academic_year_id" binding:"required"`
	EntryPeriodID     string `json:"entry_period_id" binding:"required"`
}

// StudentHandler handles student-related operations
type StudentHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewStudentHandler creates a new StudentHandler
func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// GenerateStudentFromApplicant converts an applicant to a student
// POST /api/v1/academic/students/generate
func (h *StudentHandler) GenerateStudentFromApplicant(c *gin.Context) {
	var req StudentGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	// Verify if student already exists for this applicant
	existing, err := h.repo.GetStudentByApplicantID(req.ApplicantID)
	if err == nil && existing != nil {
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"student_id":      existing.ID,
			"nim":             existing.Nim,
			"handover_status": "success",
		}).WithContext(c))
		return
	}

	var student domain.Student
	err = h.db.Transaction(func(tx *gorm.DB) error {
		yearStr := time.Now().Format("2006")
		periodID := "default-period-id"
		if req.EntryPeriodID != nil {
			periodID = *req.EntryPeriodID
		}

		nim, err := h.repo.GenerateNIM(tx, req.StudyProgramID, periodID, yearStr)
		if err != nil {
			return err
		}

		// Person ID matches applicant's person ID
		personID := "person-generated-id"

		student = domain.Student{
			PersonID:            personID,
			ApplicantID:         &req.ApplicantID,
			StudyProgramID:      req.StudyProgramID,
			Nim:                 nim,
			StudentStatus:       "active",
			EntryAcademicYearID: req.EntryAcademicYearID,
			EntryPeriodID:       req.EntryPeriodID,
			CurriculumID:        req.CurriculumID,
			CurrentSemester:     1,
		}

		if err := tx.Create(&student).Error; err != nil {
			return err
		}

		envelope := sharedevent.EventEnvelope{
			EventName:        "academic.student_created",
			EventVersion:     "v1",
			PublisherService: "academic-service",
			AggregateType:    "student",
			AggregateID:      student.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"student_id":       student.ID,
				"person_id":        student.PersonID,
				"nim":              student.Nim,
				"study_program_id": student.StudyProgramID,
				"status":           student.StudentStatus,
			},
		}

		conn := tx.Statement.ConnPool
		_, err = sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal melakukan generate mahasiswa").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.student.create",
		Module:       "academic",
		ResourceType: "student",
		ResourceID:   student.ID,
		NewValue:     student,
	})

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"student_id":      student.ID,
		"nim":             student.Nim,
		"handover_status": "success",
	}).WithContext(c))
}

// ListStudents retrieves list of students
// GET /api/v1/academic/students
func (h *StudentHandler) ListStudents(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")
	academicPeriodID := c.Query("academic_period_id")
	status := c.Query("status")
	search := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	students, total, err := h.repo.ListStudents(studyProgramID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data mahasiswa").WithContext(c))
		return
	}

	// Apply additional filters
	if academicPeriodID != "" {
		// Filter by entry period
		h.db.Where("entry_period_id = ?", academicPeriodID).Find(&students)
	}
	if status != "" {
		// Filter by status
		var filtered []domain.Student
		for _, s := range students {
			if s.StudentStatus == status {
				filtered = append(filtered, s)
			}
		}
		students = filtered
	}
	if search != "" {
		// Search by NIM or name
		var filtered []domain.Student
		for _, s := range students {
			if string(s.Nim[0]) == search || contains(s.Nim, search) {
				filtered = append(filtered, s)
			}
		}
		students = filtered
	}

	c.JSON(http.StatusOK, sharederr.Success(students).WithPagination(page, limit, int(total)).WithContext(c))
}

// GetStudentDetail retrieves student details
// GET /api/v1/academic/students/:id
func (h *StudentHandler) GetStudentDetail(c *gin.Context) {
	studentID := c.Param("id")

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(student).WithContext(c))
}

// UpdateStudent updates student information
// PUT /api/v1/academic/students/:id
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	studentID := c.Param("id")
	var req StudentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.StudyProgramID != nil {
		updates["study_program_id"] = *req.StudyProgramID
	}
	if req.StudentStatus != nil {
		updates["student_status"] = *req.StudentStatus
	}
	if req.CurrentSemester != nil {
		updates["current_semester"] = *req.CurrentSemester
	}
	if req.CurriculumID != nil {
		updates["curriculum_id"] = *req.CurriculumID
	}

	if err := h.repo.UpdateStudent(studentID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui data mahasiswa").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.student.update",
		Module:       "academic",
		ResourceType: "student",
		ResourceID:   studentID,
		OldValue:     student,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Data mahasiswa berhasil diperbarui").WithContext(c))
}

// PromoteStudent promotes student to next semester/year
// POST /api/v1/academic/students/:id/promote
func (h *StudentHandler) PromoteStudent(c *gin.Context) {
	studentID := c.Param("id")
	var req StudentPromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	newSemester := student.CurrentSemester + 1
	entryYear := student.EntryAcademicYearID
	entryPeriod := student.EntryPeriodID

	if req.EntryAcademicYearID != "" {
		entryYear = &req.EntryAcademicYearID
	}
	if req.EntryPeriodID != "" {
		entryPeriod = &req.EntryPeriodID
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"current_semester":     newSemester,
			"entry_academic_year_id": entryYear,
			"entry_period_id":      entryPeriod,
			"updated_at":           time.Now(),
		}
		return tx.Model(student).Updates(updates).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mempromosikan mahasiswa").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.student.promote",
		Module:       "academic",
		ResourceType: "student",
		ResourceID:   studentID,
		OldValue:     student,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":       studentID,
		"new_semester":    newSemester,
		"promotion_status": "success",
	}).WithContext(c))
}

// GetStudentKrs retrieves student's KRS
// GET /api/v1/academic/students/:id/krs
func (h *StudentHandler) GetStudentKrs(c *gin.Context) {
	studentID := c.Param("id")
	academicPeriodID := c.Query("academic_period_id")
	status := c.Query("status")

	krsList, err := h.repo.GetKrsByStudentID(studentID, academicPeriodID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data KRS").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(krsList).WithContext(c))
}

// GetStudentGrades retrieves student's grades
// GET /api/v1/academic/students/:id/grades
func (h *StudentHandler) GetStudentGrades(c *gin.Context) {
	studentID := c.Param("id")
	academicPeriodID := c.Query("academic_period_id")

	grades, err := h.repo.GetGradesByStudentID(studentID, academicPeriodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data nilai").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(grades).WithContext(c))
}

// AssignAdvisor assigns a PA (Pembimbing Akademik) to a student
// POST /api/v1/academic/students/:id/advisor
func (h *StudentHandler) AssignAdvisor(c *gin.Context) {
	studentID := c.Param("id")
	var req AssignAdvisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	advisorID := req.LecturerID
	err = h.repo.UpdateStudentAdvisor(studentID, &advisorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menetapkan PA").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.student.advisor.assign",
		Module:       "academic",
		ResourceType: "student",
		ResourceID:   studentID,
		OldValue:     student,
		NewValue:     advisorID,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":     studentID,
		"advisor_id":    advisorID,
		"advisor_status": "assigned",
	}).WithContext(c))
}

// RemoveAdvisor removes PA from a student
// DELETE /api/v1/academic/students/:id/advisor
func (h *StudentHandler) RemoveAdvisor(c *gin.Context) {
	studentID := c.Param("id")

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Pass nil to clear advisor
	err = h.repo.UpdateStudentAdvisor(studentID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus PA").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.student.advisor.remove",
		Module:       "academic",
		ResourceType: "student",
		ResourceID:   studentID,
		OldValue:     student,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":     studentID,
		"advisor_status": "removed",
	}).WithContext(c))
}

// GetStudentsByAdvisor retrieves all students assigned to a specific PA
// GET /api/v1/academic/advisors/:lecturer_id/students
func (h *StudentHandler) GetStudentsByAdvisor(c *gin.Context) {
	lecturerID := c.Param("lecturer_id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	students, total, err := h.repo.GetStudentsByAdvisor(lecturerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data mahasiswa").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(students).WithPagination(page, limit, int(total)).WithContext(c))
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
