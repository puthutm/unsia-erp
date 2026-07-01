# Product Requirement Document (PRD) Global - UNSIA ERP (Developer Edition)

## 1. Desain Arsitektur & Pola Integrasi Teknis
Platform ERP UNSIA beroperasi menggunakan arsitektur microservices berbasis pesan asinkron. Untuk menghindari inkonsistency data lintas database fisik modular (*distributed modular databases*), developer wajib mematuhi panduan pola integrasi berikut:

### A. Pola Outbox Transaksional (Transactional Outbox Pattern)
Setiap service yang melakukan mutasi data penting wajib menulis pesan event ke tabel outbox lokal (`integration_event_logs` / `outbox`) di dalam transaksi database yang sama dengan mutasi bisnis utama.
* **Tujuan**: Menjamin pengiriman pesan event minimal satu kali (*At-least-once delivery*) tanpa dipengaruhi kegagalan koneksi broker pesan (RabbitMQ/Kafka).
* **Skema Tabel Outbox**:
  ```sql
  CREATE TABLE integration_event_logs (
      event_id UUID PRIMARY KEY,
      event_name VARCHAR(100) NOT NULL,
      event_key VARCHAR(100) UNIQUE NOT NULL, -- Format: aggregate_type:uuid:event_version
      payload JSONB NOT NULL,
      status VARCHAR(20) NOT NULL, -- PENDING, PUBLISHED, FAILED
      retry_count INT DEFAULT 0,
      last_error TEXT,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      processed_at TIMESTAMP
  );
  ```

### B. Pola Inbox Idempotent (Idempotent Inbox Pattern)
Service penerima (consumer) wajib mencatat ID event yang berhasil diproses ke dalam tabel inbox lokal sebelum mengeksekusi logika bisnis untuk menghindari pemrosesan ganda (*At-most-once processing*).
* **Skema Tabel Inbox**:
  ```sql
  CREATE TABLE integration_inbox_logs (
      event_id UUID PRIMARY KEY,
      event_key VARCHAR(100) UNIQUE NOT NULL,
      status VARCHAR(20) NOT NULL, -- PROCESSED, FAILED
      processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
  ```

---

## 2. Event Catalog & JSON Payload Schema
Berikut adalah spesifikasi JSON schema payload untuk event-event kritis antar-modul:

### A. Event Pelunasan Tagihan (`finance.payment_paid`)
* **Publisher**: Finance Service
* **Consumer**: PMB, Academic, Portal
* **JSON Payload**:
  ```json
  {
    "event_id": "4a71b12b-3129-4b2a-886f-ebf995831201",
    "event_name": "finance.payment_paid",
    "event_version": "1.0.0",
    "event_key": "invoice:inv-2026-000451:paid",
    "timestamp": "2026-07-01T10:00:00Z",
    "data": {
      "payment_id": "8c12a87c-12b4-4b51-912f-2cb1a581e281",
      "invoice_id": "01b2a75c-31a4-4a21-998f-9a1b2c8a7f1a",
      "invoice_number": "INV-2026-000451",
      "bill_to_type": "applicant", -- applicant atau student
      "bill_to_ref_id": "9a1b2c8a-7f1a-4a21-998f-01b2a75c31a4",
      "payment_method": "va_bni",
      "amount_paid": 5000000.00,
      "paid_at": "2026-07-01T09:59:58Z"
    }
  }
  ```

