# unsia-pmb-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `pmb_db` (PostgreSQL)

## Tanggung Jawab

Modul PMB adalah **source of truth applicant** — mengelola lifecycle calon mahasiswa dari pendaftaran hingga handover ke Academic.

| Domain | Deskripsi |
|--------|-----------|
| Applicant | Data calon mahasiswa dan status pendaftaran |
| Biodata | Data lengkap personal, alamat, pendidikan, keluarga |
| Dokumen | Upload, verifikasi, dan penolakan dokumen |
| Seleksi / CBT | Integrasi hasil seleksi dari Assessment |
| Invoice Status | Read model status pembayaran dari Finance |
| Daftar Ulang | Proses re-registration setelah diterima |
| LoA | Penerbitan Letter of Acceptance |
| Handover | Serah terima applicant ke modul Academic (idempotent) |

## Endpoint Utama

```
POST   /api/v1/pmb/applicants
POST   /api/v1/pmb/applicants/{id}/submit
POST   /api/v1/pmb/applicants/{id}/documents
POST   /api/v1/pmb/applicants/{id}/documents/{doc_id}/verify
POST   /api/v1/pmb/applicants/{id}/request-invoice
POST   /api/v1/pmb/applicants/{id}/issue-loa
POST   /api/v1/pmb/applicants/{id}/handover-to-academic
GET    /api/v1/pmb/applicants
GET    /api/v1/pmb/applicants/{id}
```

## Struktur Direktori (akan diisi saat development)

```
unsia-pmb-service/
├── cmd/pmb-service/main.go
├── internal/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   │   ├── external-clients/
│   │   │   ├── finance.client.go
│   │   │   ├── academic.client.go
│   │   │   └── assessment.client.go
│   │   └── ...
│   ├── handler/
│   └── middleware/
├── migrations/
├── tests/
├── Dockerfile
├── .env.example
└── go.mod
```

## Dependencies

- **Upstream:** Core (auth), Reference (snapshot), Finance (invoice/clearance), Assessment (CBT result)
- **Event publish:** `pmb.applicant_created`, `pmb.document_verified`, `pmb.ready_for_academic`, `pmb.handover_requested`
- **Event consume:** `finance.invoice_created`, `finance.payment_paid`, `finance.clearance_changed`, `assessment.result_calculated`
- **API call:** Finance (create invoice), Academic (generate student)

## State Machine Applicant

```
DRAFT → SUBMITTED → VERIFIED → ACCEPTED → LOA_ISSUED → HANDED_OVER
```

## Aturan Penting

- PMB **tidak menyimpan payment** sebagai source of truth — hanya read model dari Finance
- Handover ke Academic wajib **idempotent** (duplicate tidak membuat student/NIM ganda)
- Jika Academic down, handover masuk **retry queue** dengan status `pending_handover`

## Owner

Phase 3, Sprint 5.
