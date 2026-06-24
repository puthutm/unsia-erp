# unsia-crm-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `crm_db` (PostgreSQL)

## Tanggung Jawab

Modul CRM mengelola **akuisisi dan konversi peminat** sebelum menjadi applicant PMB.

| Domain | Deskripsi |
|--------|-----------|
| Campaign | Kampanye pemasaran dengan periode dan budget |
| Lead / Peminat | Data calon pendaftar dari berbagai sumber |
| Agent & Mitra | Agen referral dengan scope data own-lead |
| Follow-up | Histori aktivitas follow-up per lead |
| Konversi | Convert lead qualified ke applicant PMB (idempotent) |
| Komisi | Perhitungan dan approval komisi agen |

## Endpoint Utama

```
GET    /api/v1/crm/leads
POST   /api/v1/crm/leads
POST   /api/v1/crm/leads/{id}/convert-to-applicant
GET    /api/v1/crm/campaigns
POST   /api/v1/crm/campaigns
GET    /api/v1/crm/agents
POST   /api/v1/crm/commissions
```

## Struktur Direktori (akan diisi saat development)

```
unsia-crm-service/
├── cmd/crm-service/main.go
├── internal/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── handler/
│   └── middleware/
├── migrations/
├── tests/
├── Dockerfile
├── .env.example
└── go.mod
```

## Dependencies

- **Upstream:** Core (auth), Reference (snapshot prodi)
- **Event publish:** `crm.lead_qualified`, `crm.lead_converted`
- **Event consume:** `core.person_updated`
- **API call:** PMB Service untuk convert lead → applicant

## Aturan Penting

- Agent hanya bisa melihat lead/referral **miliknya sendiri** (data scope `own_lead`)
- Konversi lead ke applicant wajib **idempotent** — duplicate convert tidak membuat applicant ganda
- CRM menyimpan `applicant_ref_id` setelah PMB berhasil membuat applicant

## Owner

Phase 3, Sprint 5.