### B. Event Handover Mahasiswa Baru (`pmb.ready_for_academic`)
* **Publisher**: PMB Service
* **Consumer**: Academic Service
* **JSON Payload**:
  ```json
  {
    "event_id": "e2ba8712-412f-412b-88ff-acbfa5129023",
    "event_name": "pmb.ready_for_academic",
    "event_version": "1.0.0",
    "event_key": "applicant:9a1b2c8a-7f1a-4a21-998f-01b2a75c31a4:handover",
    "timestamp": "2026-07-01T10:15:00Z",
    "data": {
      "applicant_id": "9a1b2c8a-7f1a-4a21-998f-01b2a75c31a4",
      "person_id": "1c7f1a2a-998f-4a21-ebf9-01b2a75c31a4",
      "full_name": "Ahmad Dani",
      "identity_number": "3174092801990002",
      "email": "ahmad.dani@example.com",
      "phone": "+6281234567890",
      "study_program_id": "771a2a8b-01b2-4a21-998f-acbfa5129023",
      "target_period_id": "1b2a75c3-1a4a-4219-98f9-a1b2c8a7f1a2",
      "curriculum_id": "8c12a87c-12b4-4b51-912f-2cb1a581e281"
    }
  }
  ```

### C. Event Persetujuan KRS (`academic.krs_approved`)
* **Publisher**: Academic Service
* **Consumer**: LMS, Portal, Finance
* **JSON Payload**:
  ```json
  {
    "event_id": "f5129023-e2ba-412f-88ff-acbfa8712412",
    "event_name": "academic.krs_approved",
    "event_version": "1.0.0",
    "event_key": "krs:8c12a87c-12b4-4b51-912f-2cb1a581e281:approved",
    "timestamp": "2026-07-01T10:30:00Z",
    "data": {
      "krs_id": "8c12a87c-12b4-4b51-912f-2cb1a581e281",
      "student_id": "01b2a75c-31a4-4a21-998f-9a1b2c8a7f1a",
      "nim": "2601010045",
      "academic_period_id": "1b2a75c3-1a4a-4219-98f9-a1b2c8a7f1a2",
      "enrolled_classes": [
        {
          "class_id": "11223344-5566-7788-9900-aabbccddeeff",
          "class_code": "INF101-A",
          "course_name": "Algoritma & Struktur Data",
          "lecturer_id": "77aa88bb-ccdd-eeff-1122-334455667788"
        }
      ]
    }
  }
  ```

---

## 3. Kebijakan Keamanan API & Otentikasi
* **JWT Signature Verification**: Setiap backend microservice wajib melakukan verifikasi integritas JWT secara lokal menggunakan JWKS (*JSON Web Key Set*). URL JWKS dikonfigurasi di env: `AUTH_JWKS_URL=http://unsia-core-service/api/v1/auth/jwks`.
* **Access Scopes**: Klaim scope pada JWT mematuhi format RBAC Core: `[read:academic, write:academic.krs]`. Middleware dilarang meloloskan akses request jika scope token tidak mencukupi.

---

## 4. Mekanisme Retry, Dead Letter Queue (DLQ), & Rekonsiliasi
Untuk mengantisipasi kegagalan jaringan saat menyalurkan pesan event:

### A. Kebijakan Retry (Retry Policy)
Consumer harus mengimplementasikan **Exponential Backoff dengan Jitter** saat pemrosesan event gagal sementara (misal: database overload).
* **Batas Maksimal Percobaan**: 5 kali retry.
* **Interval Awal**: 2 detik (berganda $2^n$ hingga batas 60 detik).

### B. Dead Letter Queue (DLQ)
Jika batas maksimal retry terlampaui, event beserta informasi error trace wajib dirutekan ke antrean khusus `dlq.failed_events` untuk proses manual troubleshooting.
* **Log Error Trace**:
  ```json
  {
    "failed_at": "2026-07-01T10:35:00Z",
    "retry_count": 5,
    "last_error": "Connection timed out to lms_db after 5000ms",
    "original_payload": { ... }
  }
  ```

### C. Rekonsiliasi Data (Reconciliation Job)
Developer wajib membangun worker rekonsiliasi terjadwal (cron job harian/mingguan) untuk mendeteksi deviasi/mismatch antara master source of truth dengan local snapshot di database modul lain.
* **Mismatch Action**: Perbedaan data akan memicu alert warning di Slack/Telegram admin dan status data mismatch ditandai `pending_review` di UI konsol teknis.
