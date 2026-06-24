# Implementation Plan — UNSIA ERP

## Tujuan

Membangun seluruh UNSIA ERP dari nol dengan pendekatan **database-first**:
1. **Database** — SQL migration per modul, semua relasi & constraint benar
2. **Backend** — Go/Gin API per service di folder `services/`
3. **Frontend** — Next.js App Router di folder `frontend/`

---

## Folder Structure Final

```
d:\Superman\Superman\Coding\New folder\Dockument\candi\unsia-docs-md\
├── services/                          ← Backend Go (1 folder = 1 modul)
│   ├── unsia-core-service/
│   │   ├── cmd/core-service/main.go
│   │   ├── internal/
│   │   ├── migrations/               ← SQL migration files
│   │   ├── tests/
│   │   ├── go.mod
│   │   └── Dockerfile
│   ├── unsia-reference-service/
│   ├── unsia-crm-service/
│   ├── unsia-pmb-service/
│   ├── unsia-finance-service/
│   ├── unsia-academic-service/
│   ├── unsia-hris-service/
│   ├── unsia-lms-service/
│   ├── unsia-assessment-service/
│   └── unsia-integration-worker/
│
├── packages/                          ← Shared Go modules
│   ├── shared-auth/
│   ├── shared-rbac/
│   ├── shared-audit/
│   ├── shared-idempotency/
│   ├── shared-event/
│   ├── shared-httpclient/
│   ├── shared-errorenvelope/
│   └── shared-observability/
│
├── frontend/                          ← Frontend Next.js
│   └── unsia-portal-web/
│       ├── app/
│       ├── src/
│       ├── package.json
│       └── next.config.ts
│
├── infra/                             ← Docker, Nginx, DB init, monitoring
│
└── (existing docs: 01-prd/, 06-erd-dbml/, etc.)
```

---

## Open Questions

> [!IMPORTANT]
> **Q1:** Apakah migration SQL digenerate langsung sebagai file `.sql` siap pakai untuk `golang-migrate`? Atau mau pakai format lain? (Status: **Selesai - Menggunakan file `.sql` untuk `golang-migrate`**)

> [!IMPORTANT]
> **Q2:** Untuk relasi lintas database (misal `finance.invoices.applicant_id → pmb.applicants.id`), DBML menuliskan sebagai `Ref` tapi kita **tidak boleh buat FK cross-database**. Rencananya field tersebut akan tetap ada sebagai `uuid` biasa tanpa FK constraint, dan integritas dijaga via API/event. Setuju? (Status: **Selesai - Menggunakan UUID biasa dengan anotasi comment**)

> [!IMPORTANT]
> **Q3:** Apakah tabel OAuth baru dari Kiro spec (B-8: `oauth_authorization_codes`, `oauth_access_tokens`, `oauth_refresh_tokens`, `client_registration_requests`) juga langsung dimasukkan ke migration `core_db`? (Status: **Selesai - Dimasukkan ke migration `core_db`**)

> [!IMPORTANT]
> **Q4:** Folder `services/` dan `frontend/` langsung di bawah `unsia-docs-md/` (workspace saat ini), atau mau di folder baru terpisah? (Status: **Selesai - Berada di bawah workspace `unsia-docs-md/`**)

> [!IMPORTANT]
> **Q5:** Bagaimana cara setup private key dan public key untuk tanda tangan JWT RS256 di `core-service`?
> *Rencana:* Jika file `.keys/private.pem` tidak ditemukan, `core-service` akan otomatis meng-generate key pair baru berukuran 2048-bit saat start dan menyimpannya ke folder tersebut. Apakah Anda setuju?

> [!IMPORTANT]
> **Q6:** Terkait penyimpanan *credential* user di tabel `users`.
> *Rencana:* Password akan di-hash menggunakan standard library `golang.org/x/crypto/bcrypt`. Apakah Anda setuju?

> [!IMPORTANT]
> **Q7:** Bagaimana seeding data awal User, Person, dan Role agar sistem bisa langsung ditest login setelah running?
> *Rencana:* Kita akan buat file SQL seed `seed-core-data.sql` di folder `services/unsia-core-service/migrations/` (atau sebagai endpoint seed/CLI seed) yang berisi minimal data user super-admin dan prodi-admin untuk testing. Apakah Anda setuju? (Status: **Selesai**)

