---
title: "Struktur Repo ERP UNSIA"
source_file: "Struktur_Repo_ERP_UNSIA.docx"
format: markdown
---

# Struktur Repo ERP UNSIA

STRUKTUR REPO ERP UNSIA

Model Multi-Repo per Modul / Service

Untuk implementasi ERP Pendidikan / SIAKAD Terintegrasi UNSIA

| Item | Keterangan |
| --- | --- |
| Dokumen | Struktur Repo ERP UNSIA |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Pendekatan | Multi-repo per modul/service |
| Stack Rekomendasi | Next.js / Node.js + TypeScript, PostgreSQL per modul, Redis, RabbitMQ/BullMQ |
| Tujuan | Menjadi acuan developer dalam membuat repository, folder, boundary, dan standar kerja antar modul. |

Ringkasan: struktur repo ini memakai pendekatan satu repo untuk satu modul/service. Setiap service memiliki database utama masing-masing dan tidak boleh melakukan direct database join atau direct write ke database modul lain. Integrasi lintas modul dilakukan melalui API contract, event contract, outbox/inbox, snapshot/read model, dan shared contracts.

# 1. Struktur Repo Utama

unsia-core-service
unsia-reference-service
unsia-crm-service
unsia-pmb-service
unsia-finance-service
unsia-academic-service
unsia-hris-service
unsia-lms-service
unsia-assessment-service
unsia-portal-web
unsia-integration-worker
unsia-shared-contracts
unsia-infra
unsia-docs

# 2. Fungsi Masing-Masing Repo

| Repo | Fungsi | Database |
| --- | --- | --- |
| unsia-core-service | Login, SSO, user, role, permission, active role, service token, audit | core_db |
| unsia-reference-service | Master data, prodi, tahun ajaran, periode akademik, status code | reference_db |
| unsia-crm-service | Lead, campaign, agent, referral, follow-up | crm_db |
| unsia-pmb-service | Applicant, biodata, dokumen, seleksi, LoA, handover | pmb_db |
| unsia-finance-service | Invoice, payment, callback, receipt, clearance | finance_db |
| unsia-academic-service | Mahasiswa, NIM, kurikulum, kelas, KRS, nilai, KHS, transkrip | academic_db |
| unsia-hris-service | Dosen, pegawai, homebase, jabatan, status aktif | hris_db |
| unsia-lms-service | Kelas online, enrollment, materi, tugas, presensi, grade input | lms_db |
| unsia-assessment-service | Bank soal, CBT, quiz, attempt, scoring, result | assessment_db |
| unsia-portal-web | Dashboard, portal mahasiswa, dosen, admin, pimpinan, notification UI | portal_db |
| unsia-integration-worker | Outbox, inbox, retry, DLQ, reconciliation, event consumer | Tidak wajib punya DB sendiri |
| unsia-shared-contracts | OpenAPI, event contract, schema, RBAC, error code | - |
| unsia-infra | Docker, Nginx, Redis, RabbitMQ, deployment, monitoring | - |
| unsia-docs | PRD, BRD, FSD, UAT, ERD, DBML, release notes | - |

# 3. Struktur Standar Setiap Repo Service

Contoh struktur berikut memakai repo unsia-pmb-service sebagai acuan. Pola yang sama dapat dipakai untuk Core, Finance, Academic, LMS, Assessment, dan modul lainnya.

