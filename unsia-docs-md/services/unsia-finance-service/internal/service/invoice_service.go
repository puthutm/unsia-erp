package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/state_machine"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedevent "github.com/unsia-erp/shared-event"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/gorm"
)

// InvoiceService handles invoice business logic
type InvoiceService struct {
	db *gorm.DB
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(db *gorm.DB) *InvoiceService {
	return &InvoiceService{db: db}
}

// InvoiceListParams represents parameters for listing invoices
type InvoiceListParams struct {
	Status           string
	TargetType       string
	ApplicantID      string
	StudentID        string
	AcademicPeriodID  string
	DueDateFrom      *time.Time
	DueDateTo        *time.Time
	Page             int
	Limit            int
}

// InvoiceListResult represents the result of listing invoices
type InvoiceListResult struct {
	Invoices   []domain.Invoice `json:"invoices"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalCount int64            `json:"total_count"`
	TotalPages int             `json:"total_pages"`
}

// ListInvoices retrieves a paginated list of invoices
func (s *InvoiceService) ListInvoices(params InvoiceListParams) (*InvoiceListResult, error) {
	query := s.db.Model(&domain.Invoice{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.TargetType != "" {
		query = query.Where("target_type = ?", params.TargetType)
	}
	if params.ApplicantID != "" {
		query = query.Where("applicant_id = ?", params.ApplicantID)
	}
	if params.StudentID != "" {
		query = query.Where("student_id = ?", params.StudentID)
	}
	if params.AcademicPeriodID != "" {
		query = query.Where("academic_period_id = ?", params.AcademicPeriodID)
	}
	if params.DueDateFrom != nil {
		query = query.Where("due_date >= ?", params.DueDateFrom)
	}
	if params.DueDateTo != nil {
		query = query.Where("due_date <= ?", params.DueDateTo)
	}

	// Get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Calculate pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Fetch invoices with items
	var invoices []domain.Invoice
	if err := query.
		Preload("Items").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&invoices).Error; err != nil {
		return nil, err
	}

	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return &InvoiceListResult{
		Invoices:   invoices,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateInvoiceStatus updates invoice status with state machine validation
func (s *InvoiceService) UpdateInvoiceStatus(invoiceID, newStatus string, actor string, c *gin.Context) error {
	// Get existing invoice
	var invoice domain.Invoice
	if err := s.db.First(&invoice, "id = ?", invoiceID).Error; err != nil {
		return err
	}

	// Validate state machine transition
	sm := state_machine.NewInvoiceStateMachine()
	oldStatus := invoice.Status
	if err := sm.ValidateTransition(state_machine.InvoiceStatus(oldStatus), state_machine.InvoiceStatus(newStatus)); err != nil {
		return err
	}

	// Update status
	now := time.Now()
	updateValues := map[string]interface{}{
		"status":      newStatus,
		"updated_at": now,
	}

	// Record the status change
	if err := s.db.Model(&invoice).Updates(updateValues).Error; err != nil {
		return err
	}

	// Audit log
	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.invoice.status_change",
		Module:      "finance",
		ResourceType: "invoice",
		ResourceID:  invoiceID,
		OldValue:    oldStatus,
		NewValue:    newStatus,
		UserID:      actor,
	})

	// Write outbox event if status changed to ISSUED
	if newStatus == "ISSUED" && oldStatus == "DRAFT" {
		correlationID, _ := c.Get("x-correlation-id")
		cid, _ := correlationID.(string)

		envelope := sharedevent.EventEnvelope{
			EventName:        "finance.invoice_issued",
			EventVersion:     "v1",
			PublisherService: "finance-service",
			AggregateType:    "invoice",
			AggregateID:      invoiceID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"invoice_id":     invoice.ID,
				"invoice_number": invoice.InvoiceNumber,
				"old_status":    oldStatus,
				"new_status":    newStatus,
			},
		}

		conn := s.db.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateInvoice creates a new invoice with validation
func (s *InvoiceService) CreateInvoice(req *InvoiceCreateRequest, c *gin.Context) (*domain.Invoice, error) {
	// Validate due date
	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		return nil, fmt.Errorf("due_date cannot be in the past")
	}

	// Validate items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("at least one invoice item is required")
	}

	// Generate invoice number
	invoiceNumber := generateInvoiceNumber()

	// Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		if item.FinalAmount <= 0 {
			return nil, fmt.Errorf("item amount must be positive")
		}
		totalAmount += item.FinalAmount
	}

	// Determine target type
	var applicantID, studentID *string
	if req.PayerType == "applicant" {
		applicantID = &req.PayerRefID
	} else {
		studentID = &req.PayerRefID
	}

	// Create invoice
	invoice := domain.Invoice{
		InvoiceNumber:    invoiceNumber,
		TargetType:       req.PayerType,
		ApplicantID:      applicantID,
		StudentID:        studentID,
		AcademicPeriodID: req.AcademicPeriodID,
		TotalAmount:      totalAmount,
		PaidAmount:       0,
		Status:           "DRAFT",
		DueDate:          req.DueDate,
	}

	// Create invoice items
	for _, itemReq := range req.Items {
		item := domain.InvoiceItem{
			PaymentComponentID: itemReq.PaymentComponentID,
			Description:        itemReq.Description,
			Amount:             itemReq.Amount,
			DiscountAmount:    itemReq.DiscountAmount,
			FinalAmount:        itemReq.FinalAmount,
		}
		invoice.Items = append(invoice.Items, item)
	}

	// Save in transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}

		// Create pending payment
		payNum := "PAY-" + invoiceNumber
		payment := domain.Payment{
			InvoiceID:     invoice.ID,
			PaymentNumber: &payNum,
			Amount:        invoice.TotalAmount,
			PaymentStatus: "pending",
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		// Write outbox event
		correlationID, _ := c.Get("x-correlation-id")
		cid, _ := correlationID.(string)

		refID := ""
		if invoice.ApplicantID != nil {
			refID = *invoice.ApplicantID
		} else if invoice.StudentID != nil {
			refID = *invoice.StudentID
		}

		envelope := sharedevent.EventEnvelope{
			EventName:        "finance.invoice_created",
			EventVersion:     "v1",
			PublisherService: "finance-service",
			AggregateType:    "invoice",
			AggregateID:      invoice.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"invoice_id":      invoice.ID,
				"invoice_number":  invoice.InvoiceNumber,
				"payer_type":      invoice.TargetType,
				"payer_ref_id":     refID,
				"total_amount":    invoice.TotalAmount,
				"due_date":       invoice.DueDate,
			},
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// InvoiceCreateRequest represents the request to create an invoice
type InvoiceCreateRequest struct {
	PayerType        string     `json:"payer_type" binding:"required,oneof=applicant student"`
	PayerRefID       string     `json:"payer_ref_id" binding:"required"`
	AcademicPeriodID *string    `json:"academic_period_id"`
	DueDate          *time.Time `json:"due_date"`
	Items            []InvoiceItemRequest `json:"items" binding:"required,gt=0"`
}

// InvoiceItemRequest represents an invoice item
type InvoiceItemRequest struct {
	PaymentComponentID *string  `json:"payment_component_id"`
	Description        string    `json:"description"`
	Amount             float64   `json:"amount" binding:"required,gt=0"`
	DiscountAmount     float64   `json:"discount_amount"`
	FinalAmount       float64   `json:"final_amount" binding:"required,gt=0"`
}

func generateInvoiceNumber() string {
	now := time.Now().Format("20060102")
	nBig, _ := rand.Int(rand.Reader, big.NewInt(900000))
	num := nBig.Int64() + 100000
	return fmt.Sprintf("INV%s%d", now, num)
}

// IdempotencyCheck checks if idempotency key exists
func (s *InvoiceService) IdempotencyCheck(idempotencyKey string) (bool, string, error) {
	type IdempotencyKey struct {
		ID         string    `gorm:"primaryKey;column:id"`
		Key        string    `gorm:"column:key;unique"`
		Response   string    `gorm:"column:response"`
		Status     string    `gorm:"column:status"`
		CreatedAt  time.Time `gorm:"column:created_at"`
	}

	var key IdempotencyKey
	err := s.db.First(&key, "key = ?", idempotencyKey).Error

	if err == gorm.ErrRecordNotFound {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	if key.Status == "completed" {
		return true, key.Response, nil
	}

	return false, key.ID, nil
}

// SaveIdempotency saves idempotency key
func (s *InvoiceService) SaveIdempotency(idempotencyKey, response string, status string) error {
	type IdempotencyKey struct {
		ID       string    `gorm:"primaryKey;column:id"`
		Key      string    `gorm:"column:key;unique"`
		Response string   `gorm:"column:response"`
		Status   string    `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	key := IdempotencyKey{
		Key:      idempotencyKey,
		Response: response,
		Status:   status,
		CreatedAt: time.Now(),
	}

	return s.db.Create(&key).Error
}

// UpdatePaidAmount updates invoice paid amount with validation
func (s *InvoiceService) UpdatePaidAmount(invoiceID string, amount float64) error {
	var invoice domain.Invoice
	if err := s.db.First(&invoice, "id = ?", invoiceID).Error; err != nil {
		return err
	}

	newPaidAmount := invoice.PaidAmount + amount

	// Check for overpayment
	if newPaidAmount > invoice.TotalAmount {
		return fmt.Errorf("overpayment not allowed: paid_amount would exceed total_amount")
	}

	// Update based on new amount
	newStatus := invoice.Status
	if newPaidAmount >= invoice.TotalAmount {
		newStatus = "PAID"
	} else if newPaidAmount > 0 {
		newStatus = "PARTIALLY_PAID"
	}

	return s.db.Model(&invoice).Updates(map[string]interface{}{
		"paid_amount": newPaidAmount,
		"status":      newStatus,
		"updated_at":  time.Now(),
	}).Error
}
