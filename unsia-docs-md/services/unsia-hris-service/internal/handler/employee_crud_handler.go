package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/domain"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type EmployeeCRUDHandler struct {
	repo *repository.HrisRepository
	db   *gorm.DB
}

func NewEmployeeCRUDHandler(db *gorm.DB) *EmployeeCRUDHandler {
	return &EmployeeCRUDHandler{
		repo: repository.NewHrisRepository(db),
		db:   db,
	}
}

// GET /api/v1/employees - List employees
func (h *EmployeeCRUDHandler) ListEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	departmentID := c.Query("department_id")
	status := c.Query("status")

	employees, total, err := h.repo.ListEmployees(page, limit, departmentID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data karyawan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(employees).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/employees - Create employee
func (h *EmployeeCRUDHandler) CreateEmployee(c *gin.Context) {
	var req struct {
		NIP           string `json:"nip" binding:"required"`
		PersonID     string `json:"person_id" binding:"required"`
		DepartmentID string `json:"department_id"`
		Position     string `json:"position"`
		EmployeeType string `json:"employee_type"` // tetap, kontrak, harian
		JoinDate     string `json:"join_date"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	employee := domain.Employee{
		NIP:           req.NIP,
		PersonID:      req.PersonID,
		DepartmentID: req.DepartmentID,
		Position:     req.Position,
		EmployeeType: req.EmployeeType,
		JoinDate:     req.JoinDate,
		Status:       req.Status,
	}

	if err := h.repo.CreateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat karyawan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(employee).WithContext(c))
}

// GET /api/v1/employees/:id - Get employee
func (h *EmployeeCRUDHandler) GetEmployee(c *gin.Context) {
	employeeID := c.Param("id")

	employee, err := h.repo.GetEmployeeByID(employeeID)
	if err != nil || employee == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Karyawan tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(employee).WithContext(c))
}

// PUT /api/v1/employees/:id - Update employee
func (h *EmployeeCRUDHandler) UpdateEmployee(c *gin.Context) {
	employeeID := c.Param("id")

	var req struct {
		DepartmentID string `json:"department_id"`
		Position    string `json:"position"`
		EmployeeType string `json:"employee_type"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.DepartmentID != "" {
		updates["department_id"] = req.DepartmentID
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.EmployeeType != "" {
		updates["employee_type"] = req.EmployeeType
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.repo.UpdateEmployee(employeeID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui karyawan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Karyawan berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/employees/:id - Delete employee
func (h *EmployeeCRUDHandler) DeleteEmployee(c *gin.Context) {
	employeeID := c.Param("id")

	if err := h.repo.DeleteEmployee(employeeID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus karyawan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Karyawan berhasil dihapus").WithContext(c))
}
