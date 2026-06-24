# Implementation Plan: unsia-finance-service

## 📋 Overview

**Service**: unsia-finance-service  
**Stack**: Go 1.22+ · Gin · GORM + golang-migrate  
**Port**: `:8005`  
**Database**: `finance_db` (PostgreSQL)

## 📝 Current State Analysis

### ✅ Already Implemented

1. **Domain Models** (`internal/domain/models.go`) - 17 entities
   - Invoice, InvoiceItem, Payment, PaymentGatewayCallback, PaymentVerification
   - StudentClearance, ClearancePolicy
   - InstallmentRequest
   - CoaAccount, Journal, JournalEntry

2. **Handler** (`internal/handler/finance_handler.go`)
   - `POST /api/v1/finance/invoices` ✅
   - `GET /api/v1/finance/invoices/:id` ✅
   - `POST /api/v1/finance/payment-callbacks/:provider` ✅
   - `POST /api/v1/finance/payment-verifications` ✅
   - `GET /api/v1/finance/clearances` ✅
   - `POST /api/v1/finance/clearance-policies` ✅
   - `PUT /api/v1/finance/clearance-policies/:id` ✅
   - `POST /api/v1/finance/installment-requests` ✅
   - `POST /api/v1/finance/installment-requests/:id/approve` ✅
   - Journal recording (auto on payment) ✅

3. **Database Migrations** (26 migrations) - Complete schema ✓

### ❌ Missing / Need Enhancement

Based on requirements.md analysis:

| Req # | Requirement | Status | Priority | Gap |
|--------|------------|--------|----------|-----|
| 1 | JWT RS256 + JWKS Validation | Partially | P0 | Need middleware |
| 1 | RBAC Permission Check | ❌ | P0 | Need middleware |
| 1 | X-Application-Code, X-Active-Role headers | ❌ | P0 | Need middleware |
| 2 | Get Invoice List (filter, pagination) | ❌ | P0 | Need endpoint |
| 2 | Idempotency-Key support | ❌ | P0 | Need middleware |
| 3 | State Machine Validation | Incomplete | P0 | Need state machine |
| 4 | Callback signature validation | ❌ | P0 | Need validation |
| 6 | Payment State Machine | ❌ | P0 | Need validation |
| 8 | Clearance State Machine | ❌ | P0 | Need validation |
| 11 | Inbound Event Processing | ❌ | P1 | Need event consumer |
| 12 | Outbox Event Publishing | Partially | P1 | Need poller |
| 13-18 | Payroll, COA, Budget, Disbursement etc | ❌ | P2 | To be implemented |
| 19 | Success/Error Envelope Format | Already | - | Using shared-errorenvelope ✓ |
| 21 | Audit Logging | Partially | - | Using shared-audit ✓ |
| 22 | Health Check | ❌ | - | Need endpoint |

---

## 🎯 Implementation Phases

### Phase 1: Security & Middleware (P0)

> **Duration**: 2 days  
> **Goal**: Complete authentication flow

#### 1.1 JWT Middleware with JWKS
```go
// internal/middleware/auth_middleware.go
func JWTAuthMiddleware() gin.HandlerFunc {
    // 1. Validate JWT RS256 using JWKS from core-service
    // 2. Check required headers (X-Application-Code, X-Active-Role, X-Correlation-Id)
    // 3. Cache JWKS with TTL (min 5 minutes)
    // 4. Return 401/400 appropriately
}
```

#### 1.2 RBAC Middleware
```go
// internal/middleware/rbac_middleware.go
func RequirePermission(permissions ...string) gin.HandlerFunc {
    // 1. Extract permissions from JWT claims
    // 2. Check if user has required permission
    // 3. Return 403 if denied
}
```

**Files to create:**
- `internal/middleware/auth_middleware.go`
- `internal/middleware/rbac_middleware.go`

**Files to update:**
- `cmd/finance-service/main.go` - Apply middleware to routes
- `internal/handler/finance_handler.go` - Add actor tracking

---

### Phase 2: Invoice Management (P0)

