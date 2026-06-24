# unsia-hris

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-hris-service` → `hris_db`

## Tanggung Jawab

Frontend modul SDM — manajemen data dosen, karyawan, jabatan, homebase, BKD, dan status kepegawaian.

## Route Structure

```
app/
├── (dosen)/                       → Dosen (self scope)
│   ├── profil/
│   ├── bkd/
│   └── sertifikasi/
│
└── (admin)/                       → Admin SDM
    ├── dashboard/
    ├── dosen/
    │   ├── page.tsx               → List dosen + filter by prodi, status
    │   ├── buat/
    │   └── [id]/
    │       ├── page.tsx           → Detail + riwayat jabatan
    │       ├── homebase/
    │       ├── bkd/
    │       └── sertifikasi/
    ├── karyawan/
    │   ├── page.tsx
    │   └── [id]/
    ├── unit-kerja/
    ├── jabatan/
    └── payroll-source/            → Hanya bisa dibaca Finance
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Dosen | Self | Profil, BKD, sertifikasi |
| Admin SDM | HRIS domain | Semua data dosen & karyawan |

## Integrasi API

- `unsia-hris-service` — semua data SDM
- `unsia-core-service` — auth/token (person identity)

## Aturan UI

- Dosen hanya bisa dihubungkan ke person yang sudah terdaftar di **Core**
- Ubah status ke nonaktif wajib **dialog konfirmasi + reason + tanggal efektif**
- Homebase prodi menggunakan dropdown dari **Reference snapshot**
- Payroll source hanya bisa diekspor ke Finance — tidak bisa diedit dari UI Finance
