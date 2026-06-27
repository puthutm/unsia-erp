# UNSIA Backend Development Plan

## Overview
Dokumen ini merangkumi rencana pengembangan backend ERP UNSIA yang berbasis **Go (Golang) Microservices** menggunakan framework **Gin Web Framework** dan ORM **GORM**. Seluruh layanan backend didesain menggunakan pendekatan **Database-First** dengan database PostgreSQL terpisah untuk setiap modul guna menjamin isolation dan scalability.

---

## Modul & Arsitektur Layanan

Setiap service di bawah folder `services/` mengikuti arsitektur bersih (Clean Architecture) standar:

```
services/unsia-{service-name}-service/
├── cmd/
│   └── {service-name}-service/
│       └── main.go              ← Inisialisasi DB, routing, middleware, & server start
├── internal/
│   ├── domain/                  ← Model GORM, Value Objects, & interface repository
│   ├── application/             ← Logika bisnis utama / usecases
│   ├── infrastructure/
│   │   └── repository/          ← Akses data / query database
│   ├── handler/                 ← API HTTP handler & JSON request/response binders
│   └── middleware/              ← Interceptors (Auth, RBAC, Logger, dll.)
├── migrations/                  ← SQL script migration (golang-migrate)
├── go.mod
└── Dockerfile
```

---

## Daftar Service & Status Implementasi

Berdasarkan status saat ini (merujuk pada `README.md` dan `API_UNPLANNED.md`), berikut adalah rincian fungsionalitas dan status dari 10 microservices utama:

| Service | Port | Database | Cakupan Fitur Utama | Status Handlers |
|---------|------|----------|---------------------|-----------------|
| **Core Service** | 8001 | `core_db` | Autentikasi SSO, JWT, active role switcher, manajemen User & Role, OAuth Client registry, Audit Log writer | ✅ DONE |
| **Reference Service** | 8007 | `reference_db` | Data master Wilayah (provinsi/kota), Program Studi, Tahun Ajaran, Periode Akademik, Metode Pembayaran | 🔄 PARTIALLY DONE |
| **PMB Service** | 8003 | `pmb_db` | Registrasi Calon Mahasiswa (biodata, dokumen, seleksi masuk, tagihan biaya pendaftaran, handover data ke akademik) | ✅ DONE |
| **Finance Service** | 8002 | `finance_db` | Manajemen Invoice tagihan kuliah, catatan pembayaran, integrasi payment gateway callback, approval cicilan, student clearance, payroll karyawan | ✅ DONE |
| **Academic Service** | 8004 | `academic_db` | Data Mahasiswa aktif, kurikulum, mata kuliah, jadwal kelas & dosen, Kartu Rencana Studi (KRS), Kartu Hasil Studi (KHS), nilai akhir | ✅ DONE |
| **LMS Service** | 8005 | `lms_db` | Sinkronisasi kelas, materi kuliah, penugasan tugas & kuis, presensi, sync nilai ke SIAKAD akademik | 🔄 PARTIALLY DONE |
| **Assessment Service** | 8006 | `assessment_db` | Bank soal ujian, paket soal (question sets), sesi ujian CBT (Computer Based Test), attempts, dan automatic scoring | 🔄 PARTIALLY DONE |
| **HRIS Service** | 8007 | `hris_db` | Kepegawaian, struktur organisasi unit kerja, record absensi harian, cuti pegawai, Beban Kerja Dosen (BKD) | 🔄 PARTIALLY DONE |
| **CRM Service** | 8008 | `crm_db` | Pengelolaan data Leads calon pendaftar, referral agent, pelacakan kampanye marketing, perhitungan komisi referral | 🔄 PARTIALLY DONE |
| **Portal Service** | 8010 | `portal_db` | Dashboard widget, pusat notifikasi real-time pengguna, user preferences, menu shortcuts, dan dynamic sidebar menu | 🔄 PARTIALLY DONE |

---

## shared-packages (Shared Go Modules)

Untuk menghindari redundansi kode, fungsionalitas lintas-layanan diekstrak ke dalam modul bersama yang terletak di folder `/packages`:

1.  **`shared-auth`**: Autentikasi token JWT RS256, caching JWKS public key dari Core Service untuk verifikasi downstream.
2.  **`shared-rbac`**: Pengecekan otorisasi berbasis peran (Role-Based Access Control) dan scope data (global, prodi, self).
3.  **`shared-errorenvelope`**: Standardisasi format JSON response untuk sukses dan error (e.g. format: `{ success: false, error: { code: "DB_ERROR", message: "..." } }`).
4.  **`shared-event`**: Implementasi pola Transactional Outbox dan Inbox untuk menjamin pengiriman event antar microservices melalui message broker (RabbitMQ).
5.  **`shared-idempotency`**: Middleware penanganan request ganda (idempotency key) pada endpoint mutasi data sensitif (misal pembayaran & submit KRS).
6.  **`shared-audit`**: Logging aktivitas audit otomatis ke database untuk pelacakan compliance sistem.
7.  **`shared-observability`**: Integrasi structured logging (ZeroLog), tracing ID, dan Prometheus metrics collector.

