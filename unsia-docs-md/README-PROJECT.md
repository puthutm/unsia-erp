# UNSIA ERP — Project Structure

**Stack:** Go 1.22+ (Gin) · Next.js 14+ (App Router) · PostgreSQL · RabbitMQ · Redis

## Struktur Folder

```
unsia-docs-md/
│
├── services/                        → Backend Go services (1 folder = 1 modul)
│   ├── unsia-core-service/          ← SSO, Auth, RBAC, OAuth 2.0 server
│   ├── unsia-reference-service/     ← Master data authority
│   ├── unsia-crm-service/           ← Lead & campaign management
│   ├── unsia-pmb-service/           ← Applicant lifecycle
│   ├── unsia-finance-service/       ← Invoice, payment, clearance
│   ├── unsia-academic-service/      ← Mahasiswa, KRS, nilai, transkrip
│   ├── unsia-hris-service/          ← Dosen & karyawan
│   ├── unsia-lms-service/           ← Pembelajaran online
│   ├── unsia-assessment-service/    ← CBT, quiz, scoring engine
│   ├── unsia-portal-service/        ← Notification & dashboard backend
│   └── unsia-integration-worker/    ← Outbox/inbox, RabbitMQ, DLQ, reconciliation
│
├── packages/                        → Shared Go modules
│   ├── shared-auth/                 ← JWT validation, JWKS cache
│   ├── shared-rbac/                 ← Permission check, data scope
│   ├── shared-audit/                ← Audit log writer
│   ├── shared-idempotency/          ← Duplicate command prevention
│   ├── shared-event/                ← Outbox writer, inbox consumer, event envelope
│   ├── shared-httpclient/           ← Service-to-service HTTP client
│   ├── shared-errorenvelope/        ← Standard error/success response format
│   └── shared-observability/        ← Logging, tracing, metrics
│
├── frontend/
│   └── unsia-portal-web/            ← Next.js App Router (semua role)
│
├── infra/                           ← Docker, Nginx, PostgreSQL, RabbitMQ, monitoring
│
├── 01-prd/                          ← Product Requirements Document
├── 02-brd/                          ← Business Requirements Document
├── 03-fsd/                          ← Functional Specification Document
├── 04-api-contract/                 ← OpenAPI / Swagger
├── 05-event-contract/               ← Event Contract
├── 06-erd-dbml/                     ← Database schema (DBML)
├── 07-uat/                          ← UAT Scenario & QA Test Plan
├── 08-developer/                    ← Developer Implementation Specification
├── 09-workplan/                     ← Sprint Plan & Delivery Roadmap
├── 10-repo-structure/               ← Repo structure guide
└── 11-srs/                          ← Software Requirements Specification
```

## Prinsip Utama

| Prinsip | Aturan |
|---------|--------|
| 1 modul = 1 database | Tidak ada cross-database FK atau join untuk transaksi online |
| Identity terpusat | Semua modul validasi JWT dari Core, tidak ada login sendiri |
| Event-driven | Perubahan penting ditulis ke outbox dan dikonsumsi via inbox |
| Idempotent | Semua command kritis aman di-retry tanpa duplikasi data |
| Degraded mode | Setiap modul punya fallback saat dependency down |
| Audit trail | Semua aksi sensitif tercatat dengan actor, role, dan timestamp |

## Phase Delivery

| Phase | Fokus | Sprint |
|-------|-------|--------|
| 0 | Architecture Foundation | Sprint 0 |
| 1 | Core Service + SSO Internal | Sprint 1 |
| 2 | SSO External App Registration | Sprint 2–3 |
| 3 | Migrate All Modules ke Go | Sprint 4–8 |
| 4 | Integration + Hardening + UAT | Sprint 9–10 |

## Dokumentasi Teknis

Lihat masing-masing `README.md` di setiap folder service/package untuk detail tanggung jawab, endpoint, dependencies, dan aturan penting.