> [!IMPORTANT]
> **Q8:** Bagaimana seeding untuk data awal Reference?
> *Rencana:* Kita akan buat file SQL seed `seed-reference-data.sql` di folder `services/unsia-reference-service/migrations/` yang berisi minimal data wilayah Indonesia, agama, status_codes standar, dan beberapa program studi awal. Apakah Anda setuju?

---

## Phase 1: Database Migration (SQL) — Database-First

Membuat file SQL migration untuk `golang-migrate` di setiap service. Urutan berdasarkan dependency:

### Aturan Migration

- Format file: `{version}_{description}.up.sql` / `{version}_{description}.down.sql`
- Setiap database punya schema sendiri (atau pakai `public` schema)
- FK **hanya** di dalam database yang sama
- Field referensi lintas database → `uuid` biasa tanpa FK, pakai comment `-- external_ref: {module}.{table}.id`
- Setiap database wajib punya 4 tabel teknis: `audit_logs`, `idempotency_keys`, `outbox_events`, `inbox_events`, `reconciliation_mismatch_logs`
- UUID sebagai PK, `gen_random_uuid()` sebagai default

---

### 1.1 `core_db` — 20 tabel domain + 5 tabel event infra + 5 tabel OAuth

| Migration | Tabel | Catatan |
|-----------|-------|---------|
| `000001_create_persons` | `persons` | Identitas dasar |
| `000002_create_users` | `users` | FK → persons |
| `000003_create_roles` | `roles` | scope_type: global/prodi/module/self |
| `000004_create_permissions` | `permissions` | Pola `module.resource.action` |
| `000005_create_user_roles` | `user_roles` | FK → users, roles; UQ(user,role,study_program_id) |
| `000006_create_role_permissions` | `role_permissions` | FK → roles, permissions |
| `000007_create_applications` | `applications` | Registry aplikasi |
| `000008_create_oauth_clients` | `oauth_clients` | FK → applications |
| `000009_create_redirect_uris` | `redirect_uris` | FK → oauth_clients |
| `000010_create_service_tokens` | `service_tokens` | FK → applications |
| `000011_create_sessions` | `sessions` | FK → users |
| `000012_create_active_role_sessions` | `active_role_sessions` | FK → users, roles, sessions, applications |
| `000013_create_impersonation_sessions` | `impersonation_sessions` | FK → users (actor/target), roles, applications |
| `000014_create_audit_logs` | `audit_logs` | FK → users, roles, impersonation_sessions, applications |
| `000015_create_idempotency_keys` | `idempotency_keys` | UQ(module, idempotency_key) |
| `000016_create_integration_event_logs` | `integration_event_logs` | UQ(source, target, type, key) |
| `000017_create_event_contracts` | `event_contracts` | UQ(event_name, event_version) |
| `000018_create_event_consumers` | `event_consumers` | FK → event_contracts |
| `000019_create_event_replay_logs` | `event_replay_logs` | DLQ replay audit |
| `000020_create_outbox_events` | `outbox_events` | Standard outbox |
| `000021_create_inbox_events` | `inbox_events` | Standard inbox |
| `000022_create_reconciliation_mismatch_logs` | `reconciliation_mismatch_logs` | Reconciliation |
| `000023_create_oauth_authorization_codes` | `oauth_authorization_codes` | PKCE, single-use code *(dari Kiro spec B-8)* |
| `000024_create_oauth_access_tokens` | `oauth_access_tokens` | JTI-based *(dari Kiro spec B-8)* |
| `000025_create_oauth_refresh_tokens` | `oauth_refresh_tokens` | Rotation *(dari Kiro spec B-8)* |
| `000026_create_client_registration_requests` | `client_registration_requests` | Audit trail *(dari Kiro spec B-8)* |

**Total: ~26 migration files untuk core_db**

---

### 1.2 `reference_db` — 16 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001_create_countries` | `countries` |
| `000002_create_provinces` | `provinces` → FK countries |
| `000003_create_cities` | `cities` → FK provinces |
| `000004_create_districts` | `districts` → FK cities |
| `000005_create_villages` | `villages` → FK districts |
| `000006_create_religions` | `religions` |
| `000007_create_study_programs` | `study_programs` |
| `000008_create_academic_years` | `academic_years` |
| `000009_create_academic_periods` | `academic_periods` → FK academic_years |
| `000010_create_admission_paths` | `admission_paths` |
| `000011_create_pmb_waves` | `pmb_waves` → FK academic_years, academic_periods, admission_paths |
| `000012_create_lead_sources` | `lead_sources` |
| `000013_create_document_types` | `document_types` |
| `000014_create_payment_components` | `payment_components` |
| `000015_create_payment_methods` | `payment_methods` |
| `000016_create_employee_types` | `employee_types` |
| `000017_create_lecturer_statuses` | `lecturer_statuses` |
| `000018_create_status_codes` | `status_codes` → UQ(module, code) |
| `000019_create_outbox_events` | event infra |
| `000020_create_inbox_events` | event infra |
| `000021_create_idempotency_keys` | event infra |
| `000022_create_reconciliation_mismatch_logs` | event infra |

