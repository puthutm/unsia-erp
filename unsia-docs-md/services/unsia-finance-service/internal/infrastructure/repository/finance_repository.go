package repository

import (
	"errors"
	"time"

	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"gorm.io/gorm"
)

type FinanceRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

// Invoice operations
func (r *FinanceRepository) CreateInvoice(inv *domain.Invoice) error {
	return r.db.Create(inv).Error
}

func (r *FinanceRepository) GetInvoiceByID(id string) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.db.Preload("Items").Preload("Payments").Where("id = ?", id).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

func (r *FinanceRepository) GetInvoiceByNumber(num string) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.db.Preload("Items").Where("invoice_number = ?", num).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

// Payment operations
func (r *FinanceRepository) CreatePayment(pay *domain.Payment) error {
	return r.db.Create(pay).Error
}

func (r *FinanceRepository) GetPaymentByID(id string) (*domain.Payment, error) {
	var pay domain.Payment
	err := r.db.Where("id = ?", id).First(&pay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pay, nil
}

func (r *FinanceRepository) UpdatePaymentStatus(id string, status string, paidAt *time.Time, extRef *string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	if extRef != nil {
		updates["external_reference"] = extRef
	}
	return r.db.Model(&domain.Payment{}).Where("id = ?", id).Updates(updates).Error
}

// Callback operations
func (r *FinanceRepository) CreateCallback(cb *domain.PaymentGatewayCallback) error {
	return r.db.Create(cb).Error
}

// Verification operations
func (r *FinanceRepository) CreateVerification(ver *domain.PaymentVerification) error {
	return r.db.Create(ver).Error
}

// Clearance operations
func (r *FinanceRepository) GetStudentClearance(studentID string, periodID string, scope string) (*domain.StudentClearance, error) {
	var cl domain.StudentClearance
	err := r.db.Where("student_id = ? AND academic_period_id = ? AND service_scope = ?", studentID, periodID, scope).First(&cl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cl, nil
}

func (r *FinanceRepository) CreateStudentClearance(cl *domain.StudentClearance) error {
	return r.db.Create(cl).Error
}

func (r *FinanceRepository) UpdateStudentClearance(id string, status string, reason *string, updaterID string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_by": updaterID,
		"updated_at": time.Now(),
	}
	if reason != nil {
		updates["reason"] = reason
	}
	return r.db.Model(&domain.StudentClearance{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinanceRepository) CreateClearancePolicy(cp *domain.ClearancePolicy) error {
	return r.db.Create(cp).Error
}

func (r *FinanceRepository) GetClearancePolicyByID(id string) (*domain.ClearancePolicy, error) {
	var cp domain.ClearancePolicy
	err := r.db.Where("id = ?", id).First(&cp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cp, nil
}

func (r *FinanceRepository) UpdateClearancePolicy(cp *domain.ClearancePolicy) error {
	return r.db.Save(cp).Error
}

func (r *FinanceRepository) CreateInstallmentRequest(ir *domain.InstallmentRequest) error {
	return r.db.Create(ir).Error
}

func (r *FinanceRepository) GetInstallmentRequestByID(id string) (*domain.InstallmentRequest, error) {
	var ir domain.InstallmentRequest
	err := r.db.Where("id = ?", id).First(&ir).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ir, nil
}

func (r *FinanceRepository) UpdateInstallmentRequestStatus(id string, status string, approverID string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":      status,
		"approved_by": &approverID,
		"approved_at": &now,
	}
	return r.db.Model(&domain.InstallmentRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinanceRepository) CreateJournal(j *domain.Journal) error {
	return r.db.Create(j).Error
}

func (r *FinanceRepository) CreateJournalEntry(je *domain.JournalEntry) error {
	return r.db.Create(je).Error
}

func (r *FinanceRepository) GetCoaAccountByCode(code string) (*domain.CoaAccount, error) {
	var acc domain.CoaAccount
	err := r.db.Where("account_code = ?", code).First(&acc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

// Invoice List with Filters and Pagination
type InvoiceListFilter struct {
	Status           string
	TargetType       string
	StudentID       string
	ApplicantID     string
	AcademicPeriodID string
	DueDateFrom     string
	DueDateTo       string
	Search          string
	Page           int
	Limit          int
}

type PaginatedResult struct {
	Data       interface{} `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int64     `json:"total_pages"`
}

func (r *FinanceRepository) GetInvoices(filter InvoiceListFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Invoice{}).Preload("Items")

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.TargetType != "" {
		query = query.Where("target_type = ?", filter.TargetType)
	}
	if filter.StudentID != "" {
		query = query.Where("student_id = ?", filter.StudentID)
	}
	if filter.ApplicantID != "" {
		query = query.Where("applicant_id = ?", filter.ApplicantID)
	}
	if filter.AcademicPeriodID != "" {
		query = query.Where("academic_period_id = ?", filter.AcademicPeriodID)
	}
	if filter.DueDateFrom != "" {
		query = query.Where("due_date >= ?", filter.DueDateFrom)
	}
	if filter.DueDateTo != "" {
		query = query.Where("due_date <= ?", filter.DueDateTo)
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("invoice_number ILIKE ?", searchPattern)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Set default pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit
	var invoices []domain.Invoice
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&invoices).Error; err != nil {
		return nil, err
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	return &PaginatedResult{
		Data:       invoices,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// Payment List with Filters and Pagination
type PaymentListFilter struct {
	PaymentStatus  string
	PaymentMethod string
	InvoiceID    string
	StudentID    string
	DateFrom     string
	DateTo       string
	Search       string
	Page         int
	Limit        int
}

func (r *FinanceRepository) GetPayments(filter PaymentListFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Payment{}).Preload("Invoice")

	// Apply filters
	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}
	if filter.PaymentMethod != "" {
		query = query.Where("payment_method = ?", filter.PaymentMethod)
	}
	if filter.InvoiceID != "" {
		query = query.Where("invoice_id = ?", filter.InvoiceID)
	}
	if filter.StudentID != "" {
		query = query.Joins("JOIN invoices ON invoices.id = payments.invoice_id").Where("invoices.student_id = ?", filter.StudentID)
	}
	if filter.DateFrom != "" {
		query = query.Where("paid_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("paid_at <= ?", filter.DateTo)
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("payment_number ILIKE ? OR external_reference ILIKE ?", searchPattern, searchPattern)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Set default pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit
	var payments []domain.Payment
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&payments).Error; err != nil {
		return nil, err
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	return &PaginatedResult{
		Data:       payments,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// Clearance List with Filters
type ClearanceListFilter struct {
	Status           string
	StudentID       string
	AcademicPeriodID string
	ServiceScope    string
	Page           int
	Limit          int
}

func (r *FinanceRepository) GetClearances(filter ClearanceListFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.StudentClearance{})

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StudentID != "" {
		query = query.Where("student_id = ?", filter.StudentID)
	}
	if filter.AcademicPeriodID != "" {
		query = query.Where("academic_period_id = ?", filter.AcademicPeriodID)
	}
	if filter.ServiceScope != "" {
		query = query.Where("service_scope = ?", filter.ServiceScope)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Set default pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit
	var clearances []domain.StudentClearance
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&clearances).Error; err != nil {
		return nil, err
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

return &PaginatedResult{
		Data:       clearances,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// ============ Scholarships ============

// ScholarshipFilter for filtering scholarships
type ScholarshipFilter struct {
	IsActive bool
	Search  string
	Page    int
	Limit   int
}

func (r *FinanceRepository) CreateScholarship(s *domain.Scholarship) error {
	return r.db.Create(s).Error
}

func (r *FinanceRepository) GetScholarshipByID(id string) (*domain.Scholarship, error) {
	var scholarship domain.Scholarship
	err := r.db.Where("id = ?", id).First(&scholarship).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &scholarship, nil
}

func (r *FinanceRepository) GetScholarships(filter ScholarshipFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Scholarship{})

	if filter.IsActive {
		query = query.Where("is_active = ?", true)
	}
	if filter.Search != "" {
		query = query.Where("code ILIKE ? OR name ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var scholarships []domain.Scholarship
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&scholarships).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       scholarships,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdateScholarship(id string, updates map[string]interface{}) error {
	return r.db.Model(&domain.Scholarship{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinanceRepository) DeleteScholarship(id string) error {
	return r.db.Delete(&domain.Scholarship{}, "id = ?", id).Error
}

// ============ Cash Accounts ============

// CashAccountFilter for filtering cash accounts
type CashAccountFilter struct {
	AccountType string
	IsActive    bool
	Search     string
	Page       int
	Limit      int
}

func (r *FinanceRepository) CreateCashAccount(c *domain.CashAccount) error {
	return r.db.Create(c).Error
}

func (r *FinanceRepository) GetCashAccountByID(id string) (*domain.CashAccount, error) {
	var account domain.CashAccount
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *FinanceRepository) GetCashAccounts(filter CashAccountFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.CashAccount{})

	if filter.AccountType != "" {
		query = query.Where("account_type = ?", filter.AccountType)
	}
	if filter.IsActive {
		query = query.Where("is_active = ?", true)
	}
	if filter.Search != "" {
		query = query.Where("account_name ILIKE ? OR account_number ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var accounts []domain.CashAccount
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       accounts,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdateCashAccount(id string, updates map[string]interface{}) error {
	return r.db.Model(&domain.CashAccount{}).Where("id = ?", id).Updates(updates).Error
}

// Cash Mutations
func (r *FinanceRepository) CreateCashMutation(c *domain.CashMutation) error {
	return r.db.Create(c).Error
}

// CashMutationFilter for filtering cash mutations
type CashMutationFilter struct {
	Page  int
	Limit int
}

func (r *FinanceRepository) GetCashMutationsByAccountID(accountID string, filter CashMutationFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.CashMutation{}).Where("cash_account_id = ?", accountID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var mutations []domain.CashMutation
	err := query.Order("mutation_date DESC").Offset(offset).Limit(limit).Find(&mutations).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       mutations,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ============ Journals ============

// JournalFilter for filtering journals
type JournalFilter struct {
	SourceType string
	DateFrom string
	DateTo  string
	Search  string
	Page    int
	Limit   int
}

func (r *FinanceRepository) GetJournals(filter JournalFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Journal{})

	if filter.SourceType != "" {
		query = query.Where("source_type = ?", filter.SourceType)
	}
	if filter.DateFrom != "" {
		query = query.Where("journal_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("journal_date <= ?", filter.DateTo)
	}
	if filter.Search != "" {
		query = query.Where("journal_number ILIKE ? OR description ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var journals []domain.Journal
	err := query.Order("journal_date DESC").Offset(offset).Limit(limit).Find(&journals).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       journals,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) GetJournalByID(id string) (*domain.Journal, error) {
	var journal domain.Journal
	err := r.db.Where("id = ?", id).First(&journal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &journal, nil
}

func (r *FinanceRepository) GetJournalEntries(journalID string) ([]domain.JournalEntry, error) {
	var entries []domain.JournalEntry
	err := r.db.Where("journal_id = ?", journalID).Find(&entries).Error
	return entries, err
}

// ============ Budgets ============

// BudgetFilter for filtering budgets
type BudgetFilter struct {
	FiscalYear    int
	BudgetType   string
	IsActive     bool
	Search      string
	Page        int
	Limit       int
}

func (r *FinanceRepository) CreateBudget(b *domain.Budget) error {
	return r.db.Create(b).Error
}

func (r *FinanceRepository) GetBudgetByID(id string) (*domain.Budget, error) {
	var budget domain.Budget
	err := r.db.Where("id = ?", id).First(&budget).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

func (r *FinanceRepository) GetBudgets(filter BudgetFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Budget{})

	if filter.FiscalYear > 0 {
		query = query.Where("fiscal_year = ?", filter.FiscalYear)
	}
	if filter.BudgetType != "" {
		query = query.Where("budget_type = ?", filter.BudgetType)
	}
	if filter.IsActive {
		query = query.Where("is_active = ?", true)
	}
	if filter.Search != "" {
		query = query.Where("budget_number ILIKE ? OR budget_name ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var budgets []domain.Budget
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&budgets).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       budgets,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) GetBudgetItems(budgetID string) ([]domain.BudgetItem, error) {
	var items []domain.BudgetItem
	err := r.db.Where("budget_id = ?", budgetID).Find(&items).Error
	return items, err
}

func (r *FinanceRepository) CreateBudgetItem(bi *domain.BudgetItem) error {
	return r.db.Create(bi).Error
}

func (r *FinanceRepository) UpdateBudgetItem(id string, realization float64) error {
	return r.db.Model(&domain.BudgetItem{}).Where("id = ?", id).Update("realization", realization).Error
}

// ============ Vendors ============

// VendorFilter for filtering vendors
type VendorFilter struct {
	IsActive bool
	Search  string
	Page    int
	Limit   int
}

func (r *FinanceRepository) CreateVendor(v *domain.Vendor) error {
	return r.db.Create(v).Error
}

func (r *FinanceRepository) GetVendorByID(id string) (*domain.Vendor, error) {
	var vendor domain.Vendor
	err := r.db.Where("id = ?", id).First(&vendor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vendor, nil
}

func (r *FinanceRepository) GetVendors(filter VendorFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Vendor{})

	if filter.IsActive {
		query = query.Where("is_active = ?", true)
	}
	if filter.Search != "" {
		query = query.Where("vendor_code ILIKE ? OR vendor_name ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var vendors []domain.Vendor
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&vendors).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       vendors,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdateVendor(id string, updates map[string]interface{}) error {
	return r.db.Model(&domain.Vendor{}).Where("id = ?", id).Updates(updates).Error
}

// ============ Purchase Orders ============

// POFIlter for filtering purchase orders
type POFilter struct {
	Status   string
	VendorID string
	Search  string
	Page    int
	Limit   int
}

func (r *FinanceRepository) CreatePO(po *domain.PurchaseOrder) error {
	return r.db.Create(po).Error
}

func (r *FinanceRepository) GetPOByID(id string) (*domain.PurchaseOrder, error) {
	var po domain.PurchaseOrder
	err := r.db.Where("id = ?", id).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &po, nil
}

func (r *FinanceRepository) GetPOs(filter POFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.PurchaseOrder{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.VendorID != "" {
		query = query.Where("vendor_id = ?", filter.VendorID)
	}
	if filter.Search != "" {
		query = query.Where("po_number ILIKE ? OR description ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var pos []domain.PurchaseOrder
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       pos,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdatePOStatus(id string, status string, approverID string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if approverID != "" {
		now := time.Now()
		updates["approved_by"] = approverID
		updates["approved_at"] = &now
	}
	return r.db.Model(&domain.PurchaseOrder{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinanceRepository) GetPOItems(poID string) ([]domain.PurchaseOrderItem, error) {
	var items []domain.PurchaseOrderItem
	err := r.db.Where("po_id = ?", poID).Find(&items).Error
	return items, err
}

func (r *FinanceRepository) CreatePOItem(poi *domain.PurchaseOrderItem) error {
	return r.db.Create(poi).Error
}

// ============ Expense Events ============

// EventFilter for filtering expense events
type EventFilter struct {
	EventType string
	Status   string
	Page     int
	Limit    int
}

func (r *FinanceRepository) CreateExpenseEvent(e *domain.ExpenseEvent) error {
	return r.db.Create(e).Error
}

func (r *FinanceRepository) GetExpenseEventByID(id string) (*domain.ExpenseEvent, error) {
	var event domain.ExpenseEvent
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *FinanceRepository) GetExpenseEvents(filter EventFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.ExpenseEvent{})

	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var events []domain.ExpenseEvent
	err := query.Order("event_date DESC").Offset(offset).Limit(limit).Find(&events).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       events,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdateExpenseEvent(id string, updates map[string]interface{}) error {
	return r.db.Model(&domain.ExpenseEvent{}).Where("id = ?", id).Updates(updates).Error
}

// ============ Payroll Runs ============

// PayrollFilter for filtering payroll runs
type PayrollFilter struct {
	Status    string
	PeriodMonth int
	PeriodYear int
	Page      int
	Limit     int
}

func (r *FinanceRepository) CreatePayrollRun(pr *domain.PayrollRun) error {
	return r.db.Create(pr).Error
}

func (r *FinanceRepository) GetPayrollRunByID(id string) (*domain.PayrollRun, error) {
	var run domain.PayrollRun
	err := r.db.Where("id = ?", id).First(&run).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &run, nil
}

func (r *FinanceRepository) GetPayrollRuns(filter PayrollFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.PayrollRun{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.PeriodMonth > 0 {
		query = query.Where("period_month = ?", filter.PeriodMonth)
	}
	if filter.PeriodYear > 0 {
		query = query.Where("period_year = ?", filter.PeriodYear)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var runs []domain.PayrollRun
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&runs).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       runs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdatePayrollRunStatus(id string, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == "processed" {
		now := time.Now()
		updates["processed_at"] = &now
	}
	return r.db.Model(&domain.PayrollRun{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinanceRepository) CreatePayrollRunItem(pri *domain.PayrollRunItem) error {
	return r.db.Create(pri).Error
}

func (r *FinanceRepository) GetPayrollRunItems(runID string) ([]domain.PayrollRunItem, error) {
	var items []domain.PayrollRunItem
	err := r.db.Where("run_id = ?", runID).Find(&items).Error
	return items, err
}

// ============ Disbursements ============

// DisbursementFilter for filtering disbursements
type DisbursementFilter struct {
	Status       string
	ReferenceType string
	Page         int
	Limit        int
}

func (r *FinanceRepository) CreateDisbursement(d *domain.Disbursement) error {
	return r.db.Create(d).Error
}

func (r *FinanceRepository) GetDisbursementByID(id string) (*domain.Disbursement, error) {
	var disb domain.Disbursement
	err := r.db.Where("id = ?", id).First(&disb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &disb, nil
}

func (r *FinanceRepository) GetDisbursements(filter DisbursementFilter) (*PaginatedResult, error) {
	query := r.db.Model(&domain.Disbursement{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ReferenceType != "" {
		query = query.Where("reference_type = ?", filter.ReferenceType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var disbursements []domain.Disbursement
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&disbursements).Error
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &PaginatedResult{
		Data:       disbursements,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *FinanceRepository) UpdateDisbursementStatus(id string, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == "processed" {
		now := time.Now()
		updates["processed_at"] = &now
	}
	return r.db.Model(&domain.Disbursement{}).Where("id = ?", id).Updates(updates).Error
}

// ============ Reports ============

func (r *FinanceRepository) GenerateBalanceSheetReport(month int, year int) (*domain.BalanceSheetReport, error) {
	// Calculate assets (cash, accounts receivable, etc.)
	var totalAssets float64
	r.db.Model(&domain.CashAccount{}).Where("is_active = ?", true).Select("COALESCE(SUM(current_balance), 0)").Scan(&totalAssets)

	// Calculate liabilities (should be tracked separately)
	var totalLiabilities float64

	// Calculate equity (assets - liabilities)
	totalEquity := totalAssets - totalLiabilities

	report := &domain.BalanceSheetReport{
		ID:                 fmt.Sprintf("BS-%d-%d", year, month),
		ReportDate:        time.Now(),
		PeriodMonth:      month,
		PeriodYear:       year,
		ReportType:       "balance_sheet",
		TotalAssets:      totalAssets,
		TotalLiabilities: totalLiabilities,
		TotalEquity:      totalEquity,
		GeneratedAt:     time.Now(),
	}

	err := r.db.Create(report).Error
	return report, err
}

func (r *FinanceRepository) GenerateIncomeStatementReport(startDate, endDate time.Time) (*domain.IncomeStatementReport, error) {
	// Calculate revenue (total payments received in period)
	var totalRevenue float64
	r.db.Model(&domain.Payment{}).Where("payment_status = ? AND paid_at BETWEEN ? AND ?", "success", startDate, endDate).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	// Calculate expenses (total journal entries with expense accounts)
	var totalExpenses float64
	r.db.Model(&domain.JournalEntry{}).Where("debit > 0 AND created_at BETWEEN ? AND ?", startDate, endDate).Select("COALESCE(SUM(debit), 0)").Scan(&totalExpenses)

	netIncome := totalRevenue - totalExpenses

	report := &domain.IncomeStatementReport{
		ID:              fmt.Sprintf("IS-%s", endDate.Format("200601")),
		ReportDate:     time.Now(),
		PeriodStart:    startDate,
		PeriodEnd:      endDate,
		ReportType:    "income_statement",
		TotalRevenue:  totalRevenue,
		TotalExpenses: totalExpenses,
		NetIncome:     netIncome,
		GeneratedAt:  time.Now(),
	}

	err := r.db.Create(report).Error
	return report, err
}

func (r *FinanceRepository) GenerateCashFlowReport(startDate, endDate time.Time) (*domain.CashFlowReport, error) {
	// Calculate cash in (debits to cash accounts)
	var cashIn float64
	r.db.Model(&domain.CashMutation{}).Where("mutation_type = ? AND mutation_date BETWEEN ? AND ?", "debit", startDate, endDate).Select("COALESCE(SUM(amount), 0)").Scan(&cashIn)

	// Calculate cash out (credits from cash accounts)
	var cashOut float64
	r.db.Model(&domain.CashMutation{}).Where("mutation_type = ? AND mutation_date BETWEEN ? AND ?", "credit", startDate, endDate).Select("COALESCE(SUM(amount), 0)").Scan(&cashOut)

	netChange := cashIn - cashOut

	report := &domain.CashFlowReport{
		ID:         fmt.Sprintf("CF-%s", endDate.Format("200601")),
		ReportDate: time.Now(),
		PeriodStart: startDate,
		PeriodEnd:  endDate,
		ReportType: "cash_flow",
		CashIn:    cashIn,
		CashOut:   cashOut,
		NetChange: netChange,
	}

	err := r.db.Create(report).Error
	return report, err
}
