# Plan: Match UI Admin Keuangan dengan Backend Finance Service

## 📋 Overview

**Tugas**: Menambahkan endpoint API yang masih missing di backend untuk match dengan UI Admin Keuangan (ADMIN KEUANGAN.html) di folder `UI/KEUANGAN/`

## 📊 Gap Analysis: UI Features vs Current Backend

### Yang SUDAH Ada di Backend:
| Endpoint | UI Panel | Status |
|----------|---------|--------|
| `POST /api/v1/finance/invoices` | Tagihan Kuliah (Create) | ✅ |
| `GET /api/v1/finance/invoices/:id` | Tagihan Kuliah (Detail) | ✅ |
| `POST /api/v1/finance/payment-callbacks/:provider` | Pembayaran Masuk | ✅ |
| `POST /api/v1/finance/payment-verifications` | Pembayaran Masuk (Verify) | ✅ |
| `GET /api/v1/finance/clearances` | Dashboard | ✅ |
| `POST /api/v1/finance/clearance-policies` | Settings | ✅ |
| `PUT /api/v1/finance/clearance-policies/:id` | Settings | ✅ |
| `POST /api/v1/finance/installment-requests` | Tagihan Kuliah (Cicilan) | ✅ |
| `POST /api/v1/finance/installment-requests/:id/approve` | Tagihan Kuliah (Approve Cicilan) | ✅ |

### Yang BELUM ADA (Need Implementation):

| Priority | Endpoint | UI Panel | Keterangan |
|----------|----------|---------|-----------|
| **P0** | `GET /api/v1/finance/invoices` (LIST) | Tagihan Kuliah | List with filter & pagination |
| **P0** | `GET /api/v1/finance/payments` (LIST) | Pembayaran Masuk | List all payments |
| **P0** | `GET /api/v1/finance/clearances` (LIST) | Data Mahasiswa | List clearances |
| **P1** | `GET/POST /api/v1/finance/scholarships` | Beasiswa | CRUD Scholarships |
| **P1** | `GET/POST /api/v1/finance/cash-accounts` | Kas & Bank | Bank accounts |
| **P1** | `GET/POST /api/v1/finance/coa-accounts` | Jurnal & Buku Besar | Chart of Accounts |
| **P1** | `GET /api/v1/finance/journals` | Jurnal & Buku Besar | Journal entries |
| **P2** | `GET/POST /api/v1/finance/budgets` | Anggaran (RAB) | Budget management |
| **P2** | `GET /api/v1/finance/reports` | Laporan Keuangan | Report generation |
| **P2** | `GET/POST /api/v1/finance/vendors` | Belanja & Operasional | Vendor management |
| **P2** | `GET/POST /api/v1/finance/purchase-orders` | Belanja & Operasional | PO management |
| **P2** | `GET/POST /api/v1/finance/payroll-runs` | Payroll Disbursement | Payroll runs |
| **P2** | `GET/POST /api/v1/finance/disbursements` | Disbursement CRM | Commission disbursement |
| **P2** | `GET/POST /api/v1/finance/events` | Wisuda & Kegiatan | Event management |

---

## 🎯 Implementation Plan

### Phase 1: Invoice & Payment Enhancement (P0)
**Duration**: 2 days

#### 1.1 Get Invoice List with Filtering
```go
// GET /api/v1/finance/invoices
// Query params:
//   - status: unpaid, partially_paid, paid, cancelled, expired
//   - target_type: applicant, student
//   - student_id: filter by student
//   - applicant_id: filter by applicant
//   - academic_period_id: filter by period
//   - due_date_from, due_date_to: filter by due date range
//   - page, limit: pagination
//   - search: search by invoice_number or student name
```

**Files to create:**
- `internal/handler/invoice_list_handler.go` - Add GetInvoices function

**Files to update:**
- `internal/handler/finance_handler.go` - Add route
- `cmd/finance-service/main.go` - Add route registration

