package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// Payroll structures - Integration with Finance Service
type PayrollGenerateRequest struct {
	PayPeriodID    string `json:"pay_period_id" binding:"required"`
	WorkUnitID    *string `json:"work_unit_id"`
}

type PayrollApprovalRequest struct {
	Status string `json:"status" binding:"required"` // approved, rejected
	Notes  string `json:"notes"`
}

type PayrollHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewPayrollHandler(db *gorm.DB) *PayrollHandler {
	return &PayrollHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ListPayrolls - GET /api/v1/payrolls
func (h *PayrollHandler) ListPayrolls(c *gin.Context) {
	payPeriodID := c.Query("pay_period_id")
	workUnitID := c.Query("work_unit_id")
	status := c.Query("status")

	payrolls, err := h.repo.ListPayrolls(payPeriodID, workUnitID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data penggajian").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(payrolls).WithContext(c))
}

// GetPayroll - GET /api/v1/payrolls/:id
func (h *PayrollHandler) GetPayroll(c *gin.Context) {
	id := c.Param("id")
	payroll, err := h.repo.GetPayrollByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data penggajian").WithContext(c))
		return
	}
	if payroll == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Data penggajian tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(payroll).WithContext(c))
}

// GeneratePayroll - POST /api/v1/payrolls/generate
// Trigger payroll calculation from HRIS side (syncs with Finance)
func (h *PayrollHandler) GeneratePayroll(c *gin.Context) {
	var req PayrollGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get employees
	employees, err := h.repo.ListEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data karyawan").WithContext(c))
		return
	}

	// Call Finance service to calculate payroll
	// This would be an internal service call
	payrolls, err := h.repo.GeneratePayroll(req.PayPeriodID, employees)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghasilkan payroll").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(payrolls).WithContext(c))
}

// ApprovePayroll - PUT /api/v1/payrolls/:id/approve
func (h *PayrollHandler) ApprovePayroll(c *gin.Context) {
	id := c.Param("id")
	var req PayrollApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	payroll, err := h.repo.GetPayrollByID(id)
	if err != nil || payroll == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Data penggajian tidak ditemukan").WithContext(c))
		return
	}

	payroll.Status = req.Status
	payroll.Notes = req.Notes

	if err := h.repo.UpdatePayroll(payroll); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui payroll").WithContext(c))
		return
	}

	// Trigger payment to Finance service
	// TODO: Call Finance service to process payment

	c.JSON(http.StatusOK, sharederr.Success(payroll).WithContext(c))
}

// GetMyPayroll - GET /api/v1/payrolls/my (Employee view own payroll)
func (h *PayrollHandler) GetMyPayroll(c *gin.Context) {
	employeeID := c.GetHeader("X-Employee-ID")
	if employeeID == "" {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Employee ID diperlukan").WithContext(c))
		return
	}

	payrolls, err := h.repo.GetEmployeePayrolls(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil slip gaji").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(payrolls).WithContext(c))
}

// GetPayrollSummary - GET /api/v1/payrolls/summary
func (h *PayrollHandler) GetPayrollSummary(c *gin.Context) {
	payPeriodID := c.Query("pay_period_id")
	workUnitID := c.Query("work_unit_id")

	summary, err := h.repo.GetPayrollSummary(payPeriodID, workUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil ringkasan payroll").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(summary).WithContext(c))
}

// ===== PAY PERIOD =====

type PayPeriodCreateRequest struct {
	Month      int    `json:"month" binding:"required"`  // 1-12
	Year      int    `json:"year" binding:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ListPayPeriods - GET /api/v1/pay-periods
func (h *PayrollHandler) ListPayPeriods(c *gin.Context) {
	periods, err := h.repo.ListPayPeriods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil periode gaji").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(periods).WithContext(c))
}

// CreatePayPeriod - POST /api/v1/pay-periods
func (h *PayrollHandler) CreatePayPeriod(c *gin.Context) {
	var req PayPeriodCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	period := domain.PayPeriod{
		Month:     req.Month,
		Year:      req.Year,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "open", // open, calculating, approved, paid
	}

	if err := h.repo.CreatePayPeriod(&period); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan periode").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(period).WithContext(c))
}

// ===== SALARY COMPONENT =====

type SalaryComponentRequest struct {
	Name         string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Type        string  `json:"type" binding:"required"` // allowance, deduction
	Amount      float64 `json:"amount"`
	IsTaxable    bool    `json:"is_taxable"`
	IsPermanent bool    `json:"is_permanent"`
}

// ListSalaryComponents - GET /api/v1/salary-components
func (h *PayrollHandler) ListSalaryComponents(c *gin.Context) {
	components, err := h.repo.ListSalaryComponents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil komponen gaji").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(components).WithContext(c))
}

// CreateSalaryComponent - POST /api/v1/salary-components
func (h *PayrollHandler) CreateSalaryComponent(c *gin.Context) {
	var req SalaryComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	component := domain.SalaryComponent{
		Name:          req.Name,
		Code:          req.Code,
		Type:          req.Type,
		Amount:        req.Amount,
		IsTaxable:     req.IsTaxable,
		IsPermanent:   req.IsPermanent,
	}

	if err := h.repo.CreateSalaryComponent(&component); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan komponen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(component).WithContext(c))
}
