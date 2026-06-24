# unsia-reference

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-reference-service` → `reference_db`

## Tanggung Jawab

Frontend modul Referensi — manajemen master data yang digunakan lintas modul ERP UNSIA.

## Route Structure

```
app/
└── (admin)/
    ├── dashboard/
    ├── prodi/
    │   ├── page.tsx               → List prodi aktif + filter
    │   ├── buat/
    │   └── [id]/
    ├── tahun-ajaran/
    │   ├── page.tsx
    │   └── [id]/
    │       └── periode/           → Daftar periode akademik dalam tahun ini
    ├── periode-akademik/
    ├── status-codes/
    │   ├── page.tsx               → Group by domain
    │   └── [domain]/
    ├── komponen-pembayaran/
    ├── metode-pembayaran/
    └── jenis-dokumen/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Admin Referensi | Global referensi | Semua master data |
| Admin Akademik Biro | Terbatas | Tahun Ajaran & Periode Akademik |
| Admin Finance | Terbatas | Komponen & Metode Pembayaran |
| Admin PMB | Terbatas | Jenis Dokumen PMB |

## Integrasi API

- `unsia-reference-service` — semua master data
- `unsia-core-service` — auth/token

## Aturan UI

- Master data yang **sudah dipakai transaksi** tidak bisa dihapus permanen — hanya nonaktifkan
- Tombol hapus hanya muncul jika data belum pernah dipakai di transaksi apapun
- Perubahan master data sensitif tampilkan **warning**: "Perubahan ini akan dipublish ke semua modul consumer"
- Tahun Ajaran dan Tahun Kurikulum adalah entitas berbeda — tampilkan label jelas di UI