#### 1.2 Get Payments List
```go
// GET /api/v1/finance/payments
// Query params:
//   - payment_status: pending, success, failed
//   - payment_method: va_bni, va_bca, va_mandiri, qris, credit_card, manual
//   - invoice_id: filter by invoice
//   - student_id: filter by student
//   - date_from, date_to: filter by payment date
//   - page, limit: pagination
```

**Files to create:**
- `internal/handler/payment_list_handler.go`

#### 1.3 Get Clearances List
```go
// GET /api/v1/finance/clearances
// Query params:
//   - status: blocked, conditional, cleared, revoked
//   - student_id: filter by student
//   - academic_period_id: filter by period
//   - service_scope: registration, krs, transcript, etc
//   - page, limit: pagination
```

**Files to update:**
- `internal/handler/finance_handler.go` - Add GetClearances endpoint

---

### Phase 2: Scholarships (P1)
**Duration**: 1 day

#### 2.1 Scholarships CRUD
```go
// GET /api/v1/finance/scholarships
// POST /api/v1/finance/scholarships
// GET /api/v1/finance/scholarships/:id
// PUT /api/v1/finance/scholarships/:id
// DELETE /api/v1/finance/scholarships/:id
// POST /api/v1/finance/scholarships/:id/approve
```

**Files to create:**
- `internal/domain/models.go` - Add Scholarship model
- `internal/handler/scholarship_handler.go`
- `internal/service/scholarship_service.go`
- `migrations/` - Add scholarship table

---

### Phase 3: Cash & Bank Management (P1)
**Duration**: 1 day

#### 3.1 Cash Accounts
```go
// GET /api/v1/finance/cash-accounts
// POST /api/v1/finance/cash-accounts
// GET /api/v1/finance/cash-accounts/:id
// POST /api/v1/finance/cash-accounts/:id/mutations
```

**Files to create:**
- `internal/domain/models.go` - Add CashAccount model
- `internal/handler/cash_account_handler.go`

#### 3.2 COA Management
```go
// GET /api/v1/finance/coa-accounts
// POST /api/v1/finance/coa-accounts
// GET /api/v1/finance/coa-accounts/:id
```

**Files to update:**
- `internal/domain/models.go` - Add COA endpoints (already has CoaAccount)

#### 3.3 Journals
```go
// GET /api/v1/finance/journals
// GET /api/v1/finance/journals/:id
// GET /api/v1/finance/journals/:id/entries
```

**Files to create:**
- `internal/handler/journal_handler.go`

---

### Phase 4: Budget & Reports (P2)
**Duration**: 2 days

#### 4.1 Budget Management
```go
// GET /api/v1/finance/budgets
// POST /api/v1/finance/budgets
// GET /api/v1/finance/budgets/:id
// GET /api/v1/finance/budgets/:id/realization
```

#### 4.2 Reports
```go
// GET /api/v1/finance/reports/position (Neraca)
// GET /api/v1/finance/reports/activity (L/R)
// GET /api/v1/finance/reports/cashflow (Arus Kas)
// GET /api/v1/finance/reports/aging (Aging Piutang)
// GET /api/v1/finance/reports/budget-realization (Anggaran vs Realisasi)
```

**Files to create:**
- `internal/handler/report_handler.go`

---

### Phase 5: Vendor & PO (P2)
**Duration**: 1 day

#### 5.1 Vendors
```go
// GET /api/v1/finance/vendors
// POST /api/v1/finance/vendors
// GET /api/v1/finance/vendors/:id
```

#### 5.2 Purchase Orders
```go
// GET /api/v1/finance/purchase-orders
// POST /api/v1/finance/purchase-orders
// GET /api/v1/finance/purchase-orders/:id
// POST /api/v1/finance/purchase-orders/:id/approve
```

---

### Phase 6: Payroll & Disbursement (P2)
**Duration**: 2 days

#### 6.1 Payroll Runs
```go
// GET /api/v1/finance/payroll-runs
// POST /api/v1/finance/payroll-runs (from HRIS sync)
// POST /api/v1/finance/payroll-runs/:id/approve
// GET /api/v1/finance/payroll-runs/:id/export
```