unsia-pmb-service/
│
├── app/
│ ├── api/
│ │ └── v1/
│ │ └── pmb/
│ │ ├── applicants/
│ │ │ ├── route.ts
│ │ │ └── [applicantId]/
│ │ │ ├── route.ts
│ │ │ ├── submit/
│ │ │ │ └── route.ts
│ │ │ ├── documents/
│ │ │ │ └── route.ts
│ │ │ ├── request-invoice/
│ │ │ │ └── route.ts
│ │ │ ├── issue-loa/
│ │ │ │ └── route.ts
│ │ │ └── handover-to-academic/
│ │ │ └── route.ts
│ │ └── health/
│ │ └── route.ts
│ │
│ └── pmb/
│ ├── applicants/
│ ├── documents/
│ ├── loa/
│ └── dashboard/
│
├── src/
│ ├── domain/
│ │ ├── entities/
│ │ │ ├── applicant.entity.ts
│ │ │ ├── applicant-document.entity.ts
│ │ │ └── loa.entity.ts
│ │ ├── value-objects/
│ │ ├── enums/
│ │ └── state-machines/
│ │ └── applicant-state.machine.ts
│ │
│ ├── application/
│ │ ├── commands/
│ │ │ ├── create-applicant.command.ts
│ │ │ ├── submit-applicant.command.ts
│ │ │ ├── verify-document.command.ts
│ │ │ ├── request-invoice.command.ts
│ │ │ ├── issue-loa.command.ts
│ │ │ └── handover-to-academic.command.ts
│ │ ├── queries/
│ │ │ ├── get-applicant.query.ts
│ │ │ └── list-applicants.query.ts
│ │ └── services/
│ │ ├── applicant.service.ts
│ │ ├── document.service.ts
│ │ └── loa.service.ts
│ │
│ ├── infrastructure/
│ │ ├── database/
│ │ │ ├── prisma.ts
│ │ │ └── transaction.ts
│ │ ├── repositories/
│ │ │ ├── applicant.repository.ts
│ │ │ ├── document.repository.ts
│ │ │ └── loa.repository.ts
│ │ ├── external-clients/
│ │ │ ├── finance.client.ts
│ │ │ ├── academic.client.ts
│ │ │ └── assessment.client.ts
│ │ ├── event-publisher/
│ │ │ └── pmb-event.publisher.ts
│ │ └── storage/
│ │ └── document-storage.ts
│ │
│ ├── interface/
│ │ ├── controllers/
│ │ ├── validators/
│ │ │ ├── applicant.schema.ts
│ │ │ ├── document.schema.ts
│ │ │ └── loa.schema.ts
│ │ └── presenters/
│ │
│ ├── rbac/
│ │ ├── permissions.ts
│ │ └── scope-policy.ts
│ │
│ ├── audit/
│ │ └── audit.service.ts
│ │
│ ├── idempotency/
│ │ └── idempotency.service.ts
│ │
│ ├── outbox/
│ │ ├── outbox.service.ts
│ │ └── outbox-event.types.ts
│ │
│ ├── inbox/
│ │ └── inbox.service.ts
│ │
│ └── shared/
│ ├── response-envelope.ts
│ ├── error-codes.ts
│ ├── logger.ts
│ └── correlation-id.ts
│
├── prisma/
│ ├── schema.prisma
│ ├── migrations/
│ └── seed.ts
│
├── tests/
│ ├── unit/
│ ├── integration/
│ ├── contract/
│ └── e2e/
│
├── public/
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── package.json
├── tsconfig.json
├── next.config.ts
└── README.md

# 4. Template yang Sama untuk Semua Service

{module}-service/
│
├── app/
│ └── api/
│ └── v1/
│ └── {module}/
│
├── src/
│ ├── domain/
│ ├── application/
│ ├── infrastructure/
│ ├── interface/
│ ├── rbac/
│ ├── audit/
│ ├── idempotency/
│ ├── outbox/
│ ├── inbox/
│ └── shared/
│
├── prisma/
├── tests/
├── Dockerfile
├── docker-compose.yml
├── package.json
└── README.md

Contoh implementasi file di modul Finance:

unsia-finance-service/src/domain/entities/invoice.entity.ts
unsia-finance-service/src/application/commands/process-payment-callback.command.ts
unsia-finance-service/src/infrastructure/repositories/payment.repository.ts

# 5. Struktur Repo unsia-core-service

unsia-core-service/
│
├── app/
│ └── api/
│ └── v1/
│ ├── auth/
│ │ ├── login/
│ │ │ └── route.ts
│ │ ├── refresh/
│ │ │ └── route.ts
│ │ ├── me/
│ │ │ └── route.ts
│ │ └── switch-role/
│ │ └── route.ts
│ │
│ ├── users/
│ ├── roles/
│ ├── permissions/
│ ├── applications/
│ ├── service-tokens/
│ ├── impersonations/
│ └── audit-logs/
│
├── src/
│ ├── domain/
│ │ ├── entities/
│ │ │ ├── user.entity.ts
│ │ │ ├── role.entity.ts
│ │ │ ├── permission.entity.ts
│ │ │ ├── session.entity.ts
│ │ │ └── service-token.entity.ts
│ │ └── value-objects/
│ │
│ ├── application/
│ │ ├── commands/
│ │ ├── queries/
│ │ └── services/
│ │
│ ├── infrastructure/
│ │ ├── database/
│ │ ├── repositories/
│ │ ├── jwt/
│ │ ├── password/
│ │ └── service-token/
│ │
│ ├── rbac/
│ │ ├── permission-checker.ts
│ │ ├── active-role-resolver.ts
│ │ └── scope-resolver.ts
│ │
│ ├── audit/
│ ├── idempotency/
│ ├── outbox/
│ └── shared/
│
├── prisma/
├── tests/
├── Dockerfile
└── README.md

