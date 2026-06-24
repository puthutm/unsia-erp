# Brainstorming Plan: Pemisahan Finance Service Menjadi Microservices

## 📋 Latar Belakang

Berdasarkan feedback user, finance service yang saat ini monolith perlu dipisahkan menjadi beberapa microservice kecil untuk mengikuti arsitektur microservice yang benar. Setiap menu/domain akan menjadi service terpisah.

## 🎯 Tujuan

Memecah `unsia-finance-service` (monolith) menjadi beberapa microservice kecil yang independen:
- Setiap service memiliki database sendiri (schema terpisah di PostgreSQL)
- Komunikasi via RabbitMQ (event-driven)
- Setiap service bisa di-deploy secara independen
- Skalabilitas per service

---

## 📊 Proposed Microservices Architecture

### 1. **unsia-invoice-service** (_port: 8005_)
Fokus: Manajemen invoice mahasiswa

| Endpoint | Deskripsi |
|----------|----------|
| POST /api/v1/finance/invoices | Create invoice |
| GET /api/v1/finance/invoices | List invoices (filter, pagination) |
| GET /api/v1/finance/invoices/:id | Get detail invoice |
| PUT /api/v1/finance/invoices/:id | Update invoice |
| POST /api/v1/finance/invoices/:id/cancel | Cancel invoice |
| POST /api/v1/finance/invoices/:id/issue | Issue invoice |

**Database**: `invoice_db`
**Dependencies**: shared-auth, shared-rbac, shared-idempotency

---

### 2. **unsia-payment-service** (port: 8006_)
Fokus: Pembayaran dan callback

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/payments | List payments |
| POST /api/v1/finance/payment-callbacks/:provider | Webhook callback |
| POST /api/v1/finance/payment-verifications | Manual verification |
| GET /api/v1/finance/payments/:id | Payment detail |

**Database**: `payment_db`
**Dependencies**: shared-auth, shared-event, shared-httpclient

---

### 3. **unsia-clearance-service** (port: 8007_)
Fokus: Clearence/status mahasiswa

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/clearances | List clearances |
| GET /api/v1/finance/clearances/check | Check student clearance |
| POST /api/v1/finance/clearance-policies | Create policy |
| PUT /api/v1/finance/clearance-policies/:id | Update policy |
| GET /api/v1/finance/clearance-policies | List policies |

**Database**: `clearance_db`
**Dependencies**: shared-auth, shared-rbac

---

### 4. **unsia-scholarship-service** (port: 8008_)
Fokus: Manajemenbeasiswa

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/scholarships | List scholarships |
| POST /api/v1/finance/scholarships | Create scholarship |
| PUT /api/v1/finance/scholarships/:id | Update scholarship |
| DELETE /api/v1/finance/scholarships/:id | Delete scholarship |

**Database**: `scholarship_db`
**Dependencies**: shared-auth, shared-rbac

---

### 5. **unsia-budget-service** (port: 8009_)
Fokus: Rencana Anggaran Budget (RAB)

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/budgets | List budgets |
| POST /api/v1/finance/budgets | Create budget |
| GET /api/v1/finance/budgets/:id | Budget detail |
| PUT /api/v1/finance/budgets/:id | Update budget |
| GET /api/v1/finance/budgets/:id/utilization | Budget utilization |

**Database**: `budget_db`
**Dependencies**: shared-auth, shared-rbac

---

### 6. **unsia-cashbook-service** (port: 8010_)
Fokus: Kas dan Bank

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/cash-accounts | List cash accounts |
| POST /api/v1/finance/cash-accounts | Create cash account |
| GET /api/v1/finance/cash-accounts/:id | Cash account detail |
| GET /api/v1/finance/cash-accounts/:id/mutations | Mutations history |
| POST /api/v1/finance/cash-accounts/:id/mutations | Create mutation |

**Database**: `cashbook_db`
**Dependencies**: shared-auth, shared-rbac

---

### 7. **unsia-journal-service** (port: 8011_)
Fokus: Buku Besar dan Jurnal

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/journals | List journals |
| GET /api/v1/finance/journals/:id | Journal detail |
| POST /api/v1/finance/journals | Create journal entry |

**Database**: `journal_db`
**Dependencies**: shared-auth

---

### 8. **unsia-payroll-service** (port: 8012_)
Fokus: Payroll karyawan

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/payroll-runs | List payroll runs |
| POST /api/v1/finance/payroll-runs | Create payroll run |
| POST /api/v1/finance/payroll-runs/:id/approve | Approve payroll |

**Database**: `payroll_db`
**Dependencies**: shared-auth, shared-rbac

---

### 9. **unsia-disbursement-service** (port: 8013_)
Fokus: Pencairan dana (komisi CRM, dll)

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/disbursements | List disbursements |
| POST /api/v1/finance/disbursements | Create disbursement |
| POST /api/v1/finance/disbursements/:id/approve | Approve |
| POST /api/v1/finance/disbursements/:id/process | Process |

**Database**: `disbursement_db`
**Dependencies**: shared-auth, shared-rbac

---

### 10. **unsia-report-service** (port: 8014_)
Fokus: Laporan Keuangan

| Endpoint | Deskripsi |
|----------|----------|
| GET /api/v1/finance/reports/position | Neraca (Balance Sheet) |
| GET /api/v1/finance/reports/activity | Laba/Rugi (Income Statement) |
| GET /api/v1/finance/reports/cashflow | Arus Kas |

