package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-pmb-service/internal/domain"
	"github.com/unsia-erp/unsia-pmb-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Study Program Handlers ============

type StudyProgramCreateRequest struct {
	ProgramName    string  `json:"program_name" binding:"required"`
	ProgramCode   string  `json:"program_code" binding:"required"`
	DepartmentID  string  `json:"department_id" binding:"required"`
	AcademicLevel string  `json:"academic_level" binding:"required,oneof=diploma sarjana magister doktor professionnel"`
	Capacity     int     `json:"capacity"`
	IsActive    *bool   `json:"is_active"`
}

type StudyProgramUpdateRequest struct {
	ProgramName   *string `json:"program_name"`
	DepartmentID *string `json:"department_id"`
	AcademicLevel *string `json:"academic_level"`
	Capacity     *int    `json:"capacity"`
	IsActive    *bool   `json:"is_active"`
}

type StudyProgramHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewStudyProgramHandler(db *gorm.DB) *StudyProgramHandler {
	return &StudyProgramHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

// CreateStudyProgram handles POST /api/v1/study-programs
func (h *StudyProgramHandler) CreateStudyProgram(c *gin.Context) {
	var req StudyProgramCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	sp := domain.StudyProgram{
		ProgramName:   req.ProgramName,
		ProgramCode:  req.ProgramCode,
		DepartmentID: req.DepartmentID,
		AcademicLevel: req.AcademicLevel,
		Capacity:    req.Capacity,
		IsActive:   isActive,
	}

	if err := h.db.Create(&sp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan program studi").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.study_program.create",
		Module:      "pmb",
		ResourceType: "study_program",
		ResourceID:  sp.ID,
		NewValue:   sp,
	})

	c.JSON(http.StatusCreated, sharederr.Success(sp).WithContext(c))
}

// GetStudyProgram handles GET /api/v1/study-programs/:id
func (h *StudyProgramHandler) GetStudyProgram(c *gin.Context) {
	id := c.Param("id")
	sp, err := h.repo.GetStudyProgramByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil program studi").WithContext(c))
		return
	}
	if sp == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Program studi tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(sp).WithContext(c))
}

// GetStudyPrograms handles GET /api/v1/study-programs
func (h *StudyProgramHandler) GetStudyPrograms(c *gin.Context) {
	filter := repository.StudyProgramFilter{
		DepartmentID: c.Query("department_id"),
		AcademicLevel: c.Query("academic_level"),
		IsActive:     c.Query("is_active") == "true",
		Search:       c.Query("search"),
		Page:         1,
		Limit:        20,
	}

	var page, limit int
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filter.Page = page
	filter.Limit = limit

	programs, total, err := h.repo.GetStudyPrograms(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar program studi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"data":  programs,
		"total": total,
		"page": filter.Page,
		"limit": filter.Limit,
	}).WithContext(c))
}

// UpdateStudyProgram handles PUT /api/v1/study-programs/:id
func (h *StudyProgramHandler) UpdateStudyProgram(c *gin.Context) {
	id := c.Param("id")
	var req StudyProgramUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	sp, err := h.repo.GetStudyProgramByID(id)
	if err != nil || sp == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Program studi tidak ditemukan").WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.ProgramName != nil {
		updates["program_name"] = *req.ProgramName
	}
	if req.DepartmentID != nil {
		updates["department_id"] = *req.DepartmentID
	}
	if req.AcademicLevel != nil {
		updates["academic_level"] = *req.AcademicLevel
	}
	if req.Capacity != nil {
		updates["capacity"] = *req.Capacity
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.repo.UpdateStudyProgram(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui program studi").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.study_program.update",
		Module:      "pmb",
		ResourceType: "study_program",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Program studi berhasil diperbarui").WithContext(c))
}

// DeleteStudyProgram handles DELETE /api/v1/study-programs/:id
func (h *StudyProgramHandler) DeleteStudyProgram(c *gin.Context) {
	id := c.Param("id")

	sp, err := h.repo.GetStudyProgramByID(id)
	if err != nil || sp == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Program studi tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.DeleteStudyProgram(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus program studi").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.study_program.delete",
		Module:      "pmb",
		ResourceType: "study_program",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Program studi berhasil dihapus").WithContext(c))
}

// GetStudyProgramQuotas handles GET /api/v1/study-programs/:id/quotas
func (h *StudyProgramHandler) GetStudyProgramQuotas(c *gin.Context) {
	id := c.Param("id")

	quotas, err := h.repo.GetStudyProgramQuotas(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil quota").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(quotas).WithContext(c))
}

// GetAvailableStudyPrograms handles GET /api/v1/study-programs/available
func (h *StudyProgramHandler) GetAvailableStudyPrograms(c *gin.Context) {
	programs, err := h.repo.GetAvailableStudyPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar program studi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(programs).WithContext(c))
}
