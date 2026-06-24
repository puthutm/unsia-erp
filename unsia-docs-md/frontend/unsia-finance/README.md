# unsia-finance

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-finance-service` → `finance_db`

## Tanggung Jawab

Frontend modul Keuangan — mengelola invoice, pembayaran, klarifikasi, clearance, dan laporan keuangan.

## Route Structure

```
app/
├── (mahasiswa)/                   → Mahasiswa / Pendaftar (self scope)
│   ├── tagihan/
│   │   ├── page.tsx               → Daftar invoice aktif
│   │   └── [invoiceId]/
│   │       ├── page.tsx           → Detail invoice + status pembayaran
│   │       └── bukti/             → Upload bukti transfer manual
│   └── riwayat/                   → Histori pembayaran
│
└── (admin)/                       → Admin Keuangan
    ├── dashboard/
    ├── invoices/
    │   ├── page.tsx               → List + filter by status, prodi, periode
    │   └── [id]/
    ├── payments/
    │   ├── callback-logs/         → Log payment gateway callback
    │   └── verifikasi/            → Antrian verifikasi manual
    ├── clearance/
    │   ├── page.tsx               → List clearance per mahasiswa per periode
    │   └── override/
    ├── beasiswa/
    ├── cicilan/
    └── laporan/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Mahasiswa / Pendaftar | Self | Lihat tagihan, upload bukti, download kwitansi |
| Admin Keuangan | Finance domain | Invoice, verifikasi, clearance, laporan |
| Pimpinan | Read-only aggregate | Dashboard KPI keuangan |

## Integrasi API

- `unsia-finance-service` — semua operasi keuangan
- `unsia-core-service` — auth/token

## Aturan UI

- Clearance **source of truth Finance** — modul lain hanya bisa baca via snapshot
- Callback log menampilkan `provider_event_id`, `status`, `retry_count`, `last_error`
- Override clearance wajib **dialog konfirmasi + reason**
- Dashboard tampilkan `refreshed_at` untuk data agregat