**Database**: Menggunakan view dari service lain (read-only replica)
**Dependencies**: shared-auth

---

## 🔗 Event Flow Diagram

```
┌──────────────────┐     Events      ┌───────────────────┐
│  unsia-invoice    │───────────────▶│  unsia-clearance  │
│    service       │  invoice.paid   │    service       │
└──────────────────┘                └───────────────────┘
         │                                     │
         │ Events                             │ Events
         ▼                                   ▼
┌──────────────────┐                ┌───────────────────┐
│  unsia-payment   │                │  unsia-journal   │
│    service     │───────────────▶│    service       │
└──────────────────┘  payment.done └───────────────────┘
         │
         │ Events
         ▼
┌──────────────────┐                ┌───────────────────┐
│  unsia-cashbook  │───────────────▶│  unsia-report     │
│    service     │  cash.updated   │    service       │
└──────────────────┘                └───────────────────┘
```

---

## 📁 Struktur Folder Target

```
unsia-docs-md/
├── services/
│   ├── unsia-invoice-service/      # NEW - PORT :8005
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── invoice_handler.go
│   │   │   ├── service/
│   │   │   │   └── invoice_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └��─ go.mod
│   │
│   ├── unsia-payment-service/       # NEW - PORT :8006
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── payment_handler.go
│   │   │   ├── service/
│   │   │   │   └── payment_gateway_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-clearance-service/      # NEW - PORT :8007
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── clearance_handler.go
│   │   │   ├── service/
│   │   │   │   └── clearance_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-scholarship-service/   # NEW - PORT :8008
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── scholarship_handler.go
│   │   │   ├── service/
│   │   │   │   └── scholarship_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-budget-service/       # NEW - PORT :8009
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── budget_handler.go
│   │   │   ├── service/
│   │   │   │   └── budget_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-cashbook-service/    # NEW - PORT :8010
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── cashbook_handler.go
│   │   │   ├── service/
│   │   │   │   └── cashbook_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-journal-service/     # NEW - PORT :8011
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── journal_handler.go
│   │   │   ├── service/
│   │   │   │   └── journal_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-payroll-service/     # NEW - PORT :8012
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── payroll_handler.go
│   │   │   ├── service/
│   │   │   │   └── payroll_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── unsia-disbursement-service/  # NEW - PORT :8013
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   │   └── models.go
│   │   │   ├── handler/
│   │   │   │   └── disbursement_handler.go
│   │   │   ├── service/
│   │   │   │   └── disbursement_service.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   └── unsia-report-service/      # NEW - PORT :8014
│       ├── cmd/
│       │   └── main.go
│       ├── internal/
│       │   ├── domain/
│       │   │   └── models.go
│       │   ├── handler/
│       │   │   └── report_handler.go
│       │   ├── service/
│       │   │   └── report_service.go
│       │   └── middleware/
│       │       └── auth_middleware.go
│       └── go.mod
```

---

## 🔄 Migration Steps

### Phase 1: Create New Services (Parallel)
- [ ] Create unsia-invoice-service structure
- [ ] Create unsia-payment-service structure  
- [ ] Create unsia-clearance-service structure

### Phase 2: Implement Core Services
- [ ] Implement unsia-invoice-service
- [ ] Implement unsia-payment-service
- [ ] Implement unsia-clearance-service

### Phase 3: Create Additional Services
- [ ] unsia-scholarship-service
- [ ] unsia-budget-service
- [ ] unsia-cashbook-service
- [ ] unsia-journal-service

### Phase 4: HR/Finance Services
- [ ] unsia-payroll-service
- [ ] unsia-disbursement-service

### Phase 5: Reporting
- [ ] unsia-report-service

---

## ✅ Checklist Persiapan

### Sebelum Memulai:
- [ ] Set up database schemas terpisah per service
- [ ] Configure RabbitMQ exchanges
- [ ] Prepare shared packages versions
- [ ] Update go.work untuk include semua services

### Dependency Graph:
```
unsia-invoice-service (independent)
    │
    ├── event: invoice.created ─────────▶ unsia-clearance-service
    │
    └── event: invoice.paid ──────────▶ unsia-payment-service
                                            │
                                            └── event: payment.completed ──▶ unsia-cashbook-service
                                                                                    │
                                                                                    └── event: cash.updated ──▶ unsia-journal-service
                                                                                                              │
                                                                                                              └── event: journal.created ──▶ unsia-report-service
```

---

## 📝 Catatan Tambahan

1. **Communication Pattern**: Gunakan event-driven via RabbitMQ untuk komunikasi antar service
2. **Shared Database**: Beberapa service boleh share database jika saling terkait erat (misal: invoice + payment)
3. **API Gateway**: pertimbangkan menggunakan API gateway (nginx/traefik) di depan
4. **Service Discovery**: Gunakan service discovery untuk production environment

---

## ❓ Pertanyaan untuk User

1. Apakah pembagian ini sudah sesuai dengan kebutuhan?
2. Ada domain spesifik yang ingin digabungkan atau dipisahkan lagi?
3. Prioritas implementasi mana yang lebih duluan?

---

*Plan ini dibuat berdasarkan brainstorming untuk memecah finance service monolith menjadi microservices.*
