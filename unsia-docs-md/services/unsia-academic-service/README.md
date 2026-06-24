# unsia-academic-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `academic_db` (PostgreSQL)

## Tanggung Jawab

Modul Academic adalah **source of truth mahasiswa, kurikulum, kelas, KRS, nilai final, KHS, transkrip, dan alumni**.

| Domain | Deskripsi |
|--------|-----------|
| Kalender Akademik | Tahun Ajaran dan Periode Akademik operasional |
| Kurikulum | Kurikulum per prodi dengan tahun versi |
| Mata Kuliah | Katalog MK per kurikulum/prodi |
| Mahasiswa & NIM | Generate NIM dari handover PMB (idempotent) |
| Kelas | Penawaran kelas per periode dengan dosen dan kapasitas |
| KRS | Rencana Studi paket (sem 1–2) dan mandiri (sem 3+) |
| Nilai Final | Finalisasi nilai dari input LMS/Assessment |
| KHS & Transkrip | Publikasi nilai akademik periodik |
| Yudisium & Alumni | Proses kelulusan |

## Endpoint Utama

```
POST   /api/v1/academic/students/generate-from-applicant
GET    /api/v1/academic/students
POST   /api/v1/academic/classes
POST   /api/v1/academic/krs
POST   /api/v1/academic/krs/{id}/submit
POST   /api/v1/academic/krs/{id}/approve
POST   /api/v1/academic/grades/source-imports
POST   /api/v1/academic/grades/{id}/finalize
GET    /api/v1/academic/khs/{student_id}
GET    /api/v1/academic/transcripts/{student_id}
```

## Struktur Direktori (akan diisi saat development)

```
unsia-academic-service/
├── cmd/academic-service/main.go
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

- **Upstream:** Core (auth), Reference (prodi, periode), Finance (clearance), HRIS (dosen snapshot), PMB (handover)
- **Event publish:** `academic.student_created`, `academic.class_opened`, `academic.krs_approved`, `academic.final_grade_published`, `academic.khs_issued`
- **Event consume:** `pmb.ready_for_academic`, `finance.clearance_changed`, `hris.lecturer_updated`, `lms.grade_input_submitted`, `assessment.result_calculated`

## State Machine KRS

```
DRAFT → SUBMITTED → APPROVED → FINALIZED → CANCELLED
```

## Aturan Penting

- `students.applicant_ref_id` harus **UNIQUE** — duplicate handover tidak membuat mahasiswa ganda
- NIM harus **unique** dan sequence dikunci saat generate (row-level lock)
- **Final grade** hanya milik Academic — LMS/Assessment hanya menjadi sumber input
- KRS Mandiri (sem 3+) wajib cek **clearance Finance** sebelum difinalisasi
- Tahun Kurikulum ≠ Tahun Ajaran — mahasiswa menyimpan `curriculum_id` saat NIM dibuat

## Owner

Phase 3, Sprint 7.
