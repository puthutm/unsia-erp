package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-finance-service/internal/service"
	"gorm.io/gorm"
)

// ============ Request Types ============

type InvoiceItemRequest struct {
	PaymentComponentID *string `json:"payment_component_id"`
	Description        string  `json:"description"`
	Amount             float64 `json:"amount" binding:"required"`
	DiscountAmount     float64 `json:"discount_amount"`
	FinalAmount        float64 `json:"final_amount"`
}

type InvoiceCreateRequest struct {
	PayerType        string               `json:"payer_type" binding:"required,oneof=applicant student"`
	PayerRefID       string               `json:"payer_ref_id" binding:"required"`
	AcademicPeriodID *string              `json:"academic_period_id"`
	DueDate          string               `json:"due_date"` // format: YYYY-MM-DD
	Items            []InvoiceItemRequest `json:"items" binding:"required,gt=0"`
	SourceModule     *string              `json:"source_module"`
	SourceRefID      *string              `json:"source_ref_id"`
}

type CallbackRequest struct {
	PaymentID         string      `json:"payment_id" binding:"required"`
	ProviderEventID   string      `json:"provider_event_id" binding:"required"`
	Amount            float64     `json:"amount" binding:"required"`
	PaymentStatus     string      `json:"payment_status" binding:"required,oneof=success failed"`
	ExternalReference string      `json:"external_reference"`
	InvoiceID         string      `json:"invoice_id"`
	SignatureStatus   string      `json:"signature_status"`
	PayloadHash       string      `json:"payload_hash"`
	RawPayload        interface{} `json:"raw_payload"`
}

type VerificationRequest struct {
	InvoiceID          string  `json:"invoice_id"`
	PaymentID          string  `json:"payment_id" binding:"required"`
	VerificationStatus string  `json:"verification_status" binding:"required,oneof=approved rejected"`
	Amount             float64 `json:"amount" binding:"required"`
	PaidAt             *string `json:"paid_at"`
	AttachmentRef      *string `json:"attachment_ref"`
	Reason             *string `json:"reason"`
	Note               *string `json:"note"`
}

type ClearancePolicyCreateRequest struct {
	Code         string `json:"code" binding:"required,alphanum"`
	Name         string `json:"name" binding:"required"`
	ServiceScope string `json:"service_scope" binding:"required,oneof=registration krs graduation"`
	RuleJson     string `json:"rule_json" binding:"required,json"`
	IsActive     *bool  `json:"is_active"`
}

type ClearancePolicyUpdateRequest struct {
	Name         string `json:"name"`
	ServiceScope string `json:"service_scope" binding:"omitempty,oneof=registration krs graduation"`
	RuleJson     string `json:"rule_json" binding:"omitempty,json"`
	IsActive     *bool  `json:"is_active"`
}

type InstallmentCreateRequest struct {
	InvoiceID string `json:"invoice_id" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

type CashAccountCreateRequest struct {
	AccountName   string  `json:"account_name" binding:"required"`
	AccountNumber string  `json:"account_number"`
	AccountType  string  `json:"account_type" binding:"required,oneof=bank cash"`
	BankName      *string `json:"bank_name"`
	Branch       *string `json:"branch"`
	IsActive     *bool   `json:"is_active"`
}

type CashMutationCreateRequest struct {
	MutationType  string  `json:"mutation_type" binding:"required,oneof=debit credit"`
	Amount        float64 `json:"amount" binding:"required"`
	MutationDate  string  `json:"mutation_date" binding:"required"` // YYYY-MM-DD
	Description  string  `json:"description"`
	Reference    *string `json:"reference"`
}

type BudgetCreateRequest struct {
	BudgetName        string  `json:"budget_name" binding:"required"`
	AcademicPeriodID *string `json:"academic_period_id"`
	FiscalYear       int     `json:"fiscal_year" binding:"required"`
	BudgetType       string  `json:"budget_type" binding:"required,oneof=operational capital investment"`
	TotalAmount      float64 `json:"total_amount" binding:"required"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
}

type VendorCreateRequest struct {
	VendorName    string  `json:"vendor_name" binding:"required"`
	ContactPerson *string `json:"contact_person"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
	TaxNumber    *string `json:"tax_number"`
}

type POCreateRequest struct {
	VendorID     string           `json:"vendor_id" binding:"required"`
	PODate       string           `json:"po_date" binding:"required"`
	DeliveryDate *string         `json:"delivery_date"`
	Description string           `json:"description"`
	TotalAmount float64          `json:"total_amount" binding:"required"`
	Items       []POItemRequest  `json:"items"`
}

type POItemRequest struct {
	Description string  `json:"description" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
}

type EventCreateRequest struct {
	EventName    string  `json:"event_name" binding:"required"`
	EventType    string  `json:"event_type" binding:"required,oneof=graduation seminar workshop"`
	EventDate    string  `json:"event_date" binding:"required"`
	Description *string `json:"description"`
	BudgetAmount float64 `json:"budget_amount"`
}

type ScholarshipCreateRequest struct {
	StudentID       string  `json:"student_id" binding:"required,uuid"`
	ScholarshipType string  `json:"scholarship_type" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	Status          string  `json:"status"`
}

type ScholarshipUpdateRequest struct {
	ScholarshipType *string  `json:"scholarship_type"`
	Amount          *float64 `json:"amount"`
	Status          *string  `json:"status"`
	ApprovedBy      *string  `json:"approved_by"`
}

// ============ FinanceHandler Struct ============

type FinanceHandler struct {
	repo                  *repository.FinanceRepository
	db                    *gorm.DB
	InvoiceService        *service.InvoiceService
	PaymentGatewayService *service.PaymentGatewayService
}

func NewFinanceHandler(db *gorm.DB) *FinanceHandler {
	return &FinanceHandler{
		repo: repository.NewFinanceRepository(db),
		db:   db,
	}
}

// ============ Shared Functions ============

func generateInvoiceNumber() string {
	now := time.Now().Format("20060102")
	nBig, _ := rand.Int(rand.Reader, big.NewInt(900000))
	num := nBig.Int64() + 100000
	return fmt.Sprintf("INV%s%d", now, num)
}

func generateJournalNumber() string {
	now := time.Now()
	nBig, _ := rand.Int(rand.Reader, big.NewInt(900000))
	num := nBig.Int64() + 100000
	return fmt.Sprintf("JV-%s%d", now.Format("20060102"), num)
}


