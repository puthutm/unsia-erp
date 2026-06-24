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

// ExpenseEventHandler handles expense event-related endpoints
type ExpenseEventHandler struct {
	*FinanceHandler
}

// NewExpenseEventHandler creates a new ExpenseEventHandler
func NewExpenseEventHandler(fh *FinanceHandler) *ExpenseEventHandler {
	return &ExpenseEventHandler{FinanceHandler: fh}
}

// GetExpenseEvents handles GET /api/v1/finance/events
func (h *ExpenseEventHandler) GetExpenseEvents(c *gin.Context) {
	filter := repository.EventFilter{
		EventType: c.Query("event_type"),
		Status:    c.Query("status"),
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

	result, err := h.repo.GetExpenseEvents(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar event").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateExpenseEvent handles POST /api/v1/finance/events
func (h *ExpenseEventHandler) CreateExpenseEvent(c *gin.Context) {
	var req EventCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE", "event_date format must be YYYY-MM-DD").WithContext(c))
		return
	}

	event := domain.ExpenseEvent{
		ID:            fmt.Sprintf("EVT-%s", time.Now().Format("20060102150405")),
		EventName:    req.EventName,
		EventType:    req.EventType,
		EventDate:    parsedDate,
		Description: req.Description,
		BudgetAmount: req.BudgetAmount,
		Status:       "planned",
	}

	if err := h.repo.CreateExpenseEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan event").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.event.create",
		Module:       "finance",
		ResourceType: "expense_event",
		ResourceID:   event.ID,
		NewValue:     event,
	})

	c.JSON(http.StatusCreated, sharederr.Success(event).WithContext(c))
}