> **Duration**: 2 days  
> **Goal**: Complete invoice CRUD + state machine

#### 2.1 Get Invoice List
```go
// GET /api/v1/finance/invoices
// Query params: status, target_type, applicant_id, student_id, academic_period_id, 
//             due_date_from, due_date_to, page, limit
```

#### 2.2 State Machine Enforcement
```go
// Valid transitions:
// DRAFT → ISSUED
// ISSUED → PARTIALLY_PAID, PAID, CANCELLED, EXPIRED
// PARTIALLY_PAID → PAID, CANCELLED, EXPIRED
// PAID, CANCELLED → (no further transitions)
```

#### 2.3 Idempotency Support
```go
// Check X-Idempotency-Key header
// Store in idempotency_keys table
// Return cached response on duplicate
```

**Files to update:**
- `internal/handler/finance_handler.go` - Add GetInvoices, UpdateInvoiceStatus
- `internal/domain/models.go` - Add IdempotencyKey model

**Files to create:**
- `internal/service/invoice_service.go` - Business logic layer
- `internal/state_machine/invoice.go` - State machine validation

---

### Phase 3: Payment Gateway Integration (P0)

> **Duration**: 1 day  
> **Goal**: Secure callback handling

#### 3.1 Signature Validation
```go
// internal/service/payment_gateway_service.go
func ValidateSignature(provider string, payload []byte, signature string) bool {
    // Validate HMAC/MD5 signature from provider
    // Return true/false
}
```

#### 3.2 Duplicate Detection
```go
// Check (provider, provider_event_id) combination
// Return HTTP 200 with status "ignored" for duplicates
```

**Files to create:**
- `internal/service/payment_gateway_service.go`

---

### Phase 4: Event Processing (P1)

> **Duration**: 2 days  
> **Goal**: Inbound + Outbound events

#### 4.1 Outbox Poller
```go
// internal/infrastructure/outbox_poller.go
// - Poll unprocessed events from outbox_events table
// - Publish to RabbitMQ
// - Mark as published
// - Retry with exponential backoff
```

#### 4.2 Inbound Event Consumer
```go
// cmd/finance-service/events/
// - Handle pmb.applicant_created
// - Handle academic.student_created
// - Idempotent processing
```

**Files to create:**
- `cmd/finance-service/events/consumer.go`
- `internal/infrastructure/outbox_poller.go`

---

### Phase 5: Clearance Enhancement (P1)

> **Duration**: 1 day  
> **Goal**: Complete clearance logic

#### 5.1 Clearance State Machine
```go
// BLOCKED → CONDITIONAL → CLEARED → REVOKED
// - Validate transitions
// - Write audit log
// - Emit event
```

#### 5.2 Clearance Policy Evaluation
```go
// Evaluate policies based on service_scope
// Return default BLOCKED if no matching policy
```

**Files to update:**
- `internal/handler/finance_handler.go` - Add Clearances POST

---

### Phase 6: Additional Features (P2)

> **Duration**: 3 days  
> **Goal**: Complete all P2 features

#### 6.1 Scholarships CRUD
```go
// POST /api/v1/finance/scholarships
// GET /api/v1/finance/scholarships
// POST /api/v1/finance/scholarships/:id/approve
```

#### 6.2 Cash Accounts
```go
// GET /api/v1/finance/cash-accounts
// POST /api/v1/finance/cash-accounts/:id/transactions
```

#### 6.3 COA & Journals
```go
// GET /api/v1/finance/coa-accounts
// POST /api/v1/finance/journals
// Validate double-entry: SUM(debit) = SUM(credit)
```

#### 6.4 Budgets
```go
// GET /api/v1/finance/budgets
// POST /api/v1/finance/budgets
// Track realized_amount per line
```

#### 6.5 Payroll
```go
// POST /api/v1/finance/payroll-runs
// Approve payroll
// Calculate net = gross - deductions
```

#### 6.6 Disbursements
```go
// POST /api/v1/finance/disbursements
// Approve disbursement
```

---

### Phase 7: Infrastructure & Observability

