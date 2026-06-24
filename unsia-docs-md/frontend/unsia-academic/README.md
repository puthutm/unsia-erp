# unsia-academic

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-academic-service` → `academic_db`

## Tanggung Jawab

Frontend modul Akademik — mengelola kalender, kurikulum, mahasiswa, kelas, KRS, nilai, KHS, transkrip, dan alumni.

## Route Structure

```
app/
├── (mahasiswa)/                   → Mahasiswa (self scope)
│   ├── krs/
│   │   ├── page.tsx               → KRS aktif periode ini
│   │   ├── paket/                 → KRS paket (sem 1–2)
│   │   └── mandiri/               → KRS mandiri (sem 3+)
│   ├── nilai/
│   ├── khs/
│   └── transkrip/
│
├── (dosen-pa)/                    → Dosen PA (advisor scope)
│   ├── mahasiswa-bimbingan/
│   └── approval-krs/
│
└── (admin)/                       → Admin Akademik
    ├── dashboard/
    ├── kalender/
    │   ├── tahun-ajaran/
    │   └── periode-akademik/
    ├── kurikulum/
    │   ├── page.tsx
    │   └── [id]/mata-kuliah/
    ├── mahasiswa/
    │   ├── page.tsx               → List + filter by prodi, angkatan
    │   └── [nim]/
    ├── kelas/
    │   ├── page.tsx
    │   └── [id]/
    ├── krs/
    │   └── monitoring/
    ├── nilai/
    │   ├── input-review/
    │   └── koreksi/
    ├── khs/
    ├── transkrip/
    └── yudisium/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Mahasiswa | Self | KRS, nilai, KHS, transkrip |
| Dosen PA | Advisor scope | Approval KRS mahasiswa bimbingan |
| Admin Akademik Biro | Global | Kalender, mahasiswa, kelas, nilai, transkrip |
| Kaprodi / Admin Prodi | `study_program_id` | Kurikulum prodi, kelas prodi, monitoring |

## Integrasi API

- `unsia-academic-service` — semua data akademik
- `unsia-finance-service` — clearance snapshot (read-only)
- `unsia-core-service` — auth/token

## Aturan UI

- KRS Mandiri (sem 3+): tampilkan **clearance status** sebelum submit
- Jika clearance tidak tersedia real-time, tampilkan label **"Data clearance mungkin tidak terbaru"** + `synced_at`
- Koreksi nilai wajib dialog **reason** + approval flow
- Transkrip hanya bisa didownload jika clearance wisuda valid
