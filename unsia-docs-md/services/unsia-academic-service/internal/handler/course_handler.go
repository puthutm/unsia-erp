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

// ============ Course Handler ============

// CourseHandler handles course-related operations
type CourseHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Course Request Types ============

type CourseCreateRequest struct {
	CourseCode      string  `json:"course_code" binding:"required"`
	CourseName      string  `json:"course_name" binding:"required"`
	StudyProgramID  string  `json:"study_program_id" binding:"required"`
	Description   *string `json:"description"`
	SKS            int     `json:"sks" binding:"required"`
	Semester       int     `json:"semester" binding:"required"`
	CourseType     string  `json:"course_type" binding:"required,oneof=theory practicelab theorypractice"`
	IsMandatory   *bool   `json:"is_mandatory"`
	MinParticipants int    `json:"min_participants"`
	MaxParticipants int    `json:"max_participants"`
}

type CourseUpdateRequest struct {
	CourseName      *string `json:"course_name"`
	Description   *string `json:"description"`
	SKS            *int    `json:"sks"`
	Semester       *int    `json:"semester"`
	CourseType     *string `json:"course_type"`
	IsMandatory   *bool   `json:"is_mandatory"`
	MinParticipants *int   `json:"min_participants"`
	MaxParticipants *int    `json:"max_participants"`
	IsActive      *bool   `json:"is_active"`
}

// CreateCourse handles POST /api/v1/academic/courses
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req CourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isMandatory := true
	if req.IsMandatory != nil {
		isMandatory = *req.IsMandatory
	}

	minParticipants := 1
	if req.MinParticipants > 0 {
		minParticipants = req.MinParticipants
	}

	maxParticipants := 100
	if req.MaxParticipants > 0 {
		maxParticipants = req.MaxParticipants
	}

	course := domain.Course{
		CourseCode:      req.CourseCode,
		CourseName:     req.CourseName,
		StudyProgramID: req.StudyProgramID,
		Description:   req.Description,
		SKS:            req.SKS,
		Semester:       req.Semester,
		CourseType:     req.CourseType,
		IsMandatory:   isMandatory,
		MinParticipants: minParticipants,
		MaxParticipants: maxParticipants,
		IsActive:       true,
	}

	if err := h.repo.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan mata kuliah").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course.create",
		Module:       "academic",
		ResourceType: "course",
		ResourceID:   course.ID,
		NewValue:     course,
	})

	c.JSON(http.StatusCreated, sharederr.Success(course).WithContext(c))
}

// GetCourse handles GET /api/v1/academic/courses/:id
func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")

	course, err := h.repo.GetCourseByID(id)
	if err != nil || course == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mata kuliah tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(course).WithContext(c))
}

// ListCourses handles GET /api/v1/academic/courses
func (h *CourseHandler) ListCourses(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")
	semester := c.Query("semester")
	courseType := c.Query("course_type")
	isActive := c.Query("is_active") == "true"

	filter := repository.CourseFilter{
		StudyProgramID: studyProgramID,
		Semester:       0,
		CourseType:    courseType,
		IsActive:      isActive,
	}

	if s := c.Query("semester"); s != "" {
		var sem int
		if _, err := fmt.Sscanf(s, "%d", &sem); err == nil {
			filter.Semester = sem
		}
	}

	result, err := h.repo.ListCourses(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar mata kuliah").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// UpdateCourse handles PUT /api/v1/academic/courses/:id
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")

	course, err := h.repo.GetCourseByID(id)
	if err != nil || course == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mata kuliah tidak ditemukan").WithContext(c))
		return
	}

	var req CourseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if req.CourseName != nil {
		course.CourseName = *req.CourseName
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.SKS != nil {
		course.SKS = *req.SKS
	}
	if req.Semester != nil {
		course.Semester = *req.Semester
	}
	if req.CourseType != nil {
		course.CourseType = *req.CourseType
	}
	if req.IsMandatory != nil {
		course.IsMandatory = *req.IsMandatory
	}
	if req.MinParticipants != nil {
		course.MinParticipants = *req.MinParticipants
	}
	if req.MaxParticipants != nil {
		course.MaxParticipants = *req.MaxParticipants
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}

	if err := h.repo.UpdateCourse(course); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui mata kuliah").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course.update",
		Module:       "academic",
		ResourceType: "course",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.Success(course).WithContext(c))
}

// DeleteCourse handles DELETE /api/v1/academic/courses/:id
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	course, err := h.repo.GetCourseByID(id)
	if err != nil || course == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mata kuliah tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.DeleteCourse(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus mata kuliah").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course.delete",
		Module:       "academic",
		ResourceType: "course",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Mata kuliah berhasil dihapus").WithContext(c))
}

