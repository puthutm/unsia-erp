package handler

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type PayrollRunCreateRequest struct {
	PayrollPeriod string               `json:"payroll_period" binding:"required"`
	RunDate       string               `json:"run_date" binding:"required"` // format: YYYY-MM-DD
	Items         []PayrollItemRequest `json:"items" binding:"required,gt=0"`
}

type PayrollItemRequest struct {
	EmployeeID      string  `json:"employee_id" binding:"required"`
	GrossAmount     float64 `json:"gross_amount"`
	DeductionAmount float64 `json:"deduction_amount"`
	NetAmount       float64 `json:"net_amount"`
}

// GetPayrollRuns handles GET /api/v1/finance/payroll-runs
func (h *FinanceHandler) GetPayrollRuns(c *gin.Context) {
	filter := repository.PayrollFilter{
		Status:      c.Query("status"),
		PeriodMonth: 0,
		PeriodYear:  0,
	}

	if pm := c.Query("period_month"); pm != "" {
		fmt.Sscanf(pm, "%d", &filter.PeriodMonth)
	}
	if py := c.Query("period_year"); py != "" {
		fmt.Sscanf(py, "%d", &filter.PeriodYear)
	}

	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filter.Page = page
	filter.Limit = limit

	result, err := h.repo.GetPayrollRuns(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar payroll").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// ApprovePayrollRun handles POST /api/v1/finance/payroll-runs/:id/approve
func (h *FinanceHandler) ApprovePayrollRun(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.UpdatePayrollRunStatus(id, "approved"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui payroll").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.payroll.approve",
		Module:       "finance",
		ResourceType: "payroll_run",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Payroll berhasil disetujui").WithContext(c))
}

// CreatePayrollRun handles POST /api/v1/finance/payroll-runs
func (h *FinanceHandler) CreatePayrollRun(c *gin.Context) {
	var req PayrollRunCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.RunDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE", "run_date format must be YYYY-MM-DD").WithContext(c))
		return
	}

	var totalAmount float64
	for _, item := range req.Items {
		if math.Abs(item.NetAmount-(item.GrossAmount-item.DeductionAmount)) > 0.001 {
			c.JSON(http.StatusUnprocessableEntity, sharederr.Error("PAYROLL_AMOUNT_INVALID", "net_amount must equal gross_amount minus deduction_amount").WithContext(c))
			return
		}
		totalAmount += item.NetAmount
	}

	run := domain.PayrollRun{
		PayrollPeriod: req.PayrollPeriod,
		RunDate:       parsedDate,
		TotalAmount:   totalAmount,
		Status:        "draft",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	errTx := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&run).Error; err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			item := domain.PayrollItem{
				PayrollRunID:    run.ID,
				EmployeeID:      itemReq.EmployeeID,
				GrossAmount:     itemReq.GrossAmount,
				DeductionAmount: itemReq.DeductionAmount,
				NetAmount:       itemReq.NetAmount,
				Status:          "draft",
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if errTx != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan payroll run").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.payroll.create",
		Module:       "finance",
		ResourceType: "payroll_run",
		ResourceID:   run.ID,
		NewValue:     run,
	})

	c.JSON(http.StatusCreated, sharederr.Success(run).WithContext(c))
}