**Total: ~22 migration files**

---

### 1.3 `crm_db` — 8 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001_create_campaigns` | `campaigns` |
| `000002_create_agents` | `agents` |
| `000003_create_referrals` | `referrals` → FK agents |
| `000004_create_leads` | `leads` → FK campaigns, referrals |
| `000005_create_lead_activities` | `lead_activities` → FK leads |
| `000006_create_lead_status_histories` | `lead_status_histories` → FK leads |
| `000007_create_commission_rules` | `commission_rules` |
| `000008_create_commission_records` | `commission_records` → FK leads, commission_rules |
| `000009–000012` | 4 tabel event infra |

**Total: ~12 migration files**

---

### 1.4 `pmb_db` — 12 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001_create_applicants` | `applicants` (person_id, study_program_id sebagai external ref) |
| `000002_create_applicant_biodata` | → FK applicants |
| `000003_create_applicant_addresses` | → FK applicants |
| `000004_create_applicant_education_backgrounds` | → FK applicants |
| `000005_create_applicant_family_members` | → FK applicants |
| `000006_create_applicant_financial_profiles` | → FK applicants |
| `000007_create_applicant_facility_profiles` | → FK applicants |
| `000008_create_applicant_documents` | → FK applicants |
| `000009_create_applicant_status_histories` | → FK applicants |
| `000010_create_re_registrations` | → FK applicants |
| `000011_create_loa_documents` | → FK applicants |
| `000012_create_handover_logs` | → FK applicants, UQ idempotency_key |
| `000013–000016` | 4 tabel event infra |

**Total: ~16 migration files**

---

### 1.5 `finance_db` — 21 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001_create_invoices` | `invoices` (applicant_id, student_id = external ref) |
| `000002_create_invoice_items` | → FK invoices |
| `000003_create_payments` | → FK invoices |
| `000004_create_payment_gateway_callbacks` | → FK payments |
| `000005_create_payment_verifications` | → FK payments |
| `000006_create_scholarships` | external ref: student_id |
| `000007_create_installment_requests` | → FK invoices |
| `000008_create_clearance_policies` | clearance rules |
| `000009_create_student_clearances` | source of truth clearance |
| `000010_create_clearance_dispensations` | → FK student_clearances |
| `000011_create_cash_accounts` | kas/bank |
| `000012_create_cash_transactions` | → FK cash_accounts |
| `000013_create_payroll_runs` | payroll |
| `000014_create_payroll_items` | → FK payroll_runs |
| `000015_create_disbursements` | pencairan komisi |
| `000016_create_tax_records` | pajak |
| `000017_create_bpjs_records` | BPJS |
| `000018_create_coa_accounts` | chart of accounts |
| `000019_create_journals` | jurnal keuangan |
| `000020_create_journal_entries` | → FK journals, coa_accounts |
| `000021_create_budgets` | anggaran |
| `000022_create_budget_lines` | → FK budgets, coa_accounts |
| `000023–000026` | 4 tabel event infra |

**Total: ~26 migration files**

---

### 1.6 `academic_db` — 25 tabel domain + 4 tabel event infra

| Migration | Tabel utama |
|-----------|-------------|
| `000001` | `students` (applicant_id, person_id = external ref) |
| `000002` | `student_advisors` → FK students |
| `000003` | `nim_format_configs` |
| `000004` | `nim_sequences` (row-level lock untuk NIM generation) |
| `000005` | `academic_period_study_program_settings` |
| `000006` | `academic_settings` |
| `000007` | `curriculums` |
| `000008` | `courses` |
| `000009` | `curriculum_courses` → FK curriculums, courses |
| `000010` | `class_packages` |
| `000011` | `class_package_items` |
| `000012` | `course_offerings` → FK courses |
| `000013` | `classes` → FK course_offerings |
| `000014` | `class_lecturers` → FK classes |
| `000015` | `class_schedules` → FK classes |
| `000016` | `krs` → FK students |
| `000017` | `krs_items` → FK krs, classes |
| `000018` | `grades` → FK krs_items |
| `000019` | `grade_histories` → FK grades |
| `000020` | `khs` → FK students |
| `000021` | `transcripts` → FK students |
| `000022` | `academic_letters` → FK students |
| `000023` | `graduation_requirements` |
| `000024` | `yudisium_records` → FK students |
| `000025` | `alumni` → FK students |
| `000026–000029` | 4 tabel event infra |

