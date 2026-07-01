package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)



// GetCashAccounts handles GET /api/v1/finance/cash-accounts
func (h *FinanceHandler) GetCashAccounts(c *gin.Context) {
	filter := repository.CashAccountFilter{
		AccountType: c.Query("account_type"),
		IsActive:    c.Query("is_active") == "true",
		Search:     c.Query("search"),
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

	result, err := h.repo.GetCashAccounts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar cash account").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateCashAccount handles POST /api/v1/finance/cash-accounts
func (h *FinanceHandler) CreateCashAccount(c *gin.Context) {
	var req CashAccountCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	account := domain.CashAccount{
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		AccountType:  req.AccountType,
		BankName:     req.BankName,
		Branch:      req.Branch,
		IsActive:     isActive,
	}

	if err := h.repo.CreateCashAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan cash account").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.cash_account.create",
		Module:       "finance",
		ResourceType: "cash_account",
		ResourceID:   account.ID,
		NewValue:     account,
	})

	c.JSON(http.StatusCreated, sharederr.Success(account).WithContext(c))
}

// GetCashMutations handles GET /api/v1/finance/cash-accounts/:id/mutations
func (h *FinanceHandler) GetCashMutations(c *gin.Context) {
	id := c.Param("id")

	// Verify cash account exists
	existing, err := h.repo.GetCashAccountByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Cash account tidak ditemukan").WithContext(c))
		return
	}

	filter := repository.CashTransactionFilter{}

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

	result, err := h.repo.GetCashMutationsByAccountID(id, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil mutations").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateCashMutation handles POST /api/v1/finance/cash-accounts/:id/mutations
func (h *FinanceHandler) CreateCashMutation(c *gin.Context) {
	id := c.Param("id")
	var req CashMutationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify cash account exists
	existing, err := h.repo.GetCashAccountByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Cash account tidak ditemukan").WithContext(c))
		return
	}

	parsed, err := time.Parse("2006-01-02", req.MutationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE", "mutation_date format must be YYYY-MM-DD").WithContext(c))
		return
	}

	mutation := domain.CashTransaction{
		CashAccountID:   id,
		TransactionType: strings.ToUpper(req.MutationType),
		Amount:          req.Amount,
		TransactionAt:   parsed,
		Description:     req.Description,
	}

	if err := h.repo.CreateCashMutation(&mutation); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan transaction").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.cash_transaction.create",
		Module:       "finance",
		ResourceType: "cash_transaction",
		ResourceID:   mutation.ID,
		NewValue:     mutation,
	})

	c.JSON(http.StatusCreated, sharederr.Success(mutation).WithContext(c))
}
