package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type LecturerCreateRequest struct {
	EmployeeID             string  `json:"employee_id" binding:"required"`
	LecturerStatusID       *string `json:"lecturer_status_id"`
	FunctionalPositionID  *string `json:"functional_position_id"`
	Nidn                   *string `json:"nidn"`
	HomebaseStudyProgramID *string `json:"homebase_study_program_id"`
	CertificationStatus   *string `json:"certification_status"`
}

type LecturerUpdateRequest struct {
	LecturerStatusID       *string `json:"lecturer_status_id"`
	FunctionalPositionID  *string `json:"functional_position_id"`
	Nidn                   *string `json:"nidn"`
	HomebaseStudyProgramID *string `json:"homebase_study_program_id"`
	CertificationStatus   *string `json:"certification_status"`
	IsActive               *bool   `json:"is_active"`
}

type LecturerHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewLecturerHandler(db *gorm.DB) *LecturerHandler {
	return &LecturerHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ListActiveLecturers - GET /api/v1/lecturers
func (h *LecturerHandler) ListActiveLecturers(c *gin.Context) {
	lecturers, err := h.repo.ListActiveLecturers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data dosen").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(lecturers).WithContext(c))
}

// GetLecturer - GET /api/v1/lecturers/:id
func (h *LecturerHandler) GetLecturer(c *gin.Context) {
	id := c.Param("id")
	lec, err := h.repo.GetLecturerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data dosen").WithContext(c))
		return
	}
	if lec == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Dosen tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(lec).WithContext(c))
}

// CreateLecturer - POST /api/v1/lecturers
func (h *LecturerHandler) CreateLecturer(c *gin.Context) {
	var req LecturerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	lec := domain.Lecturer{
		EmployeeID:             req.EmployeeID,
		LecturerStatusID:       req.LecturerStatusID,
		FunctionalPositionID:   req.FunctionalPositionID,
		Nidn:                   req.Nidn,
		HomebaseStudyProgramID: req.HomebaseStudyProgramID,
		CertificationStatus:    req.CertificationStatus,
		IsActive:              true,
	}

	if err := h.repo.CreateLecturer(&lec); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data dosen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(lec).WithContext(c))
}

// UpdateLecturer - PUT /api/v1/lecturers/:id
func (h *LecturerHandler) UpdateLecturer(c *gin.Context) {
	id := c.Param("id")
	var req LecturerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	lec, err := h.repo.GetLecturerByID(id)
	if err != nil || lec == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Dosen tidak ditemukan").WithContext(c))
		return
	}

	// Update fields
	if req.LecturerStatusID != nil {
		lec.LecturerStatusID = req.LecturerStatusID
	}
	if req.FunctionalPositionID != nil {
		lec.FunctionalPositionID = req.FunctionalPositionID
	}
	if req.Nidn != nil {
		lec.Nidn = req.Nidn
	}
	if req.HomebaseStudyProgramID != nil {
		lec.HomebaseStudyProgramID = req.HomebaseStudyProgramID
	}
	if req.CertificationStatus != nil {
		lec.CertificationStatus = req.CertificationStatus
	}
	if req.IsActive != nil {
		lec.IsActive = *req.IsActive
	}

	if err := h.repo.UpdateLecturer(lec); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui data dosen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(lec).WithContext(c))
}

// DeleteLecturer - DELETE /api/v1/lecturers/:id
func (h *LecturerHandler) DeleteLecturer(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.DeleteLecturer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus data dosen").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(nil).WithContext(c))
}

// SearchLecturers - GET /api/v1/lecturers/search?q=xxx
func (h *LecturerHandler) SearchLecturers(c *gin.Context) {
	query := c.Query("q")
	lecturers, err := h.repo.SearchLecturers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mencari dosen").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(lecturers).WithContext(c))
}
