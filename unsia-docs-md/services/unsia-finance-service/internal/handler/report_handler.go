package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
)



// GetBalanceSheet handles GET /api/v1/finance/reports/position
func (h *FinanceHandler) GetBalanceSheet(c *gin.Context) {
	month := 1
	year := time.Now().Year()

	if m := c.Query("month"); m != "" {
		fmt.Sscanf(m, "%d", &month)
	}
	if y := c.Query("year"); y != "" {
		fmt.Sscanf(y, "%d", &year)
	}

	report, err := h.repo.GenerateBalanceSheetReport(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghasilkan report neraca").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(report).WithContext(c))
}

// GetIncomeStatement handles GET /api/v1/finance/reports/activity
func (h *FinanceHandler) GetIncomeStatement(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	report, err := h.repo.GenerateIncomeStatementReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghasilkan report L/R").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(report).WithContext(c))
}

// GetCashFlow handles GET /api/v1/finance/reports/cashflow
func (h *FinanceHandler) GetCashFlow(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	report, err := h.repo.GenerateCashFlowReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghasilkan report arus kas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(report).WithContext(c))
}
