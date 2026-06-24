# unsia-assessment

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-assessment-service` → `assessment_db`

## Tanggung Jawab

Frontend modul Assessment — CBT, quiz, survey, bank soal, dan laporan hasil assessment.

## Route Structure

```
app/
├── (peserta)/                     → Mahasiswa / Pendaftar saat mengerjakan
│   └── sesi/[sessionId]/
│       ├── page.tsx               → Start attempt (tampilkan waktu, soal)
│       ├── soal/[questionIndex]/  → Tampilkan soal satu per satu / semua
│       └── selesai/               → Konfirmasi submit + summary
│
└── (admin)/                       → Admin Assessment + Dosen
    ├── dashboard/
    ├── bank-soal/
    │   ├── page.tsx               → List + search + filter by tipe
    │   ├── buat/
    │   └── [id]/
    │       └── versi/             → Riwayat versi soal
    ├── question-set/
    ├── sesi/
    │   ├── page.tsx               → List sesi CBT/quiz/survey
    │   ├── buat/
    │   └── [id]/
    │       ├── peserta/
    │       └── hasil/
    └── laporan/
        ├── distribusi-skor/
        └── hasil-per-peserta/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Peserta (PMB / Mahasiswa) | Self | Kerjakan CBT/quiz, lihat hasil |
| Dosen | Assigned session | Buat soal, kelola sesi quiz, lihat hasil |
| Admin Assessment | Assessment domain | Bank soal, sesi CBT, scoring, laporan |

## Integrasi API

- `unsia-assessment-service` — semua data assessment
- `unsia-core-service` — auth/token

## Aturan UI

- Timer CBT **full screen + tidak bisa di-pause** setelah mulai
- Soal yang sudah terpakai attempt tampilkan versi baru jika diedit
- Hasil assessment tampilkan disclaimer: "Hasil ini merupakan input — bukan nilai akhir resmi"
- Anti-cheat basic: deteksi tab switching (opsional untuk fase awal)
