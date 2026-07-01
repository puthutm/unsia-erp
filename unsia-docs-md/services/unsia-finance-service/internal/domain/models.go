package domain

import (
	"time"
)

type Invoice struct {
	ID               string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	InvoiceNumber    string        `gorm:"column:invoice_number;unique;not null"`
	TargetType       string        `gorm:"column:target_type;not null"` // applicant, student
	ApplicantID      *string       `gorm:"column:applicant_id"`         // external_ref: pmb.applicants.id
	StudentID        *string       `gorm:"column:student_id"`           // external_ref: academic.students.id
	AcademicPeriodID *string       `gorm:"column:academic_period_id"`   // external_ref: ref.academic_periods.id
	TotalAmount      float64       `gorm:"column:total_amount;default:0.00;not null"`
	PaidAmount       float64       `gorm:"column:paid_amount;default:0.00;not null"`
	Status           string        `gorm:"column:status;default:'DRAFT';not null"` // DRAFT, ISSUED, PARTIALLY_PAID, PAID, CANCELLED, EXPIRED
	DueDate          *time.Time    `gorm:"column:due_date"`
	SourceModule     *string       `gorm:"column:source_module"`
	SourceRefID      *string       `gorm:"column:source_ref_id"`
	CreatedAt        time.Time     `gorm:"column:created_at"`
	UpdatedAt        time.Time     `gorm:"column:updated_at"`
	Items            []InvoiceItem `gorm:"foreignKey:InvoiceID;references:ID"`
	Payments         []Payment     `gorm:"foreignKey:InvoiceID;references:ID"`
}

func (Invoice) TableName() string {
	return "invoices"
}

type InvoiceItem struct {
	ID                 string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	InvoiceID          string  `gorm:"column:invoice_id;not null"`
	PaymentComponentID *string `gorm:"column:payment_component_id"` // external_ref: ref.payment_components.id
	Description        string  `gorm:"column:description"`
	Amount             float64 `gorm:"column:amount;default:0.00;not null"`
	DiscountAmount     float64 `gorm:"column:discount_amount;default:0.00;not null"`
	FinalAmount        float64 `gorm:"column:final_amount;default:0.00;not null"`
}

func (InvoiceItem) TableName() string {
	return "invoice_items"
}