**Total: ~29 migration files**

---

### 1.7 `hris_db` — 11 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001` | `work_units` (self-referencing parent_unit_id) |
| `000002` | `positions` |
| `000003` | `functional_positions` |
| `000004` | `employees` → FK work_units, positions |
| `000005` | `lecturers` → FK employees, functional_positions |
| `000006` | `attendances` → FK employees |
| `000007` | `leave_requests` → FK employees |
| `000008` | `bkd_records` → FK lecturers |
| `000009` | `performance_reviews` → FK employees |
| `000010` | `certifications` → FK employees |
| `000011` | `payroll_sources` → FK employees |
| `000012–000015` | 4 tabel event infra |

**Total: ~15 migration files**

---

### 1.8 `lms_db` — 14 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001` | `classes` (academic_class_id = external ref) |
| `000002` | `enrollments` → FK classes |
| `000003` | `sessions` → FK classes |
| `000004` | `materials` → FK sessions |
| `000005` | `videos` → FK sessions |
| `000006` | `vicon_links` → FK sessions |
| `000007` | `assignments` → FK sessions |
| `000008` | `assignment_submissions` → FK assignments |
| `000009` | `quiz_activities` → FK sessions |
| `000010` | `discussions` → FK sessions |
| `000011` | `discussion_comments` → FK discussions (+ self-ref) |
| `000012` | `attendances` → FK sessions |
| `000013` | `learning_progress` → FK enrollments |
| `000014` | `grade_syncs` → FK classes |
| `000015–000018` | 4 tabel event infra |

**Total: ~18 migration files**

---

### 1.9 `assessment_db` — 14 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001` | `question_banks` |
| `000002` | `questions` → FK question_banks |
| `000003` | `question_versions` → FK questions |
| `000004` | `question_options` → FK questions |
| `000005` | `material_banks` |
| `000006` | `materials` → FK material_banks |
| `000007` | `question_sets` |
| `000008` | `question_set_items` → FK question_sets, questions |
| `000009` | `assessment_sessions` → FK question_sets |
| `000010` | `assessment_participants` → FK assessment_sessions |
| `000011` | `assessment_attempts` → FK assessment_sessions, participants |
| `000012` | `assessment_answers` → FK attempts, questions, options |
| `000013` | `assessment_scores` → FK attempts |
| `000014` | `surveys`, `survey_questions`, `survey_responses` |
| `000015–000018` | 4 tabel event infra |

**Total: ~18 migration files**

---

### 1.10 `portal_db` — 5 tabel domain + 4 tabel event infra

| Migration | Tabel |
|-----------|-------|
| `000001` | `notifications` |
| `000002` | `notification_reads` → FK notifications |
| `000003` | `user_preferences` |
| `000004` | `menu_shortcuts` |
| `000005` | `portal_activity_logs` |
| `000006–000009` | 4 tabel event infra |

**Total: ~9 migration files**

---

### Ringkasan Database

| Database | Domain Tables | Event Infra | Total Migrations |
|----------|:---:|:---:|:---:|
| `core_db` | 22 + 5 OAuth | 5 + 3 catalog | ~26 |
| `reference_db` | 18 | 4 | ~22 |
| `crm_db` | 8 | 4 | ~12 |
| `pmb_db` | 12 | 4 | ~16 |
| `finance_db` | 22 | 4 | ~26 |
| `academic_db` | 25 | 4 | ~29 |
| `hris_db` | 11 | 4 | ~15 |
| `lms_db` | 14 | 4 | ~18 |
| `assessment_db` | 14 | 4 | ~18 |
| `portal_db` | 5 | 4 | ~9 |
| **TOTAL** | **~160** | **~40** | **~191** |

---

## Phase 2: Infrastructure Setup

### Folder `infra/`

#### [NEW] `infra/docker/docker-compose.local.yml`
- 10 PostgreSQL instances (atau 1 instance + 10 databases)
- Redis
- RabbitMQ
- Nginx (API Gateway)

#### [NEW] `infra/postgres/init-databases.sql`
- Script `CREATE DATABASE` untuk semua 10 database + `CREATE EXTENSION "uuid-ossp"`

