package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type EmployeeCreateRequest struct {
	PersonID         string  `json:"person_id" binding:"required"`
	EmployeeTypeID   *string `json:"employee_type_id"`
	WorkUnitID       *string `json:"work_unit_id"`
	PositionID       *string `json:"position_id"`
	Nip              *string `json:"nip"`
	EmploymentStatus string  `json:"employment_status"`
}

type LecturerCreateRequest struct {
	EmployeeID             string  `json:"employee_id" binding:"required"`
	LecturerStatusID       *string `json:"lecturer_status_id"`
	FunctionalPositionID   *string `json:"functional_position_id"`
	Nidn                   *string `json:"nidn"`
	HomebaseStudyProgramID *string `json:"homebase_study_program_id"`
	CertificationStatus    *string `json:"certification_status"`
}

type PositionPlotRequest struct {
	PositionID string `json:"position_id" binding:"required"`
}

type BkdRecordCreateRequest struct {
	LecturerID       string  `json:"lecturer_id" binding:"required"`
	AcademicPeriodID *string `json:"academic_period_id"`
	TeachingLoad     float64 `json:"teaching_load"`
	ResearchLoad     float64 `json:"research_load"`
	ServiceLoad      float64 `json:"service_load"`
}

type HRISHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewHRISHandler(db *gorm.DB) *HRISHandler {
	return &HRISHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

func (h *HRISHandler) ListActiveLecturers(c *gin.Context) {
	lecturers, err := h.repo.ListActiveLecturers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data dosen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(lecturers).WithContext(c))
}

func (h *HRISHandler) GetLecturer(c *gin.Context) {
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

func (h *HRISHandler) CreateEmployee(c *gin.Context) {
	var req EmployeeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	empStatus := req.EmploymentStatus
	if empStatus == "" {
		empStatus = "contract"
	}

	emp := domain.Employee{
		PersonID:         req.PersonID,
		EmployeeTypeID:   req.EmployeeTypeID,
		WorkUnitID:       req.WorkUnitID,
		PositionID:       req.PositionID,
		Nip:              req.Nip,
		EmploymentStatus: empStatus,
		IsActive:         true,
	}

	if err := h.repo.CreateEmployee(&emp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data pegawai").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(emp).WithContext(c))
}

func (h *HRISHandler) CreateLecturer(c *gin.Context) {
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
		IsActive:               true,
	}

	if err := h.repo.CreateLecturer(&lec); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data dosen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(lec).WithContext(c))
}

func (h *HRISHandler) PlotPosition(c *gin.Context) {
	id := c.Param("id")
	var req PositionPlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	emp, err := h.repo.GetEmployeeByID(id)
	if err != nil || emp == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pegawai tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.UpdateEmployeePosition(id, req.PositionID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui jabatan pegawai").WithContext(c))
		return
	}

	emp.PositionID = &req.PositionID
	c.JSON(http.StatusOK, sharederr.Success(emp).WithContext(c))
}

func (h *HRISHandler) CreateBkdRecord(c *gin.Context) {
	var req BkdRecordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	br := domain.BkdRecord{
		LecturerID:       req.LecturerID,
		AcademicPeriodID: req.AcademicPeriodID,
		TeachingLoad:     req.TeachingLoad,
		ResearchLoad:     req.ResearchLoad,
		ServiceLoad:      req.ServiceLoad,
		Status:           "draft",
	}

	if err := h.repo.CreateBkdRecord(&br); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat record BKD dosen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(br).WithContext(c))
}