type Payment struct {
	ID                string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	InvoiceID         string     `gorm:"column:invoice_id;not null"`
	PaymentMethodID   *string    `gorm:"column:payment_method_id"` // external_ref: ref.payment_methods.id
	PaymentNumber     *string    `gorm:"column:payment_number;unique"`
	Amount            float64    `gorm:"column:amount;default:0.00;not null"`
	PaymentStatus     string     `gorm:"column:payment_status;default:'RECEIVED';not null"` // RECEIVED, VERIFIED, POSTED, FAILED, REVERSED
	PaidAt            *time.Time `gorm:"column:paid_at"`
	ExternalReference *string    `gorm:"column:external_reference"`
	IdempotencyKey    *string    `gorm:"column:idempotency_key;unique"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentGatewayCallback struct {
	ID                string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PaymentID         *string    `gorm:"column:payment_id"`
	Provider          string     `gorm:"column:provider;not null"`
	ProviderEventID   *string    `gorm:"column:provider_event_id"`
	ExternalReference *string    `gorm:"column:external_reference"`
	IdempotencyKey    *string    `gorm:"column:idempotency_key;unique"`
	Payload           string     `gorm:"type:jsonb;column:payload"`
	SignatureValid    bool       `gorm:"column:signature_valid;default:false;not null"`
	CallbackStatus    string     `gorm:"column:callback_status;default:'received';not null"` // received, processed, failed
	ReceivedAt        time.Time  `gorm:"column:received_at;default:now()"`
	ProcessedAt       *time.Time `gorm:"column:processed_at"`
}

func (PaymentGatewayCallback) TableName() string {
	return "payment_gateway_callbacks"
}

type PaymentVerification struct {
	ID                 string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PaymentID          string     `gorm:"column:payment_id;not null"`
	VerifiedBy         *string    `gorm:"column:verified_by"` // external_ref: core.users.id
	VerificationStatus string     `gorm:"column:verification_status;default:'pending';not null"` // pending, approved, rejected
	RejectionReason    *string    `gorm:"column:rejection_reason"`
	Note               *string    `gorm:"column:note"`
	VerifiedAt         *time.Time `gorm:"column:verified_at"`
}

func (PaymentVerification) TableName() string {
	return "payment_verifications"
}

type StudentClearance struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID       string     `gorm:"column:student_id;not null"` // external_ref: academic.students.id
	AcademicPeriodID *string    `gorm:"column:academic_period_id"`  // external_ref: ref.academic_periods.id
	ServiceScope    string     `gorm:"column:service_scope;not null"` // registration, krs, graduation
	Status          string     `gorm:"column:status;default:'BLOCKED';not null"` // BLOCKED, CONDITIONAL, CLEARED, REVOKED
	Reason          *string    `gorm:"column:reason"`
	ValidUntil      *time.Time `gorm:"column:valid_until"`
	UpdatedBy       *string    `gorm:"column:updated_by"` // external_ref: core.users.id
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
}

func (StudentClearance) TableName() string {
	return "student_clearances"
}

type ClearancePolicy struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code         string    `gorm:"column:code;unique;not null"`
	Name         string    `gorm:"column:name;not null"`
	ServiceScope string    `gorm:"column:service_scope"`
	RuleJson     string    `gorm:"type:jsonb;column:rule_json"`
	IsActive     bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (ClearancePolicy) TableName() string {
	return "clearance_policies"
}

type InstallmentRequest struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	InvoiceID   string     `gorm:"column:invoice_id;not null"`
	StudentID   *string    `gorm:"column:student_id"`
	Status      string     `gorm:"column:status;default:'pending';not null"`
	Reason      string     `gorm:"column:reason"`
	RequestedAt time.Time  `gorm:"column:requested_at;default:now()"`
	ApprovedBy  *string    `gorm:"column:approved_by"`
	ApprovedAt  *time.Time `gorm:"column:approved_at"`
}

func (InstallmentRequest) TableName() string {
	return "installment_requests"
}

type CoaAccount struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AccountCode   string    `gorm:"column:account_code;unique;not null"`
	AccountName   string    `gorm:"column:account_name;not null"`
	NormalBalance string    `gorm:"column:normal_balance;not null"`
	IsActive      bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (CoaAccount) TableName() string {
	return "coa_accounts"
}

type Journal struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	JournalNumber string    `gorm:"column:journal_number;unique;not null"`
	JournalDate   time.Time `gorm:"column:journal_date;type:date;not null"`
	SourceType    string    `gorm:"column:source_type"`
	SourceID      *string   `gorm:"column:source_id"`
	Description   string    `gorm:"column:description"`
	CreatedBy     *string   `gorm:"column:created_by"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (Journal) TableName() string {
	return "journals"
}

type JournalEntry struct {
	ID           string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	JournalID    string  `gorm:"column:journal_id;not null"`
	CoaAccountID *string `gorm:"column:coa_account_id"`
	Debit        float64 `gorm:"column:debit;default:0.00;not null"`
	Credit       float64 `gorm:"column:credit;default:0.00;not null"`
	Description  string  `gorm:"column:description"`
}

func (JournalEntry) TableName() string {
	return "journal_entries"
}

// InboxEvent represents events received from other services
type InboxEvent struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EventName     string     `gorm:"column:event_name"`
	EventKey     string     `gorm:"column:event_key"`
	CorrelationID string    `gorm:"column:correlation_id"`
	Payload     string    `gorm:"type:jsonb;column:payload"`
	Status      string    `gorm:"column:status"` // received, processed, failed, duplicate
	ErrorMessage *string   `gorm:"column:error_message"`
	ReceivedAt  time.Time `gorm:"column:received_at"`
	ProcessedAt *time.Time `gorm:"column:processed_at"`
}

