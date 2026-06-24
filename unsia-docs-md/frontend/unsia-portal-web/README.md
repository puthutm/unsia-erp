# unsia-portal-web

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

## Tanggung Jawab

Frontend utama ERP UNSIA — portal terpadu untuk semua role. **Tidak menjadi source transaksi bisnis** — hanya memanggil API backend service.

## Struktur Route

```
app/
├── (auth)/
│   ├── login/            → Halaman login (POST ke Core /auth/login)
│   └── select-role/      → Pilih active role setelah login
│
├── (portal)/
│   ├── dashboard/        → Dashboard sesuai active role
│   ├── notifications/    → Notification center
│   ├── profile/          → Profil user
│   └── applications/     → App launcher
│
├── pendaftar/            → Menu khusus calon mahasiswa (self scope)
│   ├── biodata/
│   ├── documents/
│   ├── invoice/
│   └── loa/
│
├── mahasiswa/            → Menu mahasiswa aktif (self scope)
│   ├── dashboard/
│   ├── krs/
│   ├── lms/
│   ├── khs/
│   └── transcript/
│
├── dosen/                → Menu dosen (assigned class scope)
│   ├── dashboard/
│   ├── classes/
│   ├── attendance/
│   ├── assignments/
│   └── grades/
│
├── admin/                → Menu admin per modul
│   ├── pmb/
│   ├── finance/
│   ├── academic/
│   ├── hris/
│   ├── lms/
│   ├── assessment/
│   └── oauth-clients/    → Manajemen OAuth External App
│
├── developer/            → Menu developer OAuth
│   └── oauth-credentials/
│
└── pimpinan/             → Dashboard eksekutif (read-only)
    ├── dashboard/
    ├── kpi/
    └── reports/
```

## Struktur Source

```
src/
├── components/           → Reusable UI components
│   ├── layout/
│   ├── form/
│   ├── table/
│   ├── modal/
│   └── dashboard/
├── features/             → Feature modules (satu folder per modul)
│   ├── auth/
│   ├── pmb/
│   ├── finance/
│   ├── academic/
│   ├── lms/
│   ├── assessment/
│   ├── portal/
│   └── oauth/
├── lib/
│   ├── api-client.ts     → Axios/fetch wrapper dengan auth header
│   ├── auth.ts           → Session management, token refresh
│   ├── rbac.ts           → Client-side permission check
│   └── query-client.ts   → TanStack Query global config
├── hooks/
├── stores/               → Zustand atau Jotai untuk state lokal
└── types/                → Generated dari OpenAPI schema
```

## Standard UI Rules

- Halaman yang menampilkan snapshot/read model **wajib tampilkan `synced_at`**
- Saat service backend down, tampilkan **degraded state** dengan pesan informatif dan `trace_id`
- Form dengan aksi sensitif wajib **dialog konfirmasi** sebelum submit
- Semua list page wajib: pagination, search, filter, empty state, loading state
- Badge status: PENDING (kuning), ACTIVE (hijau), SUSPENDED (oranye), REVOKED (merah)

## Owner

Phase 0 (shell) → Phase 2 Sprint 3 (OAuth UI) → Phase 3 Sprint 8 (semua fitur).