#### [NEW] `infra/rabbitmq/definitions.json`
- Exchange, queue, binding untuk semua event

---

## Phase 3: Shared Go Modules

### Folder `packages/`

| Package | File Utama | Fungsi |
|---------|-----------|--------|
| [NEW] `shared-errorenvelope/` | `envelope.go` | Success/error response format |
| [NEW] `shared-auth/` | `jwt.go`, `jwks.go`, `service_token.go` | JWT RS256 validation, JWKS cache |
| [NEW] `shared-rbac/` | `permission.go`, `scope.go` | Permission check, data scope resolver |
| [NEW] `shared-audit/` | `audit.go` | Audit log writer |
| [NEW] `shared-idempotency/` | `idempotency.go` | Request hash, response cache, lock |
| [NEW] `shared-event/` | `outbox.go`, `inbox.go`, `envelope.go` | Outbox/inbox pattern |
| [NEW] `shared-httpclient/` | `client.go` | Service-to-service HTTP, circuit breaker |
| [NEW] `shared-observability/` | `logger.go`, `tracing.go`, `metrics.go` | Structured logging, trace propagation |

---

## Phase 4: Backend Go Services

Setiap service mengikuti pola clean architecture:

```
services/unsia-{module}-service/
├── cmd/{module}-service/main.go
├── internal/
│   ├── domain/         ← entities, value objects, enums, state machines
│   ├── application/    ← commands, queries, services (use cases)
│   ├── infrastructure/ ← repositories, external clients, event publisher
│   ├── handler/        ← HTTP controllers, validators, presenters
│   └── middleware/     ← auth, RBAC, correlation, idempotency
├── migrations/         ← symlink atau copy dari Phase 1
├── tests/
│   ├── unit/
│   ├── integration/
│   └── contract/
├── go.mod
├── Dockerfile
└── .env.example
```

### Urutan implementasi (berdasarkan dependency):

| Urutan | Service | Endpoint Prioritas |
|--------|---------|-------------------|
| 1 | **core-service** | `/auth/login`, `/auth/refresh`, `/auth/me`, `/auth/switch-role`, JWKS, OAuth endpoints |
| 2 | **reference-service** | `/ref/study-programs`, `/ref/academic-years`, `/ref/academic-periods`, semua master data |
| 3 | **crm-service** | `/crm/leads`, `/crm/campaigns`, `/crm/leads/{id}/convert-to-applicant` |
| 4 | **pmb-service** | `/pmb/applicants`, submit, documents, request-invoice, issue-loa, handover |
| 5 | **finance-service** | `/finance/invoices`, payment-callbacks, clearances |
| 6 | **academic-service** | `/academic/students/generate-from-applicant`, KRS, grades, KHS, transcripts |
| 7 | **hris-service** | `/hris/lecturers`, `/hris/employees` |
| 8 | **lms-service** | `/lms/classes/sync-from-academic`, enrollments, sessions, grade-syncs |
| 9 | **assessment-service** | `/assessment/sessions`, attempts, results |
| 10 | **integration-worker** | Outbox publisher, inbox consumers, DLQ replay, reconciliation |

---

### [Component Name] unsia-core-service

Merupakan modul pusat identitas, SSO, otorisasi, dan manajemen sesi.

#### [NEW] [go.mod](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-core-service/go.mod)
Menginisialisasi modul Go `github.com/unsia-erp/unsia-core-service` dan me-replace shared modules ke local path.

#### [NEW] [.env.example](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-core-service/.env.example)
Variabel lingkungan konfigurasi database (`core_db`), port server (default `:8001`), dan rahasia enkripsi JWT/session.

#### [NEW] [main.go](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-core-service/cmd/core-service/main.go)
Inisialisasi aplikasi, database migration runner (opsional), setup routing Gin, dan inisialisasi kunci RSA.

#### [NEW] [domain](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-core-service/internal/domain)
Struktur entitas domain:
- `Person` & `User` (password di-hash dengan `bcrypt`).
- `Role` & `Permission`.
- `ActiveRoleSession` & `ImpersonationSession`.
- `ApplicationRegistry` & `ServiceToken`.

#### [NEW] [infrastructure](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-core-service/internal/infrastructure)
- Database connection pool & GORM/sqlc queries.
- Repositori untuk users, roles, sessions, dan outbox/inbox events.
- In-memory/File-based RSA key store untuk JWKS.