func (InboxEvent) TableName() string {
	return "inbox_events"
}

// OutboxEvent represents events to be published
type OutboxEvent struct {
	ID               string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EventName        string     `gorm:"column:event_name"`
	EventVersion    string     `gorm:"column:event_version"`
	EventKey        string     `gorm:"column:event_key;uniqueIndex"`
	PublisherService string   `gorm:"column:publisher_service"`
	AggregateType  string     `gorm:"column:aggregate_type"`
	AggregateID   string     `gorm:"column:aggregate_id"`
	CorrelationID  string     `gorm:"column:correlation_id"`
	CausationID   string     `gorm:"column:causation_id"`
	Payload       string     `gorm:"type:jsonb;column:payload"`
	Status        string     `gorm:"column:status"` // pending, published, failed, dead_letter
	RetryCount   int        `gorm:"column:retry_count;default:0"`
	ErrorMessage *string   `gorm:"column:error_message"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

// IdempotencyKey stores idempotency keys
type IdempotencyKey struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Key       string   `gorm:"column:key;uniqueIndex"`
	Response  string   `gorm:"column:response"`
	Status   string   `gorm:"column:status"` // processing, completed, failed
	EntityType string   `gorm:"column:entity_type"`
	EntityID string   `gorm:"column:entity_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	ExpiresAt *time.Time `gorm:"column:expires_at"`
}

func (IdempotencyKey) TableName() string {
	return "idempotency_keys"
}

// Scholarship represents a scholarship/discount for students
type Scholarship struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID       string     `gorm:"column:student_id;not null"` // external_ref: academic.students.id
	ScholarshipType string     `gorm:"column:scholarship_type"`
	Amount          float64    `gorm:"column:amount;default:0.00;not null"`
	Status          string     `gorm:"column:status;default:'PENDING_APPROVAL';not null"` // PENDING_APPROVAL, APPROVED, etc
	ApprovedBy      *string    `gorm:"column:approved_by"` // external_ref: core.users.id
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (Scholarship) TableName() string {
	return "scholarships"
}

// CashAccount represents bank accounts and cash on hand
type CashAccount struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AccountName   string    `gorm:"column:account_name;not null"`
	AccountNumber string    `gorm:"column:account_number"`
	AccountType   string    `gorm:"column:account_type;not null"` // bank, cash
	BankName      *string   `gorm:"column:bank_name"`
	Branch       *string   `gorm:"column:branch"`
	IsActive      bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (CashAccount) TableName() string {
	return "cash_accounts"
}

// CashTransaction represents cash/bank mutations (deposits, withdrawals)
type CashTransaction struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CashAccountID   string     `gorm:"column:cash_account_id;not null"`
	TransactionType string     `gorm:"column:transaction_type;not null"` // DEBIT, CREDIT
	SourceType      *string    `gorm:"column:source_type"`
	SourceID        *string    `gorm:"column:source_id"`
	Amount          float64    `gorm:"column:amount;default:0;not null"`
	Description     string     `gorm:"column:description"`
	TransactionAt   time.Time  `gorm:"column:transaction_at;default:now()"`
}

func (CashTransaction) TableName() string {
	return "cash_transactions"
}

