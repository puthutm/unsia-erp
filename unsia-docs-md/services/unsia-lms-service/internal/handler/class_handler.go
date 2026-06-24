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

type ClassHandler struct {
	repo *repository.LmsRepository
	db   *gorm.DB
}

func NewClassHandler(db *gorm.DB) *ClassHandler {
	return &ClassHandler{
		repo: repository.NewLmsRepository(db),
		db:   db,
	}
}

// GET /api/v1/lms/classes - List classes
func (h *ClassHandler) ListClasses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	courseID := c.Query("course_id")
	lecturerID := c.Query("lecturer_id")
	status := c.Query("status")

	classes, total, err := h.repo.ListClasses(page, limit, courseID, lecturerID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kelas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(classes).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/lms/classes - Create class
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req struct {
		AcademicClassID string  `json:"academic_class_id" binding:"required"`
		CourseID     string  `json:"course_id" binding:"required"`
		LecturerID   *string `json:"lecturer_id"`
		ClassCode   string  `json:"class_code"`
		Semester    string  `json:"semester"`
		AcademicYear string  `json:"academic_year"`
		MaxStudents int    `json:"max_students"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	class := domain.Class{
		AcademicClassID: req.AcademicClassID,
		CourseID:       req.CourseID,
		LecturerID:     req.LecturerID,
		ClassCode:      req.ClassCode,
		Semester:       req.Semester,
		AcademicYear:  req.AcademicYear,
		MaxStudents:   req.MaxStudents,
		Status:       "active",
	}

	if err := h.repo.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat kelas").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(class).WithContext(c))
}

// GET /api/v1/lms/classes/:id - Get class
func (h *ClassHandler) GetClass(c *gin.Context) {
	classID := c.Param("id")

	class, err := h.repo.GetClassByID(classID)
	if err != nil || class == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kelas tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(class).WithContext(c))
}

// PUT /api/v1/lms/classes/:id - Update class
func (h *ClassHandler) UpdateClass(c *gin.Context) {
	classID := c.Param("id")

	var req struct {
		LecturerID   *string `json:"lecturer_id"`
		ClassCode   string  `json:"class_code"`
		Status     string  `json:"status"`
		MaxStudents int    `json:"max_students"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.LecturerID != nil {
		updates["lecturer_id"] = *req.LecturerID
	}
	if req.ClassCode != "" {
		updates["class_code"] = req.ClassCode
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.MaxStudents > 0 {
		updates["max_students"] = req.MaxStudents
	}

	if err := h.repo.UpdateClass(classID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui kelas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Kelas berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/lms/classes/:id - Delete class
func (h *ClassHandler) DeleteClass(c *gin.Context) {
	classID := c.Param("id")

	if err := h.repo.DeleteClass(classID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus kelas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Kelas berhasil dihapus").WithContext(c))
}
