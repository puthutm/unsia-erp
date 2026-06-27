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



// GetDisbursements handles GET /api/v1/finance/disbursements
func (h *FinanceHandler) GetDisbursements(c *gin.Context) {
	filter := repository.DisbursementFilter{
		Status:       c.Query("status"),
		ReferenceType: c.Query("reference_type"),
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

	result, err := h.repo.GetDisbursements(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar disbursement").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateDisbursement handles POST /api/v1/finance/disbursements
func (h *FinanceHandler) CreateDisbursement(c *gin.Context) {
	var req struct {
		ReferenceID   string  `json:"reference_id" binding:"required"`
		ReferenceType string  `json:"reference_type" binding:"required,oneof=commission referral"`
		Amount        float64 `json:"amount" binding:"required"`
		RecipientName string  `json:"recipient_name" binding:"required"`
		BankAccount   *string `json:"bank_account"`
		BankName      *string `json:"bank_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	disbursement := domain.Disbursement{
		ID:                  fmt.Sprintf("DISB-%s", time.Now().Format("20060102150405")),
		DisbursementNumber: "DISB-" + time.Now().Format("20060102150405"),
		ReferenceID:       req.ReferenceID,
		ReferenceType:     req.ReferenceType,
		Amount:           req.Amount,
		RecipientName:     req.RecipientName,
		BankAccount:      req.BankAccount,
		BankName:         req.BankName,
		Status:            "pending",
	}

	if err := h.repo.CreateDisbursement(&disbursement); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan disbursement").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.disbursement.create",
		Module:       "finance",
		ResourceType: "disbursement",
		ResourceID:   disbursement.ID,
		NewValue:     disbursement,
	})

	c.JSON(http.StatusCreated, sharederr.Success(disbursement).WithContext(c))
}

// ApproveDisbursement handles POST /api/v1/finance/disbursements/:id/approve
func (h *FinanceHandler) ApproveDisbursement(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.UpdateDisbursementStatus(id, "approved"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui disbursement").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.disbursement.approve",
		Module:       "finance",
		ResourceType: "disbursement",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Disbursement berhasil disetujui").WithContext(c))
}

// ProcessDisbursement handles POST /api/v1/finance/disbursements/:id/process
func (h *FinanceHandler) ProcessDisbursement(c *gin.Context) {
	id := c.Param("id")

	// Here you would integrate with actual payment gateway
	// For now, we'll just mark as processed
	if err := h.repo.UpdateDisbursementStatus(id, "processed"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses disbursement").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.disbursement.process",
		Module:       "finance",
		ResourceType: "disbursement",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Disbursement berhasil diproses").WithContext(c))
}
