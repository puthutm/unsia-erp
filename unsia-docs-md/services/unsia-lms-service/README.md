# unsia-lms-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `lms_db` (PostgreSQL)

## Tanggung Jawab

Modul LMS mengelola **delivery pembelajaran online** berbasis kelas dan KRS dari Academic.

| Domain | Deskripsi |
|--------|-----------|
| LMS Class | Kelas online hasil sync dari kelas Academic |
| Enrollment | Peserta kelas hasil sync dari KRS valid |
| Sesi | Jadwal dan pertemuan kelas |
| Materi | File, video, link per sesi |
| Tugas & Submission | Penugasan mahasiswa dan pengumpulan |
| Presensi | Kehadiran per sesi |
| Diskusi | Forum diskusi per kelas |
| Progress | Tracking progress belajar mahasiswa |
| Grade Input | Nilai aktivitas (bukan final grade) dikirim ke Academic |

## Endpoint Utama

```
POST   /api/v1/lms/classes/sync-from-academic
POST   /api/v1/lms/enrollments/sync-from-krs
GET    /api/v1/lms/classes/{id}
POST   /api/v1/lms/classes/{id}/sessions
POST   /api/v1/lms/sessions/{id}/materials
POST   /api/v1/lms/sessions/{id}/assignments
POST   /api/v1/lms/assignments/{id}/submissions
POST   /api/v1/lms/sessions/{id}/attendance
POST   /api/v1/lms/grade-syncs
```

## Struktur Direktori (akan diisi saat development)

```
unsia-lms-service/
├── cmd/lms-service/main.go
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

- **Upstream:** Core (auth), Academic (class/enrollment sync), HRIS (lecturer snapshot)
- **Event publish:** `lms.progress_updated`, `lms.grade_input_submitted`
- **Event consume:** `academic.class_opened`, `academic.krs_approved`, `hris.lecturer_updated`, `finance.clearance_changed`
- **API call:** Academic (grade sync)

## Aturan Penting

- LMS **tidak membuat kelas akademik sendiri** — semua kelas berasal dari Academic
- Enrollment **hanya berasal dari KRS valid** yang sudah diapprove
- **Final grade tetap milik Academic** — LMS hanya mengirim grade input
- Class sync dan enrollment sync wajib **idempotent**
- Kelas yang sudah tersinkron tetap berjalan meskipun Academic sementara down (degraded mode)

## Owner

Phase 3, Sprint 8.
