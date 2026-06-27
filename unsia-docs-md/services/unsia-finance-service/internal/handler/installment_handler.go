package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
)

// CreateInstallmentRequest handles POST /api/v1/finance/installment-requests
func (h *FinanceHandler) CreateInstallmentRequest(c *gin.Context) {
	var req InstallmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	installment := domain.InstallmentRequest{
		ID:        fmt.Sprintf("INST-%s", time.Now().Format("20060102150405")),
		InvoiceID: req.InvoiceID,
		StudentID: &req.StudentID,
		Reason:    req.Reason,
		Status:    "pending",
	}

	if err := h.repo.CreateInstallmentRequest(&installment); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat pengajuan cicilan").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.installment.request",
		Module:       "finance",
		ResourceType: "installment_request",
		ResourceID:   installment.ID,
		NewValue:     installment,
	})

	c.JSON(http.StatusCreated, sharederr.Success(installment).WithContext(c))
}

// ApproveInstallmentRequest handles PATCH /api/v1/finance/installment-requests/:id/approve
func (h *FinanceHandler) ApproveInstallmentRequest(c *gin.Context) {
	id := c.Param("id")

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	if err := h.repo.UpdateInstallmentRequestStatus(id, "approved", actor); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui cicilan").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.installment.approve",
		Module:       "finance",
		ResourceType: "installment_request",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Pengajuan cicilan berhasil disetujui").WithContext(c))
}