#### [NEW] [handler](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-core-service/internal/handler)
Fungsi handler API sesuai OpenAPI:
- `/api/v1/auth/login` (generate JWT access token + refresh token).
- `/api/v1/auth/refresh` (validate refresh token, issue new tokens).
- `/api/v1/auth/me` (mengembalikan claims user beserta roles dan scopes).
- `/api/v1/auth/switch-role` (mengubah active role saat ini).
- `/api/v1/applications` (list modul aplikasi yang terdaftar).
- `/.well-known/jwks.json` (endpoint publik berisi public key JWK untuk verifikasi token oleh microservice downstream).

---

### [Component Name] unsia-reference-service

Merupakan modul pusat data master untuk program studi, tahun ajaran, periode akademik, metode pembayaran, dsb.

#### [NEW] [go.mod](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-reference-service/go.mod)
Menginisialisasi modul Go `github.com/unsia-erp/unsia-reference-service` dan me-replace shared modules ke local path.

#### [NEW] [.env.example](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-reference-service/.env.example)
Variabel lingkungan konfigurasi database (`reference_db`), port server (default `:8002`), dan URL JWKS Core Service untuk verifikasi token.

#### [NEW] [main.go](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-reference-service/cmd/reference-service/main.go)
Inisialisasi aplikasi, setup routing Gin, dan inisialisasi JWKS cache client.

#### [NEW] [domain](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-reference-service/internal/domain)
Struktur entitas domain:
- `StudyProgram` (prodi).
- `AcademicYear` (tahun ajaran).
- `AcademicPeriod` (periode akademik).
- `PaymentComponent` & `PaymentMethod`.
- `DocumentType`.
- `StatusCode` (managed enum).
- `Country`, `Province`, `City`, `District`, `Village` & `Religion`.

#### [NEW] [infrastructure](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-reference-service/internal/infrastructure)
- Database connection pool & GORM queries.
- Repositori untuk master data, outbox/inbox events.

#### [NEW] [handler](file:///d:/Superman/Superman/Coding/New%20folder/Dockument/candi/unsia-docs-md/services/unsia-reference-service/internal/handler)
Fungsi handler API sesuai OpenAPI:
- `/api/v1/ref/study-programs` (GET list, POST create, publish event `reference.study_program_updated`).
- `/api/v1/ref/academic-years` (GET list, POST create).
- `/api/v1/ref/academic-periods` (GET list, POST create, publish event `reference.academic_period_updated`).
- `/api/v1/ref/status-codes` (GET list).
- `/api/v1/ref/payment-components` (GET list).
- `/api/v1/ref/document-types` (GET list).

---

## Phase 5: Frontend Next.js (App Router)

### Folder `frontend/unsia-portal-web/`

```
frontend/unsia-portal-web/
├── app/
│   ├── (auth)/login/page.tsx
│   ├── (auth)/select-role/page.tsx
│   ├── (portal)/dashboard/page.tsx
│   ├── (portal)/notifications/page.tsx
│   ├── pendaftar/...
│   ├── mahasiswa/...
│   ├── dosen/...
│   ├── admin/...
│   └── pimpinan/...
├── src/
│   ├── components/
│   ├── features/
│   ├── lib/api-client.ts
│   ├── hooks/
│   ├── stores/
│   └── types/
├── package.json
├── next.config.ts
└── tsconfig.json
```

### Tech stack frontend:
- Next.js 14+ App Router
- TanStack Query untuk data fetching
- Zustand untuk client state
- Zod untuk form validation

---

## Phase 6: Testing & Hardening

### Automated Tests
- `go test ./...` per service
- Integration tests: CRM→PMB, PMB→Finance, PMB→Academic, Academic→LMS
- Event contract tests: outbox created, inbox idempotent, DLQ retry
- E2E tests: Playwright untuk frontend

### Manual Verification
- Docker compose up seluruh stack
- Test critical flows end-to-end
- Verify cross-module event propagation

---

## Verification Plan

### Database
```bash
# Per service
cd services/unsia-core-service
migrate -database "postgres://..." -path migrations up
# Verify schema
psql -d core_db -c "\dt"
```

### Backend
```bash
cd services/unsia-core-service
go test ./tests/unit/...
go test ./tests/integration/...
```

### Frontend
```bash
cd frontend/unsia-portal-web
npm run build   # build check
npm run test    # unit tests
```

---

## Estimasi Volume

| Item | Jumlah |
|------|--------|
| Migration files (up+down) | ~382 files |
| Go service repos | 10 |
| Shared Go modules | 8 |
| Frontend pages (est.) | 30+ |
| Total database tables | ~200 |
