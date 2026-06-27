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
}

type InvoiceCreateRequest struct {
	PayerType        string               `json:"payer_type" binding:"required,oneof=applicant student"`
	PayerRefID       string               `json:"payer_ref_id" binding:"required"`
	AcademicPeriodID *string              `json:"academic_period_id"`
	DueDate          string               `json:"due_date"` // format: YYYY-MM-DD
	Items            []InvoiceItemRequest `json:"items" binding:"required,gt=0"`
}

type CallbackRequest struct {
	PaymentID         string  `json:"payment_id" binding:"required"`
	ProviderEventID   string  `json:"provider_event_id" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
	PaymentStatus     string  `json:"payment_status" binding:"required,oneof=success failed"`
	ExternalReference string  `json:"external_reference"`
}

type VerificationRequest struct {
	PaymentID          string  `json:"payment_id" binding:"required"`
	VerificationStatus string  `json:"verification_status" binding:"required,oneof=approved rejected"`
	Amount             float64 `json:"amount" binding:"required"`
	RejectionReason    *string `json:"rejection_reason"`
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
	Code           string   `json:"code" binding:"required"`
	Name           string   `json:"name" binding:"required"`
	Description    *string  `json:"description"`
	DiscountType   string   `json:"discount_type" binding:"required,oneof=percentage fixed"`
	DiscountValue  float64  `json:"discount_value" binding:"required"`
	MaxAmount      *float64 `json:"max_amount"`
	StartDate      *string  `json:"start_date"` // YYYY-MM-DD
	EndDate        *string  `json:"end_date"`   // YYYY-MM-DD
	IsActive       *bool    `json:"is_active"`
}

type ScholarshipUpdateRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	DiscountType  *string  `json:"discount_type"`
	DiscountValue *float64 `json:"discount_value"`
	MaxAmount     *float64 `json:"max_amount"`
	StartDate     *string  `json:"start_date"`
	EndDate       *string  `json:"end_date"`
	IsActive      *bool    `json:"is_active"`
}

// ============ FinanceHandler Struct ============

type FinanceHandler struct {
	repo           *repository.FinanceRepository
	db             *gorm.DB
	InvoiceService *service.InvoiceService
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