> **Duration**: 1 day  
> **Goal**: Production ready

#### 7.1 Health Check
```go
// GET /health
// - Check DB connection
// - Check RabbitMQ connection
// - Return status
```

#### 7.2 Metrics
```go
// Export Prometheus metrics:
// - request_count_total
// - request_duration_seconds
// - outbox_pending_count
```

#### 7.3 Configuration
```go
// Read from env vars
// Validate on startup
```

**Files to update:**
- `cmd/finance-service/main.go` - Add health route, metrics

---

## 📦 File Structure (Target)

```
unsia-finance-service/
├── cmd/finance-service/
│   ├── main.go
│   └── events/
│       └── consumer.go          # NEW
├── internal/
│   ├── domain/
│   │   ├── models.go          # EXISTS ✓
│   │   └── validators.go     # NEW
│   ├── service/              # NEW (add service layer)
│   │   ├── invoice_service.go
│   │   ├── payment_service.go
│   │   ├── clearance_service.go
│   │   ├── scholarship_service.go
│   │   └── payment_gateway_service.go
│   ├── handler/
│   │   ├── finance_handler.go  # EXISTS ✓
│   │   └── health_handler.go # NEW
│   ├── middleware/
│   │   ├── auth_middleware.go   # NEW
│   │   └── rbac_middleware.go # NEW
│   ├── state_machine/        # NEW
│   │   ├── invoice.go
│   │   ├── payment.go
│   │   └── clearance.go
│   └── infrastructure/
│       ├── repository/
│       │   └── finance_repository.go  # EXISTS ✓
│       ├── database/
│       │   └── postgres.go     # EXISTS ✓
│       └── outbox_poller.go    # NEW
├── migrations/                     # EXISTS ✓ (26 files)
├── .env.example
├── Dockerfile
├── go.mod
└── README.md                    # EXISTS ✓
```

---

## ✅ Priority Checklist

### P0 - Critical (Must Have)
- [ ] JWT RS256 + JWKS auth middleware
- [ ] X-Application-Code, X-Active-Role, X-Correlation-Id validation
- [ ] RBAC permission middleware
- [ ] GET /api/v1/finance/invoices (list + filter + pagination)
- [ ] Invoice state machine validation
- [ ] Idempotency-Key support
- [ ] Payment callback signature validation
- [ ] Payment state machine
- [ ] Clearance state machine + policy evaluation
- [ ] Health check endpoint

### P1 - Important (Should Have)
- [ ] Outbox poller
- [ ] Inbound event consumer (pmb.applicant_created, academic.student_created)
- [ ] Scholarships CRUD
- [ ] Clearance POST endpoint

### P2 - Nice to Have
- [ ] Cash accounts
- [ ] COA + Journals with double-entry validation
- [ ] Budgets
- [ ] Payroll runs
- [ ] Disbursements

---

## 🔗 Dependencies

| Package | Usage | Status |
|---------|-------|--------|
| shared-auth | JWT validation | Required |
| shared-rbac | Permission check | Required |
| shared-errorenvelope | API response format | ✓ Using |
| shared-audit | Audit logging | ✓ Using |
| shared-idempotency | Idempotency | Required |
| shared-event | Event envelope | ✓ Using |
| shared-observability | Logging + metrics | Required |
| shared-httpclient | HTTP calls | Required |

---

## 📅 Timeline Estimate

| Phase | Duration | Cumulative |
|-------|----------|-----------|
| Phase 1: Security | 2 days | Day 1-2 |
| Phase 2: Invoice | 2 days | Day 3-4 |
| Phase 3: Payment Gateway | 1 day | Day 5 |
| Phase 4: Events | 2 days | Day 6-7 |
| Phase 5: Clearance | 1 day | Day 8 |
| Phase 6: P2 Features | 3 days | Day 9-11 |
| Phase 7: Infrastructure | 1 day | Day 12 |
| **Buffer/Testing** | 2 days | Day 13-14 |

**Total Estimate**: ~14 days (2 sprints)

---

*Plan created based on requirements.md analysis*