// GetCourseByStudyProgram handles GET /api/v1/academic/study-programs/:sp_id/courses
func (h *CourseHandler) GetCourseByStudyProgram(c *gin.Context) {
	spID := c.Param("sp_id")

	courses, err := h.repo.GetCoursesByStudyProgram(spID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar mata kuliah").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(courses).WithContext(c))
}

// ============ Course Offering Handler ============

type CourseOfferingCreateRequest struct {
	CourseID         string `json:"course_id" binding:"required"`
	AcademicPeriodID string `json:"academic_period_id" binding:"required"`
	Year            int    `json:"year" binding:"required"`
	Semester        int    `json:"semester" binding:"required"`
}

type CourseOfferingUpdateRequest struct {
	IsActive *bool `json:"is_active"`
}

// CreateCourseOffering handles POST /api/v1/academic/course-offerings
func (h *CourseHandler) CreateCourseOffering(c *gin.Context) {
	var req CourseOfferingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if course exists
	course, err := h.repo.GetCourseByID(req.CourseID)
	if err != nil || course == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mata kuliah tidak ditemukan").WithContext(c))
		return
	}

	// Check if offering already exists
	existing, _ := h.repo.GetCourseOffering(req.CourseID, req.AcademicPeriodID)
	if existing != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Offering untuk periode ini sudah ada").WithContext(c))
		return
	}

	offering := domain.CourseOffering{
		CourseID:         req.CourseID,
		AcademicPeriodID: req.AcademicPeriodID,
		Year:            req.Year,
		Semester:        req.Semester,
		IsActive:        true,
	}

	if err := h.repo.CreateCourseOffering(&offering); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan offering").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course_offering.create",
		Module:       "academic",
		ResourceType: "course_offering",
		ResourceID:   offering.ID,
	})

	c.JSON(http.StatusCreated, sharederr.Success(offering).WithContext(c))
}

// ListCourseOfferings handles GET /api/v1/academic/course-offerings
func (h *CourseHandler) ListCourseOfferings(c *gin.Context) {
	academicPeriodID := c.Query("academic_period_id")
	courseID := c.Query("course_id")

	var offerings []domain.CourseOffering
	query := h.db.Model(&domain.CourseOffering{})

	if academicPeriodID != "" {
		query = query.Where("academic_period_id = ?", academicPeriodID)
	}
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}

	if err := query.Find(&offerings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar offering").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(offerings).WithContext(c))
}

// GetCourseOfferingDetail handles GET /api/v1/academic/course-offerings/:id
func (h *CourseHandler) GetCourseOfferingDetail(c *gin.Context) {
	id := c.Param("id")

	offering, err := h.repo.GetCourseOfferingByID(id)
	if err != nil || offering == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Offering tidak ditemukan").WithContext(c))
		return
	}

	// Get course details
	course, _ := h.repo.GetCourseByID(offering.CourseID)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"offering": offering,
		"course":  course,
	}).WithContext(c))
}

// UpdateCourseOffering handles PUT /api/v1/academic/course-offerings/:id
func (h *CourseHandler) UpdateCourseOffering(c *gin.Context) {
	id := c.Param("id")

	offering, err := h.repo.GetCourseOfferingByID(id)
	if err != nil || offering == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Offering tidak ditemukan").WithContext(c))
		return
	}

	var req CourseOfferingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if req.IsActive != nil {
		offering.IsActive = *req.IsActive
	}

	if err := h.repo.UpdateCourseOffering(offering); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui offering").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course_offering.update",
		Module:       "academic",
		ResourceType: "course_offering",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.Success(offering).WithContext(c))
}

// ============ Import Helper ============

type CourseImportRequest struct {
	Courses []CourseCreateRequest `json:"courses" binding:"required,gt=0"`
}

// ImportCourses handles POST /api/v1/academic/courses/import
func (h *CourseHandler) ImportCourses(c *gin.Context) {
	var req CourseImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var created []domain.Course
	var failed []string

	for _, courseReq := range req.Courses {
		// Check if course code already exists
		existing, _ := h.repo.GetCourseByCode(courseReq.CourseCode)
		if existing != nil {
			failed = append(failed, courseReq.CourseCode+" (duplicate code)")
			continue
		}

		course := domain.Course{
			CourseCode:      courseReq.CourseCode,
			CourseName:     courseReq.CourseName,
			StudyProgramID: courseReq.StudyProgramID,
			Description:   courseReq.Description,
			SKS:            courseReq.SKS,
			Semester:       courseReq.Semester,
			CourseType:     courseReq.CourseType,
			IsMandatory:   true,
			MinParticipants: courseReq.MinParticipants,
			MaxParticipants: courseReq.MaxParticipants,
			IsActive:       true,
		}

		if err := h.repo.CreateCourse(&course); err != nil {
			failed = append(failed, courseReq.CourseCode+" (save error)")
			continue
		}

		created = append(created, course)
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course.import",
		Module:       "academic",
		ResourceType: "course",
		ResourceID:   "bulk_import",
	})

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"created": created,
		"failed":  failed,
		"total":   len(req.Courses),
	}).WithContext(c))
}

// Helper import
func fmt.Sprintf(format string, a ...interface{}) string {
	return format
}