// Budget represents budget for Rencana Anggaran Biaya (RAB)
type Budget struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	BudgetNumber    string    `gorm:"column:budget_number;unique;not null"`
	BudgetName      string    `gorm:"column:budget_name;not null"`
	AcademicPeriodID *string   `gorm:"column:academic_period_id"`
	FiscalYear     int       `gorm:"column:fiscal_year;not null"`
	BudgetType     string    `gorm:"column:budget_type;not null"` // operational, capital, investment
	TotalAmount   float64   `gorm:"column:total_amount;default:0;not null"`
	IsActive      bool      `gorm:"column:is_active;default:true;not null"`
	StartDate     *time.Time `gorm:"column:start_date"`
	EndDate       *time.Time `gorm:"column:end_date"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (Budget) TableName() string {
	return "budgets"
}

// BudgetLine represents individual budget line items
type BudgetLine struct {
	ID             string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	BudgetID       string   `gorm:"column:budget_id;not null"`
	CoaAccountID   *string  `gorm:"column:coa_account_id"`
	Description    string   `gorm:"column:description"`
	Amount         float64  `gorm:"column:amount;default:0;not null"`
	RealizedAmount float64  `gorm:"column:realized_amount;default:0;not null"`
}

func (BudgetLine) TableName() string {
	return "budget_lines"
}

// Vendor represents suppliers/vendors
type Vendor struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	VendorCode   string    `gorm:"column:vendor_code;unique;not null"`
	VendorName   string    `gorm:"column:vendor_name;not null"`
	ContactPerson *string  `gorm:"column:contact_person"`
	Phone       *string   `gorm:"column:phone"`
	Email       *string   `gorm:"column:email"`
	Address     *string   `gorm:"column:address"`
	TaxNumber   *string   `gorm:"column:tax_number"`
	IsActive    bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Vendor) TableName() string {
	return "vendors"
}

// PurchaseOrder represents PO for procurement
type PurchaseOrder struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PONumber     string    `gorm:"column:po_number;unique;not null"`
	VendorID    string    `gorm:"column:vendor_id;not null"`
	PODate      time.Time `gorm:"column:po_date;type:date;not null"`
	DeliveryDate *time.Time `gorm:"column:delivery_date"`
	Description string    `gorm:"column:description"`
	TotalAmount float64   `gorm:"column:total_amount;default:0;not null"`
	Status      string    `gorm:"column:status;default:'draft';not null"` // draft, approved, received, cancelled
	CreatedBy   *string   `gorm:"column:created_by"`
	ApprovedBy  *string   `gorm:"column:approved_by"`
	ApprovedAt *time.Time `gorm:"column:approved_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// PurchaseOrderItem represents line items in PO
