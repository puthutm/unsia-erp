package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)



// ReceivePaymentCallback handles POST /api/v1/finance/payments/callback/:provider
func (h *FinanceHandler) ReceivePaymentCallback(c *gin.Context) {
	provider := c.Param("provider")
	var req CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	payment, err := h.repo.GetPaymentByID(req.PaymentID)
	if err != nil || payment == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Payment record tidak ditemukan").WithContext(c))
		return
	}

	rawPayload, _ := json.Marshal(req)

	err = h.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 1. Log the callback
		cb := domain.PaymentGatewayCallback{
			PaymentID:         &payment.ID,
			Provider:          provider,
			ProviderEventID:   &req.ProviderEventID,
			ExternalReference: &req.ExternalReference,
			Payload:           string(rawPayload),
			SignatureValid:    true,
			CallbackStatus:    "processed",
			ProcessedAt:       &now,
		}
		if err := tx.Create(&cb).Error; err != nil {
			return err
		}

		// 2. Update payment status
		paymentStatus := "failed"
		if req.PaymentStatus == "success" {
			paymentStatus = "success"
		}
		payment.PaymentStatus = paymentStatus
		payment.PaidAt = &now
		payment.ExternalReference = &req.ExternalReference
		if err := tx.Model(payment).Updates(map[string]interface{}{
			"payment_status":     paymentStatus,
			"paid_at":            &now,
			"external_reference": &req.ExternalReference,
		}).Error; err != nil {
			return err
		}

		// 3. Update invoice paid amount
		var invoice domain.Invoice
		if err := tx.Where("id = ?", payment.InvoiceID).First(&invoice).Error; err != nil {
			return err
		}

		if paymentStatus == "success" {
			invoice.PaidAmount += payment.Amount
			if invoice.PaidAmount >= invoice.TotalAmount {
				invoice.Status = "paid"
			} else if invoice.PaidAmount > 0 {
				invoice.Status = "partially_paid"
			}
			if err := tx.Model(&invoice).Updates(map[string]interface{}{
				"paid_amount": invoice.PaidAmount,
				"status":      invoice.Status,
				"updated_at":  now,
			}).Error; err != nil {
				return err
			}

			// Record double-entry journal logs
			if err := h.RecordPaymentJournal(tx, payment, &invoice); err != nil {
				return err
			}

			// Generate event
			refID := ""
			if invoice.ApplicantID != nil {
				refID = *invoice.ApplicantID
			} else if invoice.StudentID != nil {
				refID = *invoice.StudentID
			}

			envelope := sharedevent.EventEnvelope{
				EventName:        "finance.payment_paid",
				EventVersion:     "v1",
				PublisherService: "finance-service",
				AggregateType:    "payment",
				AggregateID:      payment.ID,
				CorrelationID:    cid,
				Payload: map[string]interface{}{
					"invoice_id":   invoice.ID,
					"payment_id":   payment.ID,
					"amount":       payment.Amount,
					"payer_type":   invoice.TargetType,
					"payer_ref_id": refID,
					"status":       invoice.Status,
				},
			}
			conn := tx.Statement.ConnPool
			_, err = sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
			if err != nil {
				return err
			}

			// Auto clearance if fully paid for student
			if invoice.Status == "paid" && invoice.StudentID != nil {
				clearance := domain.StudentClearance{
					StudentID:        *invoice.StudentID,
					AcademicPeriodID: invoice.AcademicPeriodID,
					ServiceScope:     "registration",
					Status:           "cleared",
					Reason:           nil,
				}
				// Save or update clearance
				var existing domain.StudentClearance
				err := tx.Where("student_id = ? AND academic_period_id = ?", *invoice.StudentID, invoice.AcademicPeriodID).First(&existing).Error
				if err == nil {
					tx.Model(&existing).Updates(map[string]interface{}{
						"status": "cleared",
						"reason": nil,
					})
				} else {
					tx.Create(&clearance)
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses callback").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(payment).WithContext(c))
}

// VerifyManualPayment handles POST /api/v1/finance/payments/verify
func (h *FinanceHandler) VerifyManualPayment(c *gin.Context) {
	var req VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	payment, err := h.repo.GetPaymentByID(req.PaymentID)
	if err != nil || payment == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Payment tidak ditemukan").WithContext(c))
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		ver := domain.PaymentVerification{
			PaymentID:          payment.ID,
			VerifiedBy:         &actor,
			VerificationStatus: req.VerificationStatus,
			RejectionReason:    req.RejectionReason,
			Note:               req.Note,
			VerifiedAt:         &now,
		}
		if err := tx.Create(&ver).Error; err != nil {
			return err
		}

		if req.VerificationStatus == "approved" {
			payment.PaymentStatus = "success"
			payment.PaidAt = &now
			tx.Model(payment).Updates(map[string]interface{}{
				"payment_status": "success",
				"paid_at":        &now,
			})

			var invoice domain.Invoice
			tx.Where("id = ?", payment.InvoiceID).First(&invoice)

			invoice.PaidAmount += req.Amount
			if invoice.PaidAmount >= invoice.TotalAmount {
				invoice.Status = "paid"
			} else {
				invoice.Status = "partially_paid"
			}
			if err := tx.Model(&invoice).Updates(map[string]interface{}{
				"paid_amount": invoice.PaidAmount,
				"status":      invoice.Status,
				"updated_at":  now,
			}).Error; err != nil {
				return err
			}

			// Record double-entry journal logs
			if err := h.RecordPaymentJournal(tx, payment, &invoice); err != nil {
				return err
			}

			refID := ""
			if invoice.ApplicantID != nil {
				refID = *invoice.ApplicantID
			} else if invoice.StudentID != nil {
				refID = *invoice.StudentID
			}

			envelope := sharedevent.EventEnvelope{
				EventName:        "finance.payment_paid",
				EventVersion:     "v1",
				PublisherService: "finance-service",
				AggregateType:    "payment",
				AggregateID:      payment.ID,
				CorrelationID:    cid,
				Payload: map[string]interface{}{
					"invoice_id":   invoice.ID,
					"payment_id":   payment.ID,
					"amount":       payment.Amount,
					"payer_type":   invoice.TargetType,
					"payer_ref_id": refID,
					"status":       invoice.Status,
				},
			}
			conn := tx.Statement.ConnPool
			_, err = sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
			if err != nil {
				return err
			}
		} else {
			payment.PaymentStatus = "failed"
			tx.Model(payment).Update("payment_status", "failed")
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memverifikasi pembayaran").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Pembayaran berhasil diverifikasi").WithContext(c))
}

// GetPayments handles GET /api/v1/finance/payments
func (h *FinanceHandler) GetPayments(c *gin.Context) {
	filter := repository.PaymentListFilter{
		PaymentStatus:  c.Query("payment_status"),
		PaymentMethod: c.Query("payment_method"),
		InvoiceID:    c.Query("invoice_id"),
		StudentID:    c.Query("student_id"),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
		Search:       c.Query("search"),
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

	result, err := h.repo.GetPayments(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar pembayaran").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// RecordPaymentJournal records double-entry journal for payment
func (h *FinanceHandler) RecordPaymentJournal(tx *gorm.DB, payment *domain.Payment, invoice *domain.Invoice) error {
	// Find or Create Debit COA Account (Kas dan Bank - 11100)
	var cashAccount domain.CoaAccount
	err := tx.Where("account_code = ?", "11100").First(&cashAccount).Error
	if err != nil {
		cashAccount = domain.CoaAccount{
			AccountCode:   "11100",
			AccountName:   "Kas dan Bank",
			NormalBalance: "DEBIT",
			IsActive:      true,
		}
		if err := tx.Create(&cashAccount).Error; err != nil {
			return err
		}
	}

	// Find or Create Credit COA Account (Piutang Mahasiswa - 11200)
	var arAccount domain.CoaAccount
	err = tx.Where("account_code = ?", "11200").First(&arAccount).Error
	if err != nil {
		arAccount = domain.CoaAccount{
			AccountCode:   "11200",
			AccountName:   "Piutang Mahasiswa",
			NormalBalance: "CREDIT",
			IsActive:      true,
		}
		if err := tx.Create(&arAccount).Error; err != nil {
			return err
		}
	}

	// Generate JV Number
	journalNum := generateJournalNumber()

	journal := domain.Journal{
		JournalNumber: journalNum,
		JournalDate:   time.Now(),
		SourceType:    "payment",
		SourceID:      &payment.ID,
		Description:   fmt.Sprintf("Penerimaan pembayaran invoice %s", invoice.InvoiceNumber),
	}
	if err := tx.Create(&journal).Error; err != nil {
		return err
	}

	// Debit Entry (Cash)
	debitEntry := domain.JournalEntry{
		JournalID:    journal.ID,
		CoaAccountID: &cashAccount.ID,
		Debit:        payment.Amount,
		Credit:       0,
		Description:  fmt.Sprintf("Debit Kas atas Invoice %s", invoice.InvoiceNumber),
	}
	if err := tx.Create(&debitEntry).Error; err != nil {
		return err
	}

	// Credit Entry (Accounts Receivable)
	creditEntry := domain.JournalEntry{
		JournalID:    journal.ID,
		CoaAccountID: &arAccount.ID,
		Debit:        0,
		Credit:       payment.Amount,
		Description:  fmt.Sprintf("Kredit Piutang atas Invoice %s", invoice.InvoiceNumber),
	}
	if err := tx.Create(&creditEntry).Error; err != nil {
		return err
	}

	return nil
}
