package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// InvoiceHandler handles invoice-related endpoints
type InvoiceHandler struct {
	*FinanceHandler
}

// NewInvoiceHandler creates a new InvoiceHandler
func NewInvoiceHandler(fh *FinanceHandler) *InvoiceHandler {
	return &InvoiceHandler{FinanceHandler: fh}
}

// CreateInvoice handles POST /api/v1/finance/invoices
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req InvoiceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	var dueTime *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE", "due_date format must be YYYY-MM-DD").WithContext(c))
			return
		}
		dueTime = &parsed
	}

	var applicantID, studentID *string
	if req.PayerType == "applicant" {
		applicantID = &req.PayerRefID
	} else {
		studentID = &req.PayerRefID
	}

	invoice := domain.Invoice{
		InvoiceNumber:    generateInvoiceNumber(),
		TargetType:       req.PayerType,
		ApplicantID:      applicantID,
		StudentID:        studentID,
		AcademicPeriodID: req.AcademicPeriodID,
		Status:           "unpaid",
		DueDate:          dueTime,
	}

	var total float64
	var items []domain.InvoiceItem
	for _, itemReq := range req.Items {
		item := domain.InvoiceItem{
			PaymentComponentID: itemReq.PaymentComponentID,
			Description:        itemReq.Description,
			Amount:             itemReq.Amount,
			FinalAmount:        itemReq.Amount,
		}
		total += itemReq.Amount
		items = append(items, item)
	}
	invoice.TotalAmount = total
	invoice.Items = items

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}

		// Also auto create a pending payment record for this invoice
		payNum := "PAY-" + invoice.InvoiceNumber
		payment := domain.Payment{
			InvoiceID:     invoice.ID,
			PaymentNumber: &payNum,
			Amount:        invoice.TotalAmount,
			PaymentStatus: "pending",
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		envelope := sharedevent.EventEnvelope{
			EventName:        "finance.invoice_created",
			EventVersion:     "v1",
			PublisherService: "finance-service",
			AggregateType:    "invoice",
			AggregateID:      invoice.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"invoice_id":     invoice.ID,
				"invoice_number": invoice.InvoiceNumber,
				"payer_type":     invoice.TargetType,
				"payer_ref_id":   req.PayerRefID,
				"total_amount":   invoice.TotalAmount,
			},
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan invoice").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.invoice.create",
		Module:       "finance",
		ResourceType: "invoice",
		ResourceID:   invoice.ID,
		NewValue:     invoice,
	})

	c.JSON(http.StatusCreated, sharederr.Success(invoice).WithContext(c))
}

// GetInvoice handles GET /api/v1/finance/invoices/:id
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.repo.GetInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil invoice").WithContext(c))
		return
	}
	if inv == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Invoice tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(inv).WithContext(c))
}

// GetInvoices handles GET /api/v1/finance/invoices
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	filter := repository.InvoiceListFilter{
		Status:           c.Query("status"),
		TargetType:       c.Query("target_type"),
		StudentID:       c.Query("student_id"),
		ApplicantID:     c.Query("applicant_id"),
		AcademicPeriodID: c.Query("academic_period_id"),
		DueDateFrom:     c.Query("due_date_from"),
		DueDateTo:       c.Query("due_date_to"),
		Search:          c.Query("search"),
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

	result, err := h.repo.GetInvoices(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar invoice").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}
