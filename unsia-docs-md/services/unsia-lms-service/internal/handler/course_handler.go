package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type CourseHandler struct {
	repo *repository.LmsRepository
	db   *gorm.DB
}

func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{
		repo: repository.NewLmsRepository(db),
		db:   db,
	}
}

// GET /api/v1/lms/courses - List courses
func (h *CourseHandler) ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	courses, total, err := h.repo.ListCourses(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data mata kuliah").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(courses).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/lms/courses - Create course
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req struct {
		CourseID   string `json:"course_id" binding:"required"` // academic.courses.id
		CourseName string `json:"course_name" binding:"required"`
		Credit     int    `json:"credit"`
		Semester   int    `json:"semester"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	course := domain.Course{
		AcademicCourseID: req.CourseID,
		CourseName:     req.CourseName,
		Credit:         req.Credit,
		Semester:       req.Semester,
		Status:         "active",
	}

	if err := h.repo.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat mata kuliah").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(course).WithContext(c))
}

// GET /api/v1/lms/courses/:id - Get course
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	course, err := h.repo.GetCourseByID(courseID)
	if err != nil || course == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mata kuliah tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(course).WithContext(c))
}

// PUT /api/v1/lms/courses/:id - Update course
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")

	var req struct {
		CourseName string `json:"course_name"`
		Credit   int    `json:"credit"`
		Semester int    `json:"semester"`
		Status  string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.CourseName != "" {
		updates["course_name"] = req.CourseName
	}
	if req.Credit > 0 {
		updates["credit"] = req.Credit
	}
	if req.Semester > 0 {
		updates["semester"] = req.Semester
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.repo.UpdateCourse(courseID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui mata kuliah").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Mata kuliah berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/lms/courses/:id - Delete course
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")

	if err := h.repo.DeleteCourse(courseID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus mata kuliah").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Mata kuliah berhasil dihapus").WithContext(c))
}