# 6. Struktur Repo unsia-portal-web

Portal berfungsi sebagai frontend utama dan tidak boleh menjadi source transaksi. Portal hanya memanggil API modul sumber.

unsia-portal-web/
│
├── app/
│ ├── (auth)/
│ │ ├── login/
│ │ └── select-role/
│ │
│ ├── (portal)/
│ │ ├── dashboard/
│ │ ├── notifications/
│ │ ├── profile/
│ │ └── applications/
│ │
│ ├── pendaftar/
│ │ ├── biodata/
│ │ ├── documents/
│ │ ├── invoice/
│ │ └── loa/
│ │
│ ├── mahasiswa/
│ │ ├── dashboard/
│ │ ├── krs/
│ │ ├── lms/
│ │ ├── khs/
│ │ └── transcript/
│ │
│ ├── dosen/
│ │ ├── dashboard/
│ │ ├── classes/
│ │ ├── attendance/
│ │ ├── assignments/
│ │ └── grades/
│ │
│ ├── admin/
│ │ ├── pmb/
│ │ ├── finance/
│ │ ├── academic/
│ │ ├── hris/
│ │ ├── lms/
│ │ └── assessment/
│ │
│ └── pimpinan/
│ ├── dashboard/
│ ├── kpi/
│ └── reports/
│
├── src/
│ ├── components/
│ │ ├── layout/
│ │ ├── form/
│ │ ├── table/
│ │ ├── modal/
│ │ └── dashboard/
│ │
│ ├── features/
│ │ ├── auth/
│ │ ├── pmb/
│ │ ├── finance/
│ │ ├── academic/
│ │ ├── lms/
│ │ ├── assessment/
│ │ └── portal/
│ │
│ ├── lib/
│ │ ├── api-client.ts
│ │ ├── auth.ts
│ │ ├── rbac.ts
│ │ ├── query-client.ts
│ │ └── utils.ts
│ │
│ ├── hooks/
│ ├── stores/
│ ├── types/
│ └── styles/
│
├── public/
├── tests/
├── Dockerfile
├── package.json
└── README.md

# 7. Struktur Repo unsia-shared-contracts

Repo ini wajib dijadikan acuan agar semua service menggunakan format API, event, error, permission, dan schema yang sama.

unsia-shared-contracts/
│
├── openapi/
│ ├── core.openapi.yaml
│ ├── reference.openapi.yaml
│ ├── crm.openapi.yaml
│ ├── pmb.openapi.yaml
│ ├── finance.openapi.yaml
│ ├── academic.openapi.yaml
│ ├── hris.openapi.yaml
│ ├── lms.openapi.yaml
│ ├── assessment.openapi.yaml
│ ├── portal.openapi.yaml
│ └── integration.openapi.yaml
│
├── event-contracts/
│ ├── core/
│ ├── reference/
│ ├── crm/
│ ├── pmb/
│ ├── finance/
│ ├── academic/
│ ├── hris/
│ ├── lms/
│ ├── assessment/
│ └── portal/
│
├── schemas/
│ ├── response-envelope.schema.json
│ ├── error-envelope.schema.json
│ ├── event-envelope.schema.json
│ ├── audit-log.schema.json
│ └── idempotency.schema.json
│
├── rbac/
│ ├── roles.yaml
│ ├── permissions.yaml
│ ├── role-permission-matrix.yaml
│ └── data-scope-matrix.yaml
│
├── error-codes/
│ └── error-codes.yaml
│
├── typescript/
│ ├── generated/
│ └── package.json
│
└── README.md

# 8. Struktur Repo unsia-integration-worker

