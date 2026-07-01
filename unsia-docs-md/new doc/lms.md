# PRD Final - LMS Module

## 1. Pendahuluan & Latar Belakang
Modul LMS memfasilitasi kegiatan belajar mengajar online secara interaktif, menyediakan ruang materi ajar, pengunggahan tugas mingguan, kuis mandiri, forum diskusi dosen-mahasiswa, serta tracking keaktifan belajar online.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `lms_db`
* **Source of Truth (Kepemilikan Data)**:
  * `lms.classes` (Proyeksi kelas online).
  * `lms.enrollments` (Data kepesertaan kelas online).
  * `lms.sessions` (Pertemuan perkuliahan mingguan).
  * `lms.materials` & `lms.videos` (File materi perkuliahan).
  * `lms.vicon_links` (Akses video conference kuliah).
  * `lms.assignments` & `lms.assignment_submissions` (Manajemen tugas belajar).
  * `lms.quiz_activities` (Kuis mandiri terhubung dengan CBT Assessment).
  * `lms.discussions` & `lms.discussion_comments` (Forum tanya-jawab kelas).
  * `lms.attendances` (Presensi belajar sesi online).
  * `lms.learning_progress` (Rasio penyelesaian modul mahasiswa).
  * `lms.grade_syncs` (Log penyelarasan nilai tugas).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `academic_class_snapshots`: Salinan draf kelas penawaran dari Akademik.
* `student_snapshots` & `lecturer_snapshots`: Salinan profil mahasiswa dan dosen pengampu.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-LMS-001** | P0 | Larangan Penciptaan Kelas Mandiri. | LMS tidak boleh membuat kelas akademik sendiri; seluruh ruang kelas online harus merujuk snapshot kelas dari SIAKAD. |
| **PRD-LMS-002** | P0 | Kepesertaan Otomatis (Auto-Enroll). | LMS mendaftarkan peserta kelas secara otomatis sesaat setelah menerima event `academic.krs_approved`. |
| **PRD-LMS-003** | P0 | Pengiriman Nilai Tugas Mandiri. | Pengiriman nilai dari LMS ke Akademik harus dikirim via event secara idempotent untuk mencegah penimpaan nilai final. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Dikonsumsi (Consumer)**:
  * `academic.class_opened`: Digunakan untuk membuka ruang kelas online LMS baru.
  * `academic.krs_approved`: Digunakan untuk melakukan auto-enroll mahasiswa ke kelas LMS.
* **Event yang Diterbitkan (Publisher)**:
  * `lms.grade_input_submitted`: Dikirim setelah dosen menyelesaikan koreksi nilai tugas perkuliahan di LMS.
    * *Payload Minimum*: `student_ref_id`, `course_offering_ref_id`, `source_ref_id`, `score`, `submitted_at`.

## 6. Degraded Mode & Resilience Guardrails
* **Academic Outage Resilience**: Jika SIAKAD down, kelas kuliah LMS yang sudah aktif tetap berjalan normal. Dosen tetap dapat mengunggah materi dan menilai tugas. Input nilai tugas baru ditahan dalam antrean outbox LMS dan akan disinkronisasikan otomatis setelah SIAKAD pulih.
