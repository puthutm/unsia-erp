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
	"gorm.io/gorm"
)

// CreateInstallmentRequest handles POST /api/v1/finance/invoices/:id/request-installment
func (h *FinanceHandler) CreateInstallmentRequest(c *gin.Context) {
	var req InstallmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	invoiceID := c.Param("id")
	if invoiceID == "" {
		invoiceID = req.InvoiceID
	}

	// Validate invoice status (Req 10.2)
	var invoice domain.Invoice
	if err := h.db.Where("id = ?", invoiceID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Invoice tidak ditemukan").WithContext(c))
		return
	}

	if invoice.Status == "PAID" || invoice.Status == "CANCELLED" || invoice.Status == "EXPIRED" {
		c.JSON(http.StatusUnprocessableEntity, sharederr.Error("INVOICE_NOT_PAYABLE", "Invoice has already been paid, cancelled, or expired").WithContext(c))
		return
	}
	if invoice.Status != "ISSUED" && invoice.Status != "PARTIALLY_PAID" {
		c.JSON(http.StatusUnprocessableEntity, sharederr.Error("INVALID_INVOICE_STATUS", "Invoice status must be ISSUED or PARTIALLY_PAID").WithContext(c))
		return
	}

	installment := domain.InstallmentRequest{
		ID:          fmt.Sprintf("INST-%s", time.Now().Format("20060102150405")),
		InvoiceID:   invoiceID,
		StudentID:   &req.StudentID,
		Reason:      req.Reason,
		Status:      "PENDING",
		RequestedAt: time.Now(),
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

	ir, err := h.repo.GetInstallmentRequestByID(id)
	if err != nil || ir == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengajuan cicilan tidak ditemukan").WithContext(c))
		return
	}

	var invoice domain.Invoice
	if err := h.db.Where("id = ?", ir.InvoiceID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Invoice tidak ditemukan").WithContext(c))
		return
	}

	studentID := ""
	if ir.StudentID != nil {
		studentID = *ir.StudentID
	} else if invoice.StudentID != nil {
		studentID = *invoice.StudentID
	}

	if studentID == "" {
		c.JSON(http.StatusUnprocessableEntity, sharederr.Error("VALIDATION_ERROR", "Student ID is missing").WithContext(c))
		return
	}

	errTx := h.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		updates := map[string]interface{}{
			"status":      "APPROVED",
			"approved_by": &actor,
			"approved_at": &now,
		}
		if err := tx.Model(ir).Updates(updates).Error; err != nil {
			return err
		}

		var clearance domain.StudentClearance
		errCl := tx.Where("student_id = ? AND academic_period_id = ?", studentID, invoice.AcademicPeriodID).First(&clearance).Error
		if errCl == gorm.ErrRecordNotFound {
			clearance = domain.StudentClearance{
				StudentID:        studentID,
				AcademicPeriodID: invoice.AcademicPeriodID,
				ServiceScope:     "registration",
				Status:           "CONDITIONAL",
				Reason:           nil,
				UpdatedBy:        &actor,
				UpdatedAt:        now,
				CreatedAt:        now,
			}
			if err := tx.Create(&clearance).Error; err != nil {
				return err
			}
		} else if errCl == nil {
			clearance.Status = "CONDITIONAL"
			clearance.UpdatedBy = &actor
			clearance.UpdatedAt = now
			if err := tx.Save(&clearance).Error; err != nil {
				return err
			}
		} else {
			return errCl
		}

		return nil
	})

	if errTx != nil {
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
