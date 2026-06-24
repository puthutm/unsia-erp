# unsia-reference-service

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `reference_db` (PostgreSQL)

## Tanggung Jawab

Modul Referensi adalah **master data authority** untuk seluruh ekosistem ERP UNSIA. Menyediakan data standar yang dipakai lintas modul agar dropdown, status, periode, dan kode bisnis konsisten.

| Domain | Deskripsi |
|--------|-----------|
| Program Studi | Kode, nama, jenjang, fakultas, status |
| Tahun Ajaran | Kalender operasional kampus |
| Periode Akademik | Semester di bawah tahun ajaran |
| Status Code | Managed enum untuk status bisnis lintas modul |
| Komponen Pembayaran | Jenis tagihan yang digunakan Finance |
| Metode Pembayaran | Provider payment gateway |
| Jenis Dokumen PMB | Aturan dokumen per jalur/prodi |
| Wilayah & Agama | Data referensi demografis |

## Endpoint Utama

```
GET    /api/v1/ref/study-programs
POST   /api/v1/ref/study-programs
GET    /api/v1/ref/academic-years
POST   /api/v1/ref/academic-years
GET    /api/v1/ref/academic-periods
POST   /api/v1/ref/academic-periods
GET    /api/v1/ref/status-codes
GET    /api/v1/ref/payment-components
GET    /api/v1/ref/document-types
```

## Struktur Direktori (akan diisi saat development)

```
unsia-reference-service/
├── cmd/reference-service/main.go
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

- **Upstream:** Core Service (auth/token validation)
- **Event publish:** `reference.study_program_updated`, `reference.academic_period_updated`, `reference.master_data_changed`
- **Consumer:** PMB, Finance, Academic, HRIS, LMS, Assessment, Portal

## Aturan Penting

- Master data yang sudah dipakai transaksi **tidak boleh hard delete** — gunakan `inactive` atau `archived`
- Perubahan master data wajib **publish event** agar consumer dapat update snapshot lokal
- Tahun Kurikulum ≠ Tahun Ajaran — keduanya entitas terpisah

## Owner

Harus selesai di **Phase 3, Sprint 4** sebelum modul PMB, Finance, Academic dapat berjalan penuh.