#### 6.2 Disbursements
```go
// GET /api/v1/finance/disbursements
// POST /api/v1/finance/disbursements
// POST /api/v1/finance/disbursements/:id/approve
// POST /api/v1/finance/disbursements/:id/process
```

---

### Phase 7: Events (P2)
**Duration**: 1 day

#### 7.1 Event Management
```go
// GET /api/v1/finance/events
// POST /api/v1/finance/events
// GET /api/v1/finance/events/:id
// POST /api/v1/finance/events/:id/generate-invoices
```

---

## 📦 File Structure Target

```
unsia-finance-service/
├── cmd/finance-service/
│   ├── main.go
│   └── events/
│       └── consumer.go (EXISTING ���)
���── internal/
│   ├── domain/
│   │   ├── models.go (NEED UPDATE - add more models)
│   │   └── validators.go
│   ├── service/
│   │   ├── invoice_service.go (EXISTING ✓)
│   │   ├── payment_gateway_service.go (EXISTING ✓)
│   │   ├── scholarship_service.go (NEW)
│   │   ├── cash_account_service.go (NEW)
│   │   ├── journal_service.go (NEW)
│   │   ├── budget_service.go (NEW)
│   │   ├── report_service.go (NEW)
│   │   ├── vendor_service.go (NEW)
│   │   └── payroll_service.go (NEW)
│   ├── handler/
│   │   ├── finance_handler.go (EXISTING ✓)
│   │   ├── health_handler.go (EXISTING ✓)
│   │   ├── invoice_list_handler.go (NEW)
│   │   ├── payment_list_handler.go (NEW)
│   │   ├── scholarship_handler.go (NEW)
│   │   ├── cash_account_handler.go (NEW)
│   │   ├── journal_handler.go (NEW)
│   │   ├── budget_handler.go (NEW)
│   │   ├── report_handler.go (NEW)
│   │   ├── vendor_handler.go (NEW)
│   │   └── payroll_handler.go (NEW)
│   ├── middleware/
│   │   └── auth_middleware.go (EXISTING ✓)
│   ├── state_machine/
│   │   └── invoice.go (EXISTING ✓)
│   └── infrastructure/
│       ├── repository/
│       │   └── finance_repository.go (EXISTING ✓)
│       └── outbox_poller.go (EXISTING ✓)
├── migrations/
└── go.mod
```

---

## ✅ Priority Checklist

### P0 - Critical (Week 1)
- [x] GET /api/v1/finance/invoices (list + filter + pagination)
- [x] GET /api/v1/finance/payments (list + filter)
- [x] GET /api/v1/finance/clearances (list)

### P1 - Important (Week 2)
- [x] Scholarships CRUD
- [x] Cash Accounts & Mutations
- [x] COA & Journals

### P2 - Nice to Have (Week 3-4)
- [x] Budget Management
- [ ] Reports (Neraca, L/R, Arus Kas)
- [x] Vendors & Purchase Orders
- [ ] Payroll Runs
- [ ] Disbursements
- [x] Event Management

---

## 📅 Timeline Estimate

| Phase | Duration | Cumulative |
|-------|----------|-----------|
| Phase 1: Invoice & Payment | 2 days | Day 1-2 |
| Phase 2: Scholarships | 1 day | Day 3 |
| Phase 3: Cash & Bank | 1 day | Day 4 |
| Phase 4: Budget & Reports | 2 days | Day 5-6 |
| Phase 5: Vendor & PO | 1 day | Day 7 |
| Phase 6: Payroll | 2 days | Day 8-9 |
| Phase 7: Events | 1 day | Day 10 |
| **Buffer & Testing** | 2 days | Day 11-12 |

**Total Estimate**: ~12 days (2-3 sprints)

---

*Plan dibuat berdasarkan analisis gap antara UI ADMIN KEUANGAN.html dengan backend unsia-finance-service yang sudah ada*
