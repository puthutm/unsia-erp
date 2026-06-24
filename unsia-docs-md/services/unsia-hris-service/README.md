# unsia-hris-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `hris_db` (PostgreSQL)

## Tanggung Jawab

Modul HRIS adalah **source of truth dosen dan karyawan** untuk seluruh ekosistem ERP UNSIA.

| Domain | Deskripsi |
|--------|-----------|
| Pegawai | Data karyawan dengan unit kerja dan jabatan |
| Dosen | NIDN/NIDK, homebase prodi, jenjang akademik, status aktif |
| Homebase | Penugasan dosen ke prodi utama |
| Unit Kerja & Jabatan | Struktur organisasi dan riwayat jabatan |
| Status Kepegawaian | Aktif, cuti, tidak aktif, pensiun |
| BKD | Beban Kerja Dosen per periode |
| Sertifikasi & Kinerja | Rekam jejak kompetensi dosen |
| Payroll Source | Data dasar penggajian (dibaca Finance) |

## Endpoint Utama

```
GET    /api/v1/hris/lecturers
GET    /api/v1/hris/lecturers/{id}
POST   /api/v1/hris/lecturers
PUT    /api/v1/hris/lecturers/{id}
PATCH  /api/v1/hris/lecturers/{id}/status
GET    /api/v1/hris/employees
POST   /api/v1/hris/employees
GET    /api/v1/hris/bkd
```

## Struktur Direktori (akan diisi saat development)

```
unsia-hris-service/
├── cmd/hris-service/main.go
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

- **Upstream:** Core (person identity), Reference (prodi snapshot)
- **Event publish:** `hris.lecturer_updated`, `hris.lecturer_status_changed`, `hris.employee_updated`
- **Event consume:** `core.person_updated`
- **Downstream consumer:** Academic (lecturer snapshot), LMS (lecturer snapshot)

## Aturan Penting

- Academic dan LMS **tidak boleh** membuat data dosen mandiri — harus menggunakan `lecturer_ref_id` + snapshot dari HRIS
- Dosen **nonaktif** tidak boleh diplot ke kelas baru
- Perubahan data dosen harus **publish event** agar consumer update snapshot lokal

## Owner

Phase 3, Sprint 8.
