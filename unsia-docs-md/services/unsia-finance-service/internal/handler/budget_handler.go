package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)

// BudgetHandler handles budget-related endpoints
type BudgetHandler struct {
	*FinanceHandler
}

// NewBudgetHandler creates a new BudgetHandler
func NewBudgetHandler(fh *FinanceHandler) *BudgetHandler {
	return &BudgetHandler{FinanceHandler: fh}
}

// GetBudgets handles GET /api/v1/finance/budgets
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	filter := repository.BudgetFilter{
		FiscalYear:  0,
		BudgetType: c.Query("budget_type"),
		IsActive:   c.Query("is_active") == "true",
		Search:     c.Query("search"),
	}

	if fy := c.Query("fiscal_year"); fy != "" {
		fmt.Sscanf(fy, "%d", &filter.FiscalYear)
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

	result, err := h.repo.GetBudgets(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar budget").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateBudget handles POST /api/v1/finance/budgets
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	var req BudgetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	budget := domain.Budget{
		ID:                fmt.Sprintf("BUD-%s", time.Now().Format("20060102150405")),
		BudgetNumber:      "BUD-" + time.Now().Format("200601"),
		BudgetName:        req.BudgetName,
		AcademicPeriodID: req.AcademicPeriodID,
		FiscalYear:       req.FiscalYear,
		BudgetType:       req.BudgetType,
		TotalAmount:      req.TotalAmount,
		IsActive:         true,
	}

	if err := h.repo.CreateBudget(&budget); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan budget").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.budget.create",
		Module:       "finance",
		ResourceType: "budget",
		ResourceID:   budget.ID,
		NewValue:     budget,
	})

	c.JSON(http.StatusCreated, sharederr.Success(budget).WithContext(c))
}

// GetBudgetDetail handles GET /api/v1/finance/budgets/:id
func (h *BudgetHandler) GetBudgetDetail(c *gin.Context) {
	id := c.Param("id")

	budget, err := h.repo.GetBudgetByID(id)
	if err != nil || budget == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Budget tidak ditemukan").WithContext(c))
		return
	}

	items, err := h.repo.GetBudgetItems(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil budget items").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"budget": budget,
		"items":  items,
	})).WithContext(c))
}
