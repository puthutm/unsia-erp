package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)



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
