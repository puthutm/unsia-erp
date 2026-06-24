# unsia-lms

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-lms-service` → `lms_db`

## Tanggung Jawab

Frontend modul LMS — platform pembelajaran online untuk dosen dan mahasiswa.

## Route Structure

```
app/
├── (mahasiswa)/
│   ├── dashboard/
│   ├── kelas/
│   │   ├── page.tsx               → List kelas aktif semester ini
│   │   └── [classId]/
│   │       ├── page.tsx           → Overview kelas
│   │       ├── sesi/[sessionId]/  → Detail sesi + materi
│   │       ├── tugas/
│   │       │   ├── page.tsx
│   │       │   └── [taskId]/submit/
│   │       ├── diskusi/
│   │       └── presensi/
│   └── progress/
│
└── (dosen)/
    ├── dashboard/
    ├── kelas/
    │   ├── page.tsx               → List kelas yang diajar
    │   └── [classId]/
    │       ├── page.tsx
    │       ├── sesi/
    │       │   ├── page.tsx       → Buat/edit sesi
    │       │   └── [sessionId]/
    │       │       ├── materi/
    │       │       ├── tugas/
    │       │       └── attendance/
    │       ├── diskusi/
    │       ├── peserta/
    │       └── nilai-input/
    └── laporan/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Mahasiswa | Self | Kelas aktif, materi, tugas, diskusi, presensi, progress |
| Dosen | Assigned class | Sesi, materi, tugas, presensi, grade input |

## Integrasi API

- `unsia-lms-service` — semua data LMS
- `unsia-core-service` — auth/token

## Aturan UI

- Kelas dan peserta berasal dari **Academic** via sync — jika sync tertunda, tampilkan label + `synced_at`
- Grade input **bukan nilai final** — tampilkan disclaimer "Nilai final ditetapkan oleh Admin Akademik"
- Upload materi: tampilkan progress upload + validasi tipe file
- Presensi: QR code atau input manual per mahasiswa
