# unsia-finance-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `finance_db` (PostgreSQL)

## Tanggung Jawab

Modul Finance adalah **source of truth transaksi keuangan dan clearance** layanan akademik.

| Domain | Deskripsi |
|--------|-----------|
| Invoice | Tagihan per komponen pembayaran |
| Payment | Pembayaran via gateway atau manual |
| Payment Callback | Callback dari payment gateway (idempotent) |
| Manual Verification | Verifikasi bukti transfer manual |
| Receipt | Bukti pembayaran resmi |
| Clearance | Status kelayakan layanan akademik per mahasiswa per periode |
| Beasiswa & Diskon | Potongan tagihan dengan approval |
| Cicilan / Dispensasi | Jadwal pembayaran bertahap |

## Endpoint Utama

```
POST   /api/v1/finance/invoices
GET    /api/v1/finance/invoices/{id}
POST   /api/v1/finance/payment-callbacks/{provider}
POST   /api/v1/finance/payment-verifications
POST   /api/v1/finance/receipts/{payment_id}
GET    /api/v1/finance/clearances
POST   /api/v1/finance/clearances/override
POST   /api/v1/finance/scholarships
```

## Struktur Direktori (akan diisi saat development)

```
unsia-finance-service/
├── cmd/finance-service/main.go
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

- **Upstream:** Core (auth), Reference (payment component, periode)
- **Event publish:** `finance.invoice_created`, `finance.payment_paid`, `finance.clearance_changed`
- **Event consume:** `academic.student_created`, `pmb.applicant_created`
- **API call:** —

## State Machine Invoice

```
DRAFT → ISSUED → PARTIALLY_PAID → PAID → CANCELLED / EXPIRED
```

## State Machine Clearance

```
BLOCKED → CONDITIONAL → CLEARED → REVOKED
```

## Aturan Penting

- Payment callback wajib validasi **signature provider** dan idempotent berdasarkan `provider_event_id`
- Clearance adalah **satu-satunya source of truth** untuk status kelayakan layanan akademik
- Modul lain (PMB, Academic, LMS) hanya boleh **membaca clearance via API atau snapshot event**

## Owner

Phase 3, Sprint 6.
