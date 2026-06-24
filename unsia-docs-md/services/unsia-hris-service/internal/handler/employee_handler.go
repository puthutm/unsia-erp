package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type EmployeeCreateRequest struct {
	PersonID          string  `json:"person_id" binding:"required"`
	EmployeeTypeID   *string `json:"employee_type_id"`
	WorkUnitID       *string `json:"work_unit_id"`
	PositionID       *string `json:"position_id"`
	Nip              *string `json:"nip"`
	EmploymentStatus string  `json:"employment_status"`
}

type EmployeeUpdateRequest struct {
	EmployeeTypeID    *string `json:"employee_type_id"`
	WorkUnitID        *string `json:"work_unit_id"`
	PositionID       *string `json:"position_id"`
	EmploymentStatus  string  `json:"employment_status"`
	IsActive          *bool   `json:"is_active"`
}

type EmployeeHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ListEmployees - GET /api/v1/employees
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.repo.ListEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pegawai").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(employees).WithContext(c))
}

// GetEmployee - GET /api/v1/employees/:id
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	emp, err := h.repo.GetEmployeeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pegawai").WithContext(c))
		return
	}
	if emp == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pegawai tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(emp).WithContext(c))
}

// CreateEmployee - POST /api/v1/employees
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
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
		PersonID:          req.PersonID,
		EmployeeTypeID:   req.EmployeeTypeID,
		WorkUnitID:        req.WorkUnitID,
		PositionID:        req.PositionID,
		Nip:               req.Nip,
		EmploymentStatus:  empStatus,
		IsActive:          true,
	}

	if err := h.repo.CreateEmployee(&emp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data pegawai").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(emp).WithContext(c))
}

// UpdateEmployee - PUT /api/v1/employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var req EmployeeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	emp, err := h.repo.GetEmployeeByID(id)
	if err != nil || emp == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pegawai tidak ditemukan").WithContext(c))
		return
	}

	// Update fields
	if req.EmployeeTypeID != nil {
		emp.EmployeeTypeID = req.EmployeeTypeID
	}
	if req.WorkUnitID != nil {
		emp.WorkUnitID = req.WorkUnitID
	}
	if req.PositionID != nil {
		emp.PositionID = req.PositionID
	}
	if req.EmploymentStatus != "" {
		emp.EmploymentStatus = req.EmploymentStatus
	}
	if req.IsActive != nil {
		emp.IsActive = *req.IsActive
	}

	if err := h.repo.UpdateEmployee(emp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui data pegawai").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(emp).WithContext(c))
}

// DeleteEmployee - DELETE /api/v1/employees/:id
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.DeleteEmployee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus datapegawai").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(nil).WithContext(c))
}

// SearchEmployees - GET /api/v1/employees/search?q=xxx
func (h *EmployeeHandler) SearchEmployees(c *gin.Context) {
	query := c.Query("q")
	employees, err := h.repo.SearchEmployees(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mencaripegawai").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(employees).WithContext(c))
}
