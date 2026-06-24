# unsia-assessment-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `assessment_db` (PostgreSQL)

## Tanggung Jawab

Modul Assessment adalah **reusable assessment engine** untuk CBT, quiz, survey, dan scoring — dapat digunakan oleh PMB, LMS, maupun Academic.

| Domain | Deskripsi |
|--------|-----------|
| Bank Soal | Soal dengan tipe, opsi, kunci, dan tingkat kesulitan |
| Versi Soal | Soal yang sudah dipakai attempt tidak bisa diedit — buat versi baru |
| Question Set | Paket soal untuk satu sesi assessment |
| Assessment Session | Sesi CBT/quiz/survey dengan konteks dan jadwal |
| Attempt | Jawaban peserta per sesi |
| Scoring | Auto-scoring dan manual review |
| Survey | Kuesioner tanpa jawaban benar |
| Result API | Publish hasil ke consumer (PMB/LMS/Academic) |

## Endpoint Utama

```
POST   /api/v1/assessment/question-banks
POST   /api/v1/assessment/question-banks/{id}/versions
POST   /api/v1/assessment/sessions
POST   /api/v1/assessment/sessions/{id}/participants
POST   /api/v1/assessment/attempts
POST   /api/v1/assessment/attempts/{id}/submit
POST   /api/v1/assessment/results/publish
GET    /api/v1/assessment/sessions/{id}/results
```

## Struktur Direktori (akan diisi saat development)

```
unsia-assessment-service/
├── cmd/assessment-service/main.go
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

- **Upstream:** Core (auth), PMB/LMS (context owner saat membuat sesi)
- **Event publish:** `assessment.result_calculated`
- **Event consume:** —
- **API call:** PMB / LMS (notify result)

## Aturan Penting

- Soal yang sudah dipakai attempt **tidak boleh diedit langsung** — buat `question_version` baru
- Attempt yang sudah submitted **immutable**
- Result publish wajib **idempotent**
- Assessment result hanya menjadi **input** ke PMB/LMS/Academic — bukan otomatis final decision
- Konteks assessment (CBT PMB vs quiz LMS) dibedakan via `context_type` dan `context_ref_id`

## Owner

Phase 3, Sprint 8.