unsia-integration-worker/
│
├── src/
│ ├── workers/
│ │ ├── outbox-publisher.worker.ts
│ │ ├── inbox-consumer.worker.ts
│ │ ├── dlq-replay.worker.ts
│ │ ├── reconciliation.worker.ts
│ │ ├── notification.worker.ts
│ │ └── snapshot-refresh.worker.ts
│ │
│ ├── consumers/
│ │ ├── finance-payment-paid.consumer.ts
│ │ ├── pmb-ready-for-academic.consumer.ts
│ │ ├── academic-student-created.consumer.ts
│ │ ├── academic-krs-approved.consumer.ts
│ │ ├── lms-grade-input-submitted.consumer.ts
│ │ └── assessment-result-calculated.consumer.ts
│ │
│ ├── publishers/
│ │ └── event-bus.publisher.ts
│ │
│ ├── queues/
│ │ ├── rabbitmq.ts
│ │ └── redis.ts
│ │
│ ├── clients/
│ │ ├── core.client.ts
│ │ ├── pmb.client.ts
│ │ ├── finance.client.ts
│ │ ├── academic.client.ts
│ │ ├── lms.client.ts
│ │ └── portal.client.ts
│ │
│ ├── observability/
│ │ ├── logger.ts
│ │ ├── metrics.ts
│ │ └── tracing.ts
│ │
│ └── shared/
│
├── tests/
├── Dockerfile
├── package.json
└── README.md

# 9. Struktur Repo unsia-infra

unsia-infra/
│
├── docker/
│ ├── docker-compose.local.yml
│ ├── docker-compose.staging.yml
│ └── docker-compose.prod.yml
│
├── nginx/
│ ├── api-gateway.conf
│ ├── portal.conf
│ └── services.conf
│
├── postgres/
│ ├── core-db/
│ ├── reference-db/
│ ├── crm-db/
│ ├── pmb-db/
│ ├── finance-db/
│ ├── academic-db/
│ ├── hris-db/
│ ├── lms-db/
│ ├── assessment-db/
│ └── portal-db/
│
├── redis/
│ └── redis.conf
│
├── rabbitmq/
│ ├── definitions.json
│ └── rabbitmq.conf
│
├── monitoring/
│ ├── prometheus/
│ ├── grafana/
│ ├── loki/
│ └── alertmanager/
│
├── ci-cd/
│ ├── github-actions/
│ └── gitlab-ci/
│
├── secrets/
│ └── .env.example
│
├── backup/
│ ├── backup-postgres.sh
│ ├── restore-postgres.sh
│ └── backup-minio.sh
│
└── README.md

# 10. Struktur Repo unsia-docs

unsia-docs/
│
├── 01-prd/
├── 02-brd/
├── 03-fsd/
├── 04-api-contract/
├── 05-event-contract/
├── 06-erd-dbml/
├── 07-rbac/
├── 08-state-machine/
├── 09-uat/
├── 10-release-plan/
├── 11-runbook/
├── 12-architecture-decision-record/
└── README.md

# 11. Naming Branch

main
develop
release/v1.0.0
feature/pmb-applicant-registration
feature/finance-payment-callback
feature/academic-krs-approval
bugfix/pmb-document-upload-validation
hotfix/payment-callback-duplicate

# 12. Naming Commit

Gunakan conventional commit agar riwayat perubahan mudah dibaca dan bisa dipakai untuk changelog otomatis.

feat(pmb): add applicant registration API
fix(finance): prevent duplicate payment callback
refactor(core): improve active role resolver
test(academic): add KRS approval scope test
docs(api): update PMB handover contract
chore(infra): add RabbitMQ config

# 13. Urutan Pembuatan Repo

unsia-shared-contracts

unsia-infra

unsia-core-service

unsia-reference-service

unsia-pmb-service

unsia-finance-service

unsia-academic-service

unsia-portal-web

unsia-lms-service

unsia-assessment-service

unsia-hris-service

unsia-crm-service

unsia-integration-worker

unsia-docs

# 14. Keputusan Final Struktur Repo

UNSIA ERP Git Organization
│
├── unsia-core-service
├── unsia-reference-service
├── unsia-crm-service
├── unsia-pmb-service
├── unsia-finance-service
├── unsia-academic-service
├── unsia-hris-service
├── unsia-lms-service
├── unsia-assessment-service
├── unsia-portal-web
├── unsia-integration-worker
├── unsia-shared-contracts
├── unsia-infra
└── unsia-docs

Aturan utama:

Satu repo mewakili satu modul atau service.

Satu service memiliki satu database utama.

Tidak boleh ada cross-database join untuk transaksi online.

Tidak boleh ada write langsung ke database modul lain.

Semua integrasi lintas modul wajib melalui API contract dan event contract.

Shared contracts wajib menjadi acuan semua repo.

Infra dan docs dipisah dari source code aplikasi.

Catatan implementasi: struktur ini dapat digunakan untuk model full Next.js/Node.js per service, selama setiap service tetap menjaga module boundary, database ownership, RBAC, idempotency, outbox/inbox, dan audit trail.