type PurchaseOrderItem struct {
	ID             string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	POID           string   `gorm:"column:po_id;not null"`
	Description    string   `gorm:"column:description"`
	Quantity      float64  `gorm:"column:quantity;default:0;not null"`
	UnitPrice     float64  `gorm:"column:unit_price;default:0;not null"`
	TotalPrice    float64  `gorm:"column:total_price;default:0;not null"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

// ExpenseEvent represents graduation or other events
type ExpenseEvent struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EventName      string    `gorm:"column:event_name;not null"`
	EventType      string    `gorm:"column:event_type;not null"` // graduation, seminar, workshop
	EventDate      time.Time `gorm:"column:event_date;type:date;not null"`
	Description   *string   `gorm:"column:description"`
	BudgetAmount  float64   `gorm:"column:budget_amount;default:0;not null"`
	Status       string    `gorm:"column:status;default:'planned';not null"` // planned, approved, executed, cancelled
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (ExpenseEvent) TableName() string {
	return "expense_events"
}

// PayrollRun represents payroll runs synced from HRIS
type PayrollRun struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PayrollPeriod string     `gorm:"column:payroll_period;not null"`
	RunDate       time.Time  `gorm:"column:run_date;not null"`
	TotalAmount   float64    `gorm:"column:total_amount;default:0;not null"`
	Status        string     `gorm:"column:status;default:'draft';not null"`
	ApprovedBy    *string    `gorm:"column:approved_by"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (PayrollRun) TableName() string {
	return "payroll_runs"
}

// PayrollItem represents individual employee payroll
type PayrollItem struct {
	ID              string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PayrollRunID    string   `gorm:"column:payroll_run_id;not null"`
	EmployeeID      string   `gorm:"column:employee_id;not null"`
	GrossAmount     float64  `gorm:"column:gross_amount;default:0;not null"`
	DeductionAmount float64  `gorm:"column:deduction_amount;default:0;not null"`
	NetAmount       float64  `gorm:"column:net_amount;default:0;not null"`
	Status          string   `gorm:"column:status;default:'draft';not null"`
}

func (PayrollItem) TableName() string {
	return "payroll_items"
}

// Disbursement represents CRM commission disbursements
type Disbursement struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	DisbursementNumber string   `gorm:"column:disbursement_number;unique;not null"`
	ReferenceID    string    `gorm:"column:reference_id;not null"` // CRM lead/sale reference
	ReferenceType  string    `gorm:"column:reference_type;not null"` // commission, referral
	Amount       float64   `gorm:"column:amount;default:0;not null"`
	RecipientName string   `gorm:"column:recipient_name;not null"`
	BankAccount  *string  `gorm:"column:bank_account"`
	BankName     *string  `gorm:"column:bank_name"`
	Status       string   `gorm:"column:status;default:'pending';not null"` // pending, approved, processed, failed
	ProcessedAt  *time.Time `gorm:"column:processed_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (Disbursement) TableName() string {
	return "disbursements"
}

// ============ Report Models ============

// BalanceSheetReport represents Neraca
type BalanceSheetReport struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ReportDate   time.Time `gorm:"column:report_date;type:date;not null"`
	PeriodMonth int       `gorm:"column:period_month;not null"`
	PeriodYear  int       `gorm:"column:period_year;not null"`
	ReportType  string    `gorm:"column:report_type;default:'balance_sheet';not null"`
	TotalAssets float64   `gorm:"column:total_assets;default:0;not null"`
	TotalLiabilities float64 `gorm:"column:total_liabilities;default:0;not null"`
	TotalEquity float64   `gorm:"column:total_equity;default:0;not null"`
	GeneratedAt time.Time `gorm:"column:generated_at"`
}

func (BalanceSheetReport) TableName() string {
	return "balance_sheet_reports"
}

// IncomeStatementReport represents Laba/Rugi
type IncomeStatementReport struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ReportDate   time.Time `gorm:"column:report_date;type:date;not null"`
	PeriodStart time.Time `gorm:"column:period_start;type:date;not null"`
	PeriodEnd   time.Time `gorm:"column:period_end;type:date;not null"`
	ReportType  string    `gorm:"column:report_type;default:'income_statement';not null"`
	TotalRevenue float64   `gorm:"column:total_revenue;default:0;not null"`
	TotalExpenses float64  `gorm:"column:total_expenses;default:0;not null"`
	NetIncome  float64   `gorm:"column:net_income;default:0;not null"`
	GeneratedAt time.Time `gorm:"column:generated_at"`
}

func (IncomeStatementReport) TableName() string {
	return "income_statement_reports"
}

// CashFlowReport represents Arus Kas
type CashFlowReport struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ReportDate   time.Time `gorm:"column:report_date;type:date;not null"`
	PeriodStart time.Time `gorm:"column:period_start;type:date;not null"`
	PeriodEnd   time.Time `gorm:"column:period_end;type:date;not null"`
	ReportType  string    `gorm:"column:report_type;default:'cash_flow';not null"`
	CashIn     float64   `gorm:"column:cash_in;default:0;not null"`
	CashOut    float64   `gorm:"column:cash_out;default:0;not null"`
	NetChange  float64   `gorm:"column:net_change;default:0;not null"`
	GeneratedAt time.Time `gorm:"column:generated_at"`
}

func (CashFlowReport) TableName() string {
	return "cash_flow_reports"
}

// ClearanceDispensation represents temporary dispensation for clearance block
type ClearanceDispensation struct {
	ID                 string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentClearanceID string     `gorm:"column:student_clearance_id;not null"`
	Reason             string     `gorm:"column:reason"`
	ApprovedBy         *string    `gorm:"column:approved_by"`
	ApprovedAt         *time.Time `gorm:"column:approved_at"`
	ValidUntil         *time.Time `gorm:"column:valid_until;type:date"`
	Status             string     `gorm:"column:status;default:'pending';not null"`
}

func (ClearanceDispensation) TableName() string {
	return "clearance_dispensations"
}
