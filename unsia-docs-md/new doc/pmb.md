# PRD Final - PMB Module

## 1. Pendahuluan & Latar Belakang
Modul Penerimaan Mahasiswa Baru (PMB) melayani proses seleksi penerimaan mahasiswa secara transparan dan teratur, mulai dari pengisian biodata diri, pemenuhan berkas administrasi, ujian CBT masuk, hingga handover data mahasiswa baru ke modul SIAKAD.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `pmb_db`
* **Source of Truth (Kepemilikan Data)**:
  * `pmb.applicants` (Data pendaftar PMB).
  * `pmb.applicant_biodata` (Informasi biodata pendaftar).
  * `pmb.applicant_addresses` (Data alamat pendaftar).
  * `pmb.applicant_education_backgrounds` (Data sekolah/pendidikan asal).
  * `pmb.applicant_family_members` (Data keluarga pendaftar).
  * `pmb.applicant_financial_profiles` & `pmb.applicant_facility_profiles` (Profil penunjang belajar).
  * `pmb.applicant_documents` (Unggahan berkas pendaftaran).
  * `pmb.applicant_status_histories` (Histori transisi status PMB).
  * `pmb.re_registrations` (Registrasi daftar ulang).
  * `pmb.loa_documents` (Surat keputusan kelulusan).
  * `pmb.handover_logs` (Rekaman serah terima mahasiswa baru).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `person_snapshots` & `reference_snapshots`: Informasi personal dari Core dan lookup master prodi dari Referensi.
* `applicant_invoice_statuses`: Read model status tagihan keuangan pendaftar yang bersumber dari modul Finance secara asinkron.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-PMB-001** | P0 | Kemandirian Data Pendaftaran. | PMB menyimpan pelamar dan dokumen secara terpisah menggunakan person_ref_id dan reference snapshot lokal. |
| **PRD-PMB-002** | P0 | Sinkronisasi Asinkron Status Pembayaran. | PMB tidak melalukan query langsung ke finance_db; PMB memperbarui read model `applicant_invoice_statuses` dari event pembayaran. |
| **PRD-PMB-003** | P0 | Handover Mahasiswa Baru Idempotent. | Pengiriman draf mahasiswa baru ke modul SIAKAD harus bersifat idempotent dan aman ketika SIAKAD down. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `pmb.applicant_created`: Diterbitkan ketika pendaftaran awal pendaftar disimpan.
    * *Payload Minimum*: `applicant_id`, `person_ref_id`, `applicant_no`, `target_period_ref_id`, `study_program_ref_id`.
  * `pmb.ready_for_academic`: Diterbitkan setelah pendaftar melunasi tagihan registrasi awal dan siap digenerasikan NIM.
    * *Payload Minimum*: `applicant_id`, `person_ref_id`, `target_period_ref_id`, `study_program_ref_id`, `curriculum_candidate_ref`.

## 6. Degraded Mode & Resilience Guardrails
* **Finance Outage Resilience**: Jika Finance down, PMB tetap menerima pengisian biodata, dokumen, dan proses seleksi. Request tagihan biaya pendaftaran akan masuk ke tabel outbox PMB dan dikirim ulang secara otomatis setelah Finance online.
