# PRD Final - Referensi Module

## 1. Pendahuluan & Latar Belakang
Modul Referensi bertanggung jawab mengelola master data umum yang digunakan secara kolektif oleh modul-modul lain di dalam ekosistem UNSIA ERP, seperti data program studi, tahun akademik, periode semester, wilayah geografis, agama, komponen biaya, dan kode status.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `reference_db`
* **Source of Truth (Kepemilikan Data)**:
  * `ref.countries`, `ref.provinces`, `ref.cities`, `ref.districts`, `ref.villages` (Master wilayah geografis).
  * `ref.religions` (Master agama).
  * `ref.study_programs` (Master program studi & fakultas).
  * `ref.academic_years` (Tahun ajaran kalender operasional).
  * `ref.academic_periods` (Periode semester aktif).
  * `ref.admission_paths` (Master jalur penerimaan PMB).
  * `ref.pmb_waves` (Master gelombang pendaftaran PMB).
  * `ref.document_types` (Jenis dokumen syarat PMB/Akademik).
  * `ref.payment_components` (Komponen biaya tagihan).
  * `ref.payment_methods` (Metode pembayaran valid).
  * `ref.employee_types` & `ref.lecturer_statuses` (Master status kepegawaian).
  * `ref.status_codes` (Kode status terstandarisasi).

## 3. Data Lintas Modul (Snapshot/Read Model)
Modul Referensi bersifat mandiri dan tidak menyimpan data snapshot dari modul lain secara lokal.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-REF-001** | P0 | Penerbitan event perubahan program studi dan periode akademik. | Modul Referensi wajib menerbitkan event outbox ketika data prodi atau periode diubah agar modul consumer dapat mengupdate snapshot. |
| **PRD-REF-002** | P0 | Local Reference Snapshot di consumer. | Setiap modul consumer (PMB, SIAKAD, dsb.) wajib menyimpan draf referensi lokal minimum agar operasi dasar berjalan saat database Referensi down. |
| **PRD-REF-003** | P1 | Aturan Valid Tanggal Efektif (*Effective Dating*). | Perubahan komponen master data sensitif tidak boleh merusak/memutus relasi historis pada transaksi keuangan/akademik lama. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `reference.study_program_updated`: Dikirim ketika program studi baru didaftarkan atau diperbarui.
    * *Payload Minimum*: `study_program_id`, `code`, `name`, `status_code`, `occurred_at`.
  * `reference.academic_period_updated`: Dikirim ketika periode akademik baru diaktifkan.
    * *Payload Minimum*: `academic_period_id`, `academic_year_id`, `code`, `name`, `term_code`, `status_code`.

## 6. Degraded Mode & Resilience Guardrails
* **Reference Outage Fallback**: Saat database Referensi down, modul lain tetap dapat beroperasi normal menggunakan data dari *local reference snapshot* terakhir. Pengeditan master referensi baru ditangguhkan sampai modul Referensi kembali online.
