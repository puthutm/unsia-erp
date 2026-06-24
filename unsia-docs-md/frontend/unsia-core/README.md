# unsia-core

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-core-service` → `core_db`

## Tanggung Jawab

Frontend modul Core — Super Admin panel untuk manajemen user, role, permission, aplikasi, service token, impersonation, audit, dan OAuth client registry.

## Route Structure

```
app/
├── (auth)/                        → Halaman login (semua role masuk sini)
│   ├── login/
│   ├── select-role/
│   └── forgot-password/
│
└── (super-admin)/
    ├── dashboard/
    ├── users/
    │   ├── page.tsx               → List user
    │   ├── buat/
    │   └── [id]/
    │       ├── roles/
    │       └── impersonate/
    ├── roles/
    │   ├── page.tsx
    │   └── [id]/
    │       └── permissions/
    ├── permissions/
    ├── data-scope/
    ├── applications/              → Application registry
    ├── service-tokens/
    ├── impersonation-logs/
    ├── audit-logs/
    │
    ├── oauth-clients/             → SSO External App Management
    │   ├── page.tsx               → List + filter by status
    │   └── [id]/
    │       ├── page.tsx           → Detail client + histori status
    │       ├── approve/
    │       ├── suspend/
    │       └── revoke/
    │
    └── integration-logs/          → Event log, retry, DLQ viewer
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Super Admin / Admin BPPTI | Global | User, role, permission, OAuth, audit |

## Halaman Khusus OAuth Developer

```
app/(developer)/
└── oauth-credentials/             → Developer lihat client_id & one-time secret
```

## Integrasi API

- `unsia-core-service` — semua endpoint auth, user, role, OAuth
- Login dari halaman ini digunakan oleh **semua role** di semua modul FE

## Aturan UI

- OAuth client **PENDING**: badge kuning, tampilkan tombol Approve + Reject
- OAuth client **ACTIVE**: badge hijau, tampilkan tombol Suspend
- OAuth client **SUSPENDED**: badge oranye, tampilkan tombol Revoke + Reactivate
- OAuth client **REVOKED**: badge merah, no action
- `client_secret` tampil satu kali saja saat approval — tidak bisa dilihat ulang
- Semua aksi admin OAuth wajib **dialog konfirmasi + reason field**
- Audit log: filter by actor, action, module, date range — export CSV
