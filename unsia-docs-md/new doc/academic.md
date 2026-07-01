# PRD Final - Akademik Module

## 1. Pendahuluan & Latar Belakang
Modul Akademik (SIAKAD) bertindak sebagai jantung operasional perkuliahan mahasiswa, kurikulum prodi, persetujuan KRS oleh dosen pembimbing akademik, serta pelaporan kelulusan nilai mahasiswa (KHS/Transkrip).

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `academic_db`
* **Source of Truth (Kepemilikan Data)**:
  * `academic.students` (Data induk profil mahasiswa).
  * `academic.student_advisors` (Pemetaan pembimbing akademik / dosen wali).
  * `academic.nim_format_configs` & `academic.nim_sequences` (Konfigurasi pembuatan NIM otomatis).
  * `academic.curriculums` (Kurikulum pendidikan).
  * `academic.courses` & `academic.curriculum_courses` (Mata kuliah & pemetaan semester).
  * `academic.course_offerings` (Penawaran mata kuliah semester).
  * `academic.classes` & `academic.class_lecturers` (Kelas paralel & plot pengampu).
  * `academic.class_schedules` (Jadwal hari & jam kelas).
  * `academic.krs` & `academic.krs_items` (Rancangan KRS diajukan).
  * `academic.grades` & `academic.grade_histories` (Nilai akhir perkuliahan).
  * `academic.khs` & `academic.transcripts` (Kartu hasil studi & transkrip nilai akumulatif).
  * `academic.yudisium_records` & `academic.alumni` (Yudisium kelulusan alumni).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `person_snapshots` & `reference_snapshots` dari Core dan Referensi.
* `student_clearance_snapshots`: Salinan status clearance keuangan mahasiswa dari Finance untuk validasi kelaikan isi KRS.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-ACA-001** | P0 | SIAKAD Otoritas Nilai Kelulusan. | Final grade mahasiswa sepenuhnya dikuasai SIAKAD; nilai mentah dari LMS hanya diperlakukan sebagai usulan. |
| **PRD-ACA-002** | P0 | Integrasi Clearance Asinkron. | SIAKAD dilarang menembak database Finance langsung; validasi keuangan mematuhi snapshot clearance lokal. |
| **PRD-ACA-003** | P0 | Provisioning Kelas & KRS ke LMS. | Setiap persetujuan KRS wajib memicu trigger pendaftaran otomatis akun dan enroll mahasiswa ke ruang kelas LMS terkait. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `academic.student_created`: Dikirim sesaat setelah data mahasiswa baru/NIM di-generate.
    * *Payload Minimum*: `student_id`, `person_ref_id`, `nim`, `entry_period_ref_id`, `study_program_ref_id`, `curriculum_id`.
  * `academic.krs_approved`: Dikirim saat draf KRS mahasiswa selesai disetujui dosen PA.
    * *Payload Minimum*: `student_id`, `krs_id`, `academic_period_ref_id`, `krs_item_ids`.

## 6. Degraded Mode & Resilience Guardrails
* **Finance Clearance Fallback**: Jika Finance down, pengisian KRS tetap dapat dilangsungkan. Status kelayakan pembayaran akan membaca data `student_clearance_snapshots` lokal terakhir. Jika data snapshot tidak ditemukan, KRS mahasiswa dialihkan menjadi status `PENDING_REVIEW` untuk diuji ulang saat Finance pulih.
