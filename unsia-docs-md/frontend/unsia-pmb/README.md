# unsia-pmb

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-pmb-service` → `pmb_db`

## Tanggung Jawab

Frontend khusus modul PMB — mengelola lifecycle calon mahasiswa dari pendaftaran publik hingga handover ke Akademik.

## Route Structure

```
app/
├── (public)/
│   └── daftar/                    → Pendaftaran publik (tanpa login)
│       ├── page.tsx               → Pilih gelombang & prodi
│       └── [applicantId]/
│           ├── biodata/
│           ├── dokumen/
│           ├── invoice/
│           └── seleksi/
│
├── (pendaftar)/                   → Applicant dashboard (self scope)
│   ├── dashboard/
│   ├── biodata/
│   ├── dokumen/
│   ├── invoice/
│   ├── seleksi/
│   └── loa/
│
└── (admin)/                       → Admin PMB
    ├── dashboard/
    ├── applicants/
    │   ├── page.tsx               → List + filter + search
    │   └── [id]/
    │       ├── page.tsx           → Detail applicant
    │       ├── dokumen/
    │       ├── seleksi/
    │       └── handover/
    ├── gelombang/
    └── laporan/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Pendaftar | Self | Isi biodata, upload dokumen, cek invoice, download LoA |
| Admin PMB | PMB domain | Verifikasi applicant, dokumen, seleksi, LoA, handover |

## Integrasi API

- `unsia-pmb-service` — endpoint utama PMB
- `unsia-finance-service` — read-only invoice/payment status
- `unsia-assessment-service` — hasil CBT/seleksi
- `unsia-core-service` — auth/token

## Aturan UI

- Status pembayaran berasal dari **snapshot Finance** — tampilkan `synced_at`
- Handover button hanya muncul jika status `LOA_ISSUED` dan payment policy valid
- Form biodata panjang: gunakan **save draft** otomatis tiap 30 detik