---

## Roadmap Pengembangan Backend (To-Do List)

### Phase 1: High Priority (LMS & Assessment CBT)
*   **LMS Service (Port 8005)**
    *   [ ] Menyelesaikan endpoint mata kuliah: `GET /api/v1/lms/courses` dan `POST /api/v1/lms/courses`
    *   [ ] Menyelesaikan endpoint kelas & pendaftaran mahasiswa: `GET /api/v1/lms/enrollments` dan `POST /api/v1/lms/enrollments`
    *   [ ] Menyelesaikan endpoint sesi perkuliahan & upload materi: `POST /api/v1/lms/sessions` dan `POST /api/v1/lms/materials`
    *   [ ] Menambahkan handler tugas: `GET/POST /api/v1/lms/assignments`
*   **Assessment Service (Port 8006)**
    *   [ ] Menyelesaikan manajemen sesi CBT: `GET/POST /api/v1/assessment-sessions`
    *   [ ] Menyelesaikan registrasi peserta ujian: `POST /api/v1/participants`
    *   [ ] Menyelesaikan endpoint CBT attempt: `POST /api/v1/attempts` (memulai ujian) dan `POST /api/v1/attempts/:id/submit` (selesai ujian)

### Phase 2: Medium Priority (HRIS & CRM)
*   **HRIS Service (Port 8007)**
    *   [ ] Menyelesaikan CRUD karyawan dan dosen: `GET/POST/PUT /api/v1/employees`
    *   [ ] Menyelesaikan absensi kepegawaian: `POST /api/v1/attendances`
    *   [ ] Menyelesaikan pengajuan cuti: `POST /api/v1/leave-requests`
    *   [ ] Mengimplementasikan modul BKD & performance review
*   **CRM Service (Port 8008)**
    *   [ ] Menyelesaikan endpoint manajemen kontak Leads: `GET/POST /api/v1/contacts`
    *   [ ] Menyelesaikan pipeline opportunities pendaftaran: `GET/POST /api/v1/opportunities`
    *   [ ] Menyelesaikan campaign tracking: `GET/POST /api/v1/campaigns`

### Phase 3: Low Priority (Reference & Portal Enhancement)
*   **Reference Service (Port 8009)**
    *   [ ] Menyelesaikan data wilayah lengkap: `GET /api/v1/regencies`, `GET /api/v1/districts`, `GET /api/v1/villages`
*   **Portal Service (Port 8010)**
    *   [ ] Menyelesaikan endpoint pengumuman global: `GET/POST /api/v1/portal/announcements`
    *   [ ] Menyelesaikan berita internal: `GET/POST /api/v1/portal/news`
    *   [ ] Menyelesaikan event kalender: `GET/POST /api/v1/portal/events`
    *   [ ] Mengintegrasikan Dynamic Sidebar Menu (`GET /api/v1/portal/menus`) ke dashboard

---

## Keputusan Desain & Integrasi Utama

1.  **Database Isolation**: Tidak boleh melakukan query lintas database PostgreSQL (`join`). Data referensi dari service lain wajib dikonsumsi melalui API call (dengan HTTP Client cache) atau dikirim secara asinkron menggunakan event/Outbox pattern ke Read-Model (Snapshot) lokal.
2.  **JWT Autentikasi Downstream**: Setiap request masuk ke Gateway (Nginx) akan diteruskan ke service masing-masing. Service downstream memvalidasi JWT menggunakan JWKS client cache secara in-memory untuk meningkatkan performa (tanpa perlu hit Core Service terus-menerus).
3.  **Transactional Outbox**: Mutasi state bisnis penting wajib menulis entri ke tabel `outbox_events` di dalam satu transaksi DB yang sama. Integration worker kemudian akan membaca tabel outbox secara periodik untuk diteruskan ke RabbitMQ.

---

## Environment Variables Utama (File `.env` Per Service)

```env
# Database Configuration
DATABASE_URL=postgres://postgres:@localhost:5432/{db_name}?sslmode=disable

# Service Execution Port
PORT={port_number}

# Authentication (JWKS Verification)
JWKS_URL=http://localhost:8001/.well-known/jwks.json

# Redis Configuration (For caching & rate limiting)
REDIS_URL=redis://localhost:6379/0

# Message Broker (RabbitMQ)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```
