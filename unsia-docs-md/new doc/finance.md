# PRD Final - Finance Module

## 1. Pendahuluan & Latar Belakang
Modul Finance menjamin transparansi dan keakuratan pengelolaan keuangan kampus, mencakup pembuatan invoice tagihan, penerimaan pembayaran dari payment gateway luar, rekonsiliasi kas bank, jurnal akuntansi pembukuan berpasangan, hingga evaluasi status pembebasan akademik mahasiswa.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `finance_db`
* **Source of Truth (Kepemilikan Data)**:
  * `finance.invoices` & `finance.invoice_items` (Invoice tagihan biaya pendaftaran/UKT).
  * `finance.payments` (Transaksi pembayaran masuk).
  * `finance.payment_gateway_callbacks` & `finance.payment_verifications` (Log callback PG dan verifikasi manual).
  * `finance.scholarships` (Potongan biaya beasiswa).
  * `finance.installment_requests` (Pengajuan cicilan pembayaran).
  * `finance.student_clearances` & `finance.clearance_dispensations` (Status kelayakan finansial mahasiswa).
  * `finance.cash_accounts` & `finance.cash_transactions` (Mutasi kas).
  * `finance.coa_accounts`, `finance.journals` & `finance.journal_entries` (Buku besar akuntansi).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `customer_snapshots`: Menyimpan profil identitas pembayar (NIM/No Registrasi, nama lengkap, status keaktifan) dari PMB atau Akademik untuk mendukung pembuatan invoice tanpa cross-DB query.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-FIN-001** | P0 | Finance Otoritas Tunggal Transaksi. | Modul Keuangan merupakan satu-satunya source of truth untuk tagihan, pembayaran, kuitansi, dan jurnal akuntansi. |
| **PRD-FIN-002** | P0 | Idempotency Callback Payment Gateway. | Callback transaksi PG harus diamankan menggunakan idempotency key berdasarkan kode provider dan event ID. |
| **PRD-FIN-003** | P0 | Penerbitan Clearance Lintas Modul. | Setiap pelunasan biaya wajib mempublikasikan perubahan status kelayakan akademik mahasiswa ke modul Akademik/LMS. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `finance.invoice_created`: Dikirim sesaat setelah invoice resmi diterbitkan.
    * *Payload Minimum*: `invoice_id`, `invoice_no`, `bill_to_type`, `bill_to_ref_id`, `amount_total`, `status_code`.
  * `finance.payment_paid`: Dikirim saat transaksi pembayaran dikonfirmasi lunas.
    * *Payload Minimum*: `payment_id`, `invoice_id`, `bill_to_type`, `bill_to_ref_id`, `paid_amount`, `paid_at`.
  * `finance.clearance_changed`: Dikirim ketika status kelayakan layanan akademik mahasiswa mengalami perubahan.
    * *Payload Minimum*: `subject_type`, `subject_ref_id`, `academic_period_ref_id`, `service_code`, `status_code`.

## 6. Degraded Mode & Resilience Guardrails
* **Standalone Billing Operation**: Jika modul PMB atau Akademik down, modul Finance tetap dapat memproses transaksi pembayaran atas invoice yang sudah terbit menggunakan data `customer_snapshots` lokal terakhir. Jurnal akuntansi dan kuitansi tetap diterbitkan secara normal.
